// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./IPQVerify.sol";

/// @title PQVerifyTestHelper — Helper contract for end-to-end testing
/// @notice Wraps the PQ_VERIFY precompile for integration tests
contract PQVerifyTestHelper {
    address constant PQ_VERIFY_PRECOMPILE = address(0x0300000000000000000000000000000000000000);

    event VerificationResult(bool valid, uint256 gasUsed);

    /// @notice Verify a PQ signature and emit the result with gas used
    function verifyAndLog(
        bytes calldata pubkey,
        bytes calldata signature,
        bytes calldata message,
        uint8 algorithm
    ) external returns (bool valid) {
        uint256 gasBefore = gasleft();

        (bool success, bytes memory result) = PQ_VERIFY_PRECOMPILE.staticcall(
            abi.encode(pubkey, signature, message, algorithm)
        );

        uint256 gasUsed = gasBefore - gasleft();

        if (!success || result.length < 32) {
            emit VerificationResult(false, gasUsed);
            return false;
        }

        valid = abi.decode(result, (bool));
        emit VerificationResult(valid, gasUsed);
        return valid;
    }

    /// @notice Batch verify multiple signatures
    function batchVerify(
        bytes[] calldata pubkeys,
        bytes[] calldata signatures,
        bytes[] calldata messages,
        uint8[] calldata algorithms
    ) external view returns (bool[] memory results) {
        require(
            pubkeys.length == signatures.length &&
            signatures.length == messages.length &&
            messages.length == algorithms.length,
            "array length mismatch"
        );

        results = new bool[](pubkeys.length);
        for (uint256 i = 0; i < pubkeys.length; i++) {
            (bool success, bytes memory result) = PQ_VERIFY_PRECOMPILE.staticcall(
                abi.encode(pubkeys[i], signatures[i], messages[i], algorithms[i])
            );
            if (success && result.length >= 32) {
                results[i] = abi.decode(result, (bool));
            }
        }
    }
}
