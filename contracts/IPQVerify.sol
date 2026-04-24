// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/// @title IPQVerify — Post-Quantum Signature Verification Precompile
/// @notice Interface for the PQ_VERIFY precompile deployed at 0x0000...0b00
/// @dev This precompile verifies ML-DSA-65 (Dilithium) and SLH-DSA-128s (SPHINCS+)
///      signatures natively on the Beam Subnet-EVM. It is stateless and read-only.
///
///      Algorithm IDs:
///        0 = ML-DSA-65   (NIST FIPS 204) — primary, lattice-based
///        1 = SLH-DSA-128s (NIST FIPS 205) — fallback, hash-based
///
///      Gas costs (benchmarked on Apple M1 Pro, 10x safety margin):
///        ML-DSA-65:    ~133,600 gas
///        SLH-DSA-128s: ~523,600 gas
interface IPQVerify {
    /// @notice Verify a post-quantum digital signature
    /// @param pubkey The signer's public key (ML-DSA-65: 1952 bytes, SLH-DSA-128s: 32 bytes)
    /// @param signature The signature to verify (ML-DSA-65: 3309 bytes, SLH-DSA-128s: 7856 bytes)
    /// @param message The message that was signed
    /// @param algorithm The PQ algorithm ID (0 = ML-DSA-65, 1 = SLH-DSA-128s)
    /// @return valid True if the signature is valid, false otherwise
    function pqVerify(
        bytes calldata pubkey,
        bytes calldata signature,
        bytes calldata message,
        uint8 algorithm
    ) external view returns (bool valid);
}
