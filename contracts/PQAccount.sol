// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./IPQVerify.sol";

/// @title UserOperation — ERC-4337 user operation struct
/// @dev Simplified for Beam's account abstraction layer.
///      Full ERC-4337 uses PackedUserOperation; this covers the core fields.
struct UserOperation {
    address sender;
    uint256 nonce;
    bytes initCode;
    bytes callData;
    uint256 callGasLimit;
    uint256 verificationGasLimit;
    uint256 preVerificationGas;
    uint256 maxFeePerGas;
    uint256 maxPriorityFeePerGas;
    bytes paymasterAndData;
    bytes signature;
}

/// @title IAccount — ERC-4337 account interface
/// @dev See https://eips.ethereum.org/EIPS/eip-4337
interface IAccount {
    /// @notice Validate a user operation's signature and nonce
    /// @param userOp The user operation to validate
    /// @param userOpHash Hash of the user operation (computed by EntryPoint)
    /// @param missingAccountFunds Funds the account must deposit to the EntryPoint
    /// @return validationData 0 for success, 1 for failure (SIG_VALIDATION_FAILED)
    function validateUserOp(
        UserOperation calldata userOp,
        bytes32 userOpHash,
        uint256 missingAccountFunds
    ) external returns (uint256 validationData);
}

/// @title IEntryPoint — Minimal EntryPoint interface for fund deposits
interface IEntryPoint {
    function depositTo(address account) external payable;
    function getNonce(address sender, uint192 key) external view returns (uint256);
}

