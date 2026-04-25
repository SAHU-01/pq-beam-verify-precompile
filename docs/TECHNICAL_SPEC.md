# PQ_VERIFY Precompile — Technical Specification

**Version:** 0.1.0  
**Status:** Draft  
**Author:** Ankita Sahu (@SAHU-01)  
**Date:** April 2026

---

## 1. Overview

The PQ_VERIFY precompile adds native post-quantum signature verification to the Beam Network's Subnet-EVM. It supports two NIST-standardized algorithms:

| Algorithm | Standard | Type | Security Level |
|-----------|----------|------|----------------|
| ML-DSA-65 (Dilithium) | NIST FIPS 204 | Lattice-based | NIST Level 3 (AES-192 equivalent) |
| SLH-DSA-128s (SPHINCS+) | NIST FIPS 205 | Hash-based | NIST Level 1 (AES-128 equivalent) |

ML-DSA-65 is the primary algorithm (compact signatures, fast verification). SLH-DSA-128s is the conservative fallback (hash-based, minimal cryptographic assumptions).

## 2. Precompile Interface

### Address

```
0x0300000000000000000000000000000000000000
```

### ABI

```solidity
function pqVerify(
    bytes calldata pubkey,
    bytes calldata signature,
    bytes calldata message,
    uint8 algorithm
) external view returns (bool valid);
```

### Algorithm IDs

| ID | Algorithm | Public Key Size | Signature Size |
|----|-----------|----------------|----------------|
| 0  | ML-DSA-65 | 1,952 bytes | 3,309 bytes |
| 1  | SLH-DSA-128s | 32 bytes | 7,856 bytes |

### Input Encoding (ABI)

The input follows standard Solidity ABI encoding for `(bytes, bytes, bytes, uint8)`:

```
Offset 0x00:  offset to pubkey data (uint256)
Offset 0x20:  offset to signature data (uint256)
Offset 0x40:  offset to message data (uint256)
Offset 0x60:  algorithm ID (uint8, right-padded to 32 bytes)
Offset 0x80+: dynamic data (length-prefixed, 32-byte padded)
```

### Output Encoding

```
Offset 0x00:  valid (bool, right-padded to 32 bytes)
```

Returns `0x01` if verification succeeds, `0x00` if verification fails.

### Error Behavior

- Invalid algorithm ID → returns `false` (no revert)
- Malformed input (too short) → reverts
- Invalid key/signature sizes → returns `false`
- Empty message → valid (if signature matches)

## 3. Gas Schedule

### Gas Costs

| Operation | Gas Cost | Rationale |
|-----------|----------|-----------|
| Base overhead | 3,600 | ABI decoding + dispatch |
| ML-DSA-65 verify | 130,000 | 4.2x ecrecover (3,000) with 10x safety margin |
| SLH-DSA-128s verify | 520,000 | 17.3x ecrecover (3,000) with 10x safety margin |

Total gas = `base + algorithm_cost`.

### Benchmark Methodology

Gas costs are derived from CPU time benchmarks on Apple M1 Pro reference hardware:
- ecrecover: ~25μs = 3,000 gas (Ethereum standard)
- ML-DSA-65 verify: ~105μs → 4.2x ecrecover → 12,600 raw gas → 130,000 with 10x margin
- SLH-DSA-128s verify: ~432μs → 17.3x ecrecover → 51,900 raw gas → 520,000 with 10x margin

The 10x safety margin accounts for validator hardware variance (ARM vs x86, different CPU generations), different liboqs builds, and provides conservative DoS protection. SLH-DSA-128s correctly costs more gas than ML-DSA-65 because it is the slower algorithm (~4x slower at verification).

The current gas formula accounts for CPU compute time only. Memory expansion costs for large signatures (SLH-DSA-128s at 7,856 bytes) are not separately metered — the EVM's native `memory_cost = (memory_size_word ** 2) / 512 + (3 * memory_size_word)` applies to calldata loading, but the precompile's internal buffer allocation is covered by the safety margin. This assumption must be validated on production hardware.

Final gas values will be calibrated on Beam validator hardware before mainnet activation. A benchmark matrix across AWS EC2 instances (t3.large, c5.xlarge) and bare-metal Linux servers is required to ensure gas costs are safe for the weakest validator in the set.

## 4. PQ Account Format

### Address Derivation

PQ addresses are derived identically to ECDSA addresses but from PQ key material:

```
pq_address = keccak256(pq_public_key)[12:32]  // last 20 bytes
```

This ensures PQ addresses are the same length as standard Ethereum addresses and compatible with existing infrastructure (explorers, wallets, contracts).

### Account Distinguishing

PQ accounts are distinguished by their transaction type, not their address format. The EIP-2718 typed envelope (type `0x50`) signals that the transaction carries PQ signature fields.

## 5. PQ Transaction Type

### Type ID: `0x50`

Chosen as "P" for Post-Quantum. This is a new EIP-2718 typed transaction envelope.

### Envelope Format

```
0x50 || RLP([
    chainId,
    nonce,
    gasPrice,
    gasLimit,
    to,
    value,
    data,
    pqAlgorithm,    // uint8
    pqPublicKey,    // bytes
    pqSignature     // bytes
])
```

### Signing

The transaction hash for signing is:

