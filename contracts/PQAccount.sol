// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./IPQVerify.sol";

/// @title PQAccount — Post-Quantum Smart Account (Proof of Concept)
/// @notice Minimal smart account that uses PQ signatures for transaction authorization.
///         Demonstrates how game studios and dApps can leverage PQ signing without
///         changing their existing code — the SDK handles key management.
/// @dev This is a PoC for the grant. Production version will be integrated into
///      Beam's account abstraction layer via SDK updates.
contract PQAccount {
    /// @notice Address of the PQ_VERIFY precompile
    address constant PQ_VERIFY_PRECOMPILE = address(0x0b00);

    /// @notice The PQ public key that controls this account
    bytes public pqPublicKey;

    /// @notice The algorithm used for this account's signatures
    uint8 public pqAlgorithm;

    /// @notice Nonce for replay protection
    uint256 public nonce;

    /// @notice Owner address (derived from PQ public key hash)
    address public owner;

    event PQAccountCreated(address indexed account, uint8 algorithm, bytes32 pubkeyHash);
    event PQTransactionExecuted(address indexed target, uint256 value, uint256 nonce);

    error InvalidSignature();
    error InvalidNonce();
    error ExecutionFailed();
    error NotOwner();

    /// @param _pubkey The PQ public key for this account
    /// @param _algorithm The PQ algorithm (0 = ML-DSA-65, 1 = SLH-DSA-128s)
    constructor(bytes memory _pubkey, uint8 _algorithm) {
        pqPublicKey = _pubkey;
        pqAlgorithm = _algorithm;
        owner = derivePQAddress(_pubkey);
        emit PQAccountCreated(address(this), _algorithm, keccak256(_pubkey));
    }

    /// @notice Execute a transaction after verifying PQ signature
    /// @param target The contract/address to call
    /// @param value ETH/BEAM value to send
    /// @param data Calldata for the target
    /// @param _nonce Expected nonce (for replay protection)
    /// @param signature PQ signature over (target, value, data, nonce, chainId)
    function execute(
        address target,
        uint256 value,
        bytes calldata data,
        uint256 _nonce,
        bytes calldata signature
    ) external {
        if (_nonce != nonce) revert InvalidNonce();

        // Construct the message that was signed
        bytes memory message = abi.encodePacked(
            target,
            value,
            keccak256(data),
            _nonce,
            block.chainid
        );

        // Verify PQ signature via precompile
        (bool success, bytes memory result) = PQ_VERIFY_PRECOMPILE.staticcall(
            abi.encode(pqPublicKey, signature, message, pqAlgorithm)
        );

        if (!success || result.length < 32) revert InvalidSignature();
        bool valid = abi.decode(result, (bool));
        if (!valid) revert InvalidSignature();

        // Increment nonce
        nonce = _nonce + 1;

        // Execute the call
        (bool callSuccess, ) = target.call{value: value}(data);
        if (!callSuccess) revert ExecutionFailed();

        emit PQTransactionExecuted(target, value, _nonce);
    }

    /// @notice Derive an address from a PQ public key (keccak256 → last 20 bytes)
    /// @dev Uses the same derivation as Ethereum but from PQ key material
    function derivePQAddress(bytes memory pubkey) public pure returns (address) {
        return address(uint160(uint256(keccak256(pubkey))));
    }

    /// @notice Check if a PQ signature is valid for a message
    function isValidSignature(
        bytes calldata message,
        bytes calldata signature
    ) external view returns (bool) {
        (bool success, bytes memory result) = PQ_VERIFY_PRECOMPILE.staticcall(
            abi.encode(pqPublicKey, signature, message, pqAlgorithm)
        );
        if (!success || result.length < 32) return false;
        return abi.decode(result, (bool));
    }

    receive() external payable {}
}