/// @title PQAccount — Post-Quantum Smart Account with ERC-4337 Support
/// @notice Smart account that uses PQ signatures for transaction authorization,
///         compatible with ERC-4337 account abstraction infrastructure.
///         Demonstrates how game studios and dApps can leverage PQ signing without
///         changing their existing code — the SDK handles key management.
/// @dev Implements IAccount for EntryPoint compatibility. The PQ signature is
///      verified via the PQ_VERIFY precompile at 0x0300...0000.
contract PQAccount is IAccount {
    /// @notice ERC-4337 signature validation success/failure constants
    uint256 internal constant SIG_VALIDATION_SUCCESS = 0;
    uint256 internal constant SIG_VALIDATION_FAILED = 1;

    /// @notice Address of the PQ_VERIFY precompile
    address constant PQ_VERIFY_PRECOMPILE = address(0x0300000000000000000000000000000000000000);

    /// @notice The PQ public key that controls this account
    bytes public pqPublicKey;

    /// @notice The algorithm used for this account's signatures
    uint8 public pqAlgorithm;

    /// @notice Nonce for replay protection (used by direct execute())
    uint256 public nonce;

    /// @notice Owner address (derived from PQ public key hash)
    address public owner;

    /// @notice ERC-4337 EntryPoint address
    address public entryPoint;

    event PQAccountCreated(address indexed account, uint8 algorithm, bytes32 pubkeyHash);
    event PQTransactionExecuted(address indexed target, uint256 value, uint256 nonce);
    event EntryPointChanged(address indexed oldEntryPoint, address indexed newEntryPoint);

    error InvalidSignature();
    error InvalidNonce();
    error ExecutionFailed();
    error NotOwnerOrEntryPoint();
    error OnlyEntryPoint();

    modifier onlyOwnerOrEntryPoint() {
        if (msg.sender != owner && msg.sender != entryPoint) revert NotOwnerOrEntryPoint();
        _;
    }

    modifier onlyEntryPoint() {
        if (msg.sender != entryPoint) revert OnlyEntryPoint();
        _;
    }

    /// @param _pubkey The PQ public key for this account
    /// @param _algorithm The PQ algorithm (0 = ML-DSA-65, 1 = SLH-DSA-128s)
    /// @param _entryPoint The ERC-4337 EntryPoint contract address
    constructor(bytes memory _pubkey, uint8 _algorithm, address _entryPoint) {
        pqPublicKey = _pubkey;
        pqAlgorithm = _algorithm;
        owner = derivePQAddress(_pubkey);
        entryPoint = _entryPoint;
        emit PQAccountCreated(address(this), _algorithm, keccak256(_pubkey));
    }

    // ── ERC-4337 IAccount ──────────────────────────────────────────────

    /// @notice Validate a user operation (called by EntryPoint)
    /// @dev Verifies the PQ signature over the userOpHash. The signature field
    ///      in the UserOperation contains the raw PQ signature bytes.
    /// @param userOp The user operation containing the PQ signature
    /// @param userOpHash Hash of the user operation (domain-separated by EntryPoint)
    /// @param missingAccountFunds Funds to deposit to EntryPoint if needed
    /// @return validationData 0 for success, 1 for signature failure
    function validateUserOp(
        UserOperation calldata userOp,
        bytes32 userOpHash,
        uint256 missingAccountFunds
    ) external onlyEntryPoint returns (uint256 validationData) {
        // Verify PQ signature over the EntryPoint-provided hash
        bool valid = _verifyPQSignature(abi.encodePacked(userOpHash), userOp.signature);

        // Deposit missing funds to EntryPoint if required
        if (missingAccountFunds > 0) {
            (bool sent, ) = payable(msg.sender).call{value: missingAccountFunds}("");
            (sent); // ignore failure — EntryPoint will revert if deposit insufficient
        }

        return valid ? SIG_VALIDATION_SUCCESS : SIG_VALIDATION_FAILED;
    }

    // ── Direct Execution (non-4337 path) ───────────────────────────────

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
        if (!_verifyPQSignature(message, signature)) revert InvalidSignature();

        // Increment nonce
        nonce = _nonce + 1;

        // Execute the call
        (bool callSuccess, ) = target.call{value: value}(data);
        if (!callSuccess) revert ExecutionFailed();

        emit PQTransactionExecuted(target, value, _nonce);
    }

    /// @notice Execute a transaction via EntryPoint (4337 path, already validated)
    /// @param target The contract/address to call
    /// @param value ETH/BEAM value to send
    /// @param data Calldata for the target
    function executeFromEntryPoint(
        address target,
        uint256 value,
        bytes calldata data
    ) external onlyEntryPoint {
        (bool success, ) = target.call{value: value}(data);
        if (!success) revert ExecutionFailed();
        emit PQTransactionExecuted(target, value, nonce);
    }

    /// @notice Execute a batch of transactions via EntryPoint
    /// @param targets Array of addresses to call
    /// @param values Array of values to send
    /// @param datas Array of calldatas
    function executeBatchFromEntryPoint(
        address[] calldata targets,
        uint256[] calldata values,
        bytes[] calldata datas
    ) external onlyEntryPoint {
        require(targets.length == values.length && values.length == datas.length, "length mismatch");
        for (uint256 i = 0; i < targets.length; i++) {
            (bool success, ) = targets[i].call{value: values[i]}(datas[i]);
            if (!success) revert ExecutionFailed();
            emit PQTransactionExecuted(targets[i], values[i], nonce);
        }
    }

    // ── Signature Verification ─────────────────────────────────────────

    /// @notice Check if a PQ signature is valid for a message (ERC-1271 style)
    /// @param message The signed message bytes
    /// @param signature The PQ signature
    /// @return True if signature is valid
    function isValidSignature(
        bytes calldata message,
        bytes calldata signature
    ) external view returns (bool) {
        return _verifyPQSignature(message, signature);
    }

    /// @dev Internal PQ signature verification via precompile
    function _verifyPQSignature(
        bytes memory message,
        bytes calldata signature
    ) internal view returns (bool) {
        (bool success, bytes memory result) = PQ_VERIFY_PRECOMPILE.staticcall(
            abi.encode(pqPublicKey, signature, message, pqAlgorithm)
        );
        if (!success || result.length < 32) return false;
        return abi.decode(result, (bool));
    }

    // ── Admin ──────────────────────────────────────────────────────────

    /// @notice Update the EntryPoint address (only owner or current EntryPoint)
    function setEntryPoint(address _newEntryPoint) external onlyOwnerOrEntryPoint {
        emit EntryPointChanged(entryPoint, _newEntryPoint);
        entryPoint = _newEntryPoint;
    }

    /// @notice Derive an address from a PQ public key (keccak256 last 20 bytes)
    function derivePQAddress(bytes memory pubkey) public pure returns (address) {
        return address(uint160(uint256(keccak256(pubkey))));
    }

    receive() external payable {}
}