```
txHash = keccak256(0x50 || RLP([
    chainId,
    nonce,
    gasPrice,
    gasLimit,
    to,
    value,
    data,
    pqAlgorithm
]))
```

The signature is computed as:
```
pqSignature = PQ_Sign(pqSecretKey, txHash)
```

### Validation

Nodes validate PQ transactions by:
1. Decode the typed envelope
2. Reconstruct the signing hash from unsigned fields
3. Call `PQ_VERIFY(pqPublicKey, pqSignature, txHash, pqAlgorithm)`
4. Derive sender address as `keccak256(pqPublicKey)[12:32]`
5. Continue with standard EVM execution

## 6. Migration Path

### Phase 1: Precompile Only (This Grant)

- PQ_VERIFY precompile deployed at `0x0300`
- PQ accounts via smart contracts (PQAccount.sol)
- Game studios and dApps can opt into PQ signing via SDK

### Phase 2: Native PQ Transactions

- Type `0x50` transactions recognized by validators
- PQ accounts are first-class citizens
- SDK auto-creates PQ accounts for new users

### Phase 3: Validator Migration

- Validator migration toolkit (CLI) for rotating to PQ signing keys
- Existing ECDSA accounts can migrate via key rotation contract
- Governance proposal for enforcing PQ-only blocks after threshold date

### Backward Compatibility

- ECDSA transactions continue to work unchanged
- PQ accounts can interact with ECDSA accounts normally
- Existing contracts require no modifications
- Social login flow in Beam SDK is unchanged — PQ is handled at the infrastructure level

## 7. Security Considerations

### Algorithm Choice

- ML-DSA-65 is NIST's primary recommendation for general-purpose digital signatures
- SLH-DSA-128s provides a conservative fallback with minimal cryptographic assumptions (hash-based only)
- Both algorithms are standardized in NIST FIPS 204/205

### Key Sizes

ML-DSA-65 public keys (1,952 bytes) are significantly larger than ECDSA (33/65 bytes). This affects:
- Transaction size and calldata costs
- State storage for smart contract wallets
- Block size (mitigated by Beam's 4,500 TPS throughput)

### Side-Channel Resistance

The liboqs implementation used in the precompile is designed for constant-time operation. However, CGo boundary may introduce timing variations. The security audit (M3) will specifically examine this.

### Signature Malleability

Unlike ECDSA (which required EIP-2 to enforce low-S values to prevent signature malleability), both ML-DSA-65 and SLH-DSA-128s are **non-malleable by design**:

- **ML-DSA-65 (FIPS 204):** The signing algorithm produces a canonical encoding. There is no degree of freedom in the signature representation — a given (message, secret key) pair produces a unique, deterministic signature. Padding bits are fixed by the spec.
- **SLH-DSA-128s (FIPS 205):** Uses deterministic signing with no flexibility in encoding. The hash-based construction produces exactly one valid signature per (message, key) pair.

This means transaction hash uniqueness is preserved without requiring additional strictness checks analogous to EIP-2. The precompile does not need to enforce canonical form — the algorithms guarantee it.

### CGo Dependency and Precedent

The use of CGo to call liboqs is an operational concern for validator builds. However, this approach has precedent in the EVM ecosystem:

- Ethereum's `bn256` pairing precompile used a C implementation (cloudflare/bn256) for years before pure Go alternatives matured.
- The `bls12-381` precompile (EIP-2537) references both C (blst) and Go implementations.

The long-term goal is a pure Go implementation of FIPS 204/205 to eliminate CGo entirely. In the interim, static linking of liboqs reduces operator burden and removes the runtime dependency on shared libraries. The security audit (M3) will specifically examine the CGo boundary for memory safety, side-channel leakage, and cross-platform determinism.

### Precompile Safety

- The precompile is stateless — no storage reads/writes
- It is a pure function: same inputs always produce same outputs
- No reentrancy risk
- Gas costs prevent DoS via expensive verification

## 8. Implementation Details

### Library: liboqs

- Version: 0.15.0+
- License: MIT
- Binding: CGo (direct C calls, no intermediate wrapper library)
- Thread safety: Each call creates and frees its own `OQS_SIG` context

### Subnet-EVM Integration

The precompile is registered in the Subnet-EVM precompile registry:
1. Added to `precompile/registry/registry.go`
2. Configured via genesis JSON (`pqVerifyConfig`)
3. Activatable at specific block timestamp

### Dependencies

- liboqs 0.15.0 (MIT) — PQ cryptographic primitives
- OpenSSL 3.x — required by liboqs for symmetric crypto (AES, SHA)
- golang.org/x/crypto — keccak256 for address derivation

## 9. References

- [NIST FIPS 204 (ML-DSA)](https://csrc.nist.gov/pubs/fips/204/final)
- [NIST FIPS 205 (SLH-DSA)](https://csrc.nist.gov/pubs/fips/205/final)
- [Open Quantum Safe liboqs](https://github.com/open-quantum-safe/liboqs)
- [Beam SDK Documentation](https://docs.onbeam.com/sdk)
- [Subnet-EVM Custom Precompiles](https://docs.avax.network/build/subnet/upgrade/customize-a-subnet#precompiles)
- [EIP-2718: Typed Transaction Envelope](https://eips.ethereum.org/EIPS/eip-2718)
