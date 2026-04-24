// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./IPQVerify.sol";

/// @title PQKeyRotation — Post-Quantum Key Migration & Rotation Registry
/// @notice Allows users to register PQ public keys, rotate to new keys, and
///         migrate from ECDSA-controlled accounts to PQ-controlled accounts.
///         Includes a timelock on rotations to prevent immediate takeover if
///         a key is compromised.
/// @dev Uses PQ_VERIFY precompile to verify that rotation requests are signed
///      by the current key holder. ECDSA→PQ migration uses msg.sender (ECDSA)
///      as authorization; PQ→PQ rotation uses PQ signature verification.
contract PQKeyRotation {
    /// @notice Address of the PQ_VERIFY precompile
    address constant PQ_VERIFY_PRECOMPILE = address(0x0300000000000000000000000000000000000000);

    /// @notice Timelock duration for key rotations (default: 24 hours)
    /// @dev Gives key holders time to cancel a malicious rotation
    uint256 public constant ROTATION_TIMELOCK = 24 hours;

    /// @notice Represents a registered PQ key binding
    struct KeyBinding {
        bytes pqPublicKey;      // The current PQ public key
        uint8 algorithm;        // Algorithm ID (0 = ML-DSA-65, 1 = SLH-DSA-128s)
        uint256 registeredAt;   // Timestamp of registration
        bool active;            // Whether this binding is active
    }

    /// @notice Represents a pending key rotation
    struct PendingRotation {
        bytes newPublicKey;     // The proposed new PQ public key
        uint8 newAlgorithm;     // Algorithm for the new key
        uint256 executeAfter;   // Timestamp after which rotation can execute
        bool exists;            // Whether a pending rotation exists
    }

    /// @notice Mapping from account address to their PQ key binding
    mapping(address => KeyBinding) public keyBindings;

    /// @notice Mapping from account address to pending rotation
    mapping(address => PendingRotation) public pendingRotations;

    /// @notice Mapping from PQ address (derived from pubkey) to account
    mapping(address => address) public pqAddressToAccount;

    // ── Events ─────────────────────────────────────────────────────────

    event KeyRegistered(
        address indexed account,
        address indexed pqAddress,
        uint8 algorithm,
        bytes32 pubkeyHash
    );

    event RotationRequested(
        address indexed account,
        bytes32 newPubkeyHash,
        uint8 newAlgorithm,
        uint256 executeAfter
    );

    event RotationExecuted(
        address indexed account,
        address indexed newPqAddress,
        bytes32 oldPubkeyHash,
        bytes32 newPubkeyHash
    );

    event RotationCancelled(address indexed account);

    event KeyRevoked(address indexed account, bytes32 pubkeyHash);

    // ── Errors ─────────────────────────────────────────────────────────

    error AlreadyRegistered();
    error NotRegistered();
    error InvalidSignature();
    error RotationNotReady();
    error NoRotationPending();
    error RotationAlreadyPending();
    error KeyAlreadyBound();
    error NotAuthorized();

    // ── ECDSA → PQ Migration ───────────────────────────────────────────

    /// @notice Register a PQ public key for the caller's ECDSA account
    /// @dev This is the entry point for ECDSA→PQ migration. The caller
    ///      proves ownership of their ECDSA account via msg.sender.
    /// @param pubkey The PQ public key to bind to this account
    /// @param algorithm The PQ algorithm (0 = ML-DSA-65, 1 = SLH-DSA-128s)
    function registerKey(bytes calldata pubkey, uint8 algorithm) external {
        if (keyBindings[msg.sender].active) revert AlreadyRegistered();

        address pqAddr = derivePQAddress(pubkey);
        if (pqAddressToAccount[pqAddr] != address(0)) revert KeyAlreadyBound();

        keyBindings[msg.sender] = KeyBinding({
            pqPublicKey: pubkey,
            algorithm: algorithm,
            registeredAt: block.timestamp,
            active: true
        });

        pqAddressToAccount[pqAddr] = msg.sender;

        emit KeyRegistered(msg.sender, pqAddr, algorithm, keccak256(pubkey));
    }

    // ── PQ → PQ Key Rotation ───────────────────────────────────────────

    /// @notice Request rotation to a new PQ key (starts timelock)
    /// @dev Must be signed by the CURRENT PQ key to prove ownership.
    ///      The rotation executes after ROTATION_TIMELOCK seconds.
    /// @param newPubkey The new PQ public key
    /// @param newAlgorithm Algorithm for the new key
    /// @param signature PQ signature by current key over (account, newPubkeyHash, chainId)
    function requestRotation(
        bytes calldata newPubkey,
        uint8 newAlgorithm,
        bytes calldata signature
    ) external {
        KeyBinding storage binding = keyBindings[msg.sender];
        if (!binding.active) revert NotRegistered();
        if (pendingRotations[msg.sender].exists) revert RotationAlreadyPending();

        // Verify the rotation request is signed by the current key
        bytes memory message = abi.encodePacked(
            msg.sender,
            keccak256(newPubkey),
            newAlgorithm,
            block.chainid
        );

        if (!_verifyPQ(binding.pqPublicKey, signature, message, binding.algorithm)) {
            revert InvalidSignature();
        }

        // Ensure the new key isn't already bound to another account
        address newPqAddr = derivePQAddress(newPubkey);
        if (pqAddressToAccount[newPqAddr] != address(0) && pqAddressToAccount[newPqAddr] != msg.sender) {
            revert KeyAlreadyBound();
        }

        uint256 executeAfter = block.timestamp + ROTATION_TIMELOCK;

        pendingRotations[msg.sender] = PendingRotation({
            newPublicKey: newPubkey,
            newAlgorithm: newAlgorithm,
            executeAfter: executeAfter,
            exists: true
        });

        emit RotationRequested(msg.sender, keccak256(newPubkey), newAlgorithm, executeAfter);
    }

    /// @notice Execute a pending rotation after timelock expires
    /// @dev Can be called by anyone after the timelock — the authorization
    ///      was verified when the rotation was requested.
    function executeRotation(address account) external {
        PendingRotation storage pending = pendingRotations[account];
        if (!pending.exists) revert NoRotationPending();
        if (block.timestamp < pending.executeAfter) revert RotationNotReady();

        KeyBinding storage binding = keyBindings[account];
        bytes32 oldHash = keccak256(binding.pqPublicKey);

        // Remove old PQ address mapping
        address oldPqAddr = derivePQAddress(binding.pqPublicKey);
        delete pqAddressToAccount[oldPqAddr];

        // Update binding
        binding.pqPublicKey = pending.newPublicKey;
        binding.algorithm = pending.newAlgorithm;
        binding.registeredAt = block.timestamp;

        // Set new PQ address mapping
        address newPqAddr = derivePQAddress(pending.newPublicKey);
        pqAddressToAccount[newPqAddr] = account;

        bytes32 newHash = keccak256(pending.newPublicKey);

        // Clean up pending
        delete pendingRotations[account];

        emit RotationExecuted(account, newPqAddr, oldHash, newHash);
    }

    /// @notice Cancel a pending rotation (must be signed by current key)
    /// @param signature PQ signature by current key over (account, "cancel", chainId)
    function cancelRotation(bytes calldata signature) external {
        KeyBinding storage binding = keyBindings[msg.sender];
        if (!binding.active) revert NotRegistered();
        if (!pendingRotations[msg.sender].exists) revert NoRotationPending();

        // Verify cancellation is authorized by current key
        bytes memory message = abi.encodePacked(
            msg.sender,
            bytes("cancel"),
            block.chainid
        );

        if (!_verifyPQ(binding.pqPublicKey, signature, message, binding.algorithm)) {
            revert InvalidSignature();
        }

        delete pendingRotations[msg.sender];
        emit RotationCancelled(msg.sender);
    }

    // ── Key Revocation ─────────────────────────────────────────────────

    /// @notice Revoke a PQ key binding (emergency use)
    /// @dev Must be called from the ECDSA account (msg.sender) as a
    ///      fallback if the PQ key is compromised. This also cancels
    ///      any pending rotation.
    function revokeKey() external {
        KeyBinding storage binding = keyBindings[msg.sender];
        if (!binding.active) revert NotRegistered();

        bytes32 pubkeyHash = keccak256(binding.pqPublicKey);
        address pqAddr = derivePQAddress(binding.pqPublicKey);

        // Clean up mappings
        delete pqAddressToAccount[pqAddr];
        delete pendingRotations[msg.sender];
        binding.active = false;

        emit KeyRevoked(msg.sender, pubkeyHash);
    }

    // ── View Functions ─────────────────────────────────────────────────

    /// @notice Get the PQ public key for an account
    function getPublicKey(address account) external view returns (bytes memory, uint8, bool) {
        KeyBinding storage binding = keyBindings[account];
        return (binding.pqPublicKey, binding.algorithm, binding.active);
    }

    /// @notice Check if an account has a registered PQ key
    function isRegistered(address account) external view returns (bool) {
        return keyBindings[account].active;
    }

    /// @notice Check if a rotation is pending and when it can execute
    function getRotationStatus(address account) external view returns (bool pending, uint256 executeAfter) {
        PendingRotation storage rotation = pendingRotations[account];
        return (rotation.exists, rotation.executeAfter);
    }

    /// @notice Derive a PQ address from a public key
    function derivePQAddress(bytes memory pubkey) public pure returns (address) {
        return address(uint160(uint256(keccak256(pubkey))));
    }

    // ── Internal ───────────────────────────────────────────────────────

    /// @dev Verify a PQ signature via the precompile
    function _verifyPQ(
        bytes storage pubkey,
        bytes calldata signature,
        bytes memory message,
        uint8 algorithm
    ) internal view returns (bool) {
        (bool success, bytes memory result) = PQ_VERIFY_PRECOMPILE.staticcall(
            abi.encode(pubkey, signature, message, algorithm)
        );
        if (!success || result.length < 32) return false;
        return abi.decode(result, (bool));
    }
}
