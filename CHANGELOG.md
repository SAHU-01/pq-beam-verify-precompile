# Changelog

## [0.1.0] — 2026-04-25 (Phase 1: Precompile)

### Added

- **PQ_VERIFY precompile** at `0x0300000000000000000000000000000000000000`
  - ML-DSA-65 (FIPS 204) verification — 105us, 133,600 gas
  - SLH-DSA-128s (FIPS 205) verification — 432us, 523,600 gas
  - Stateless, gas-metered, staticcall-only
- **CGo bridge to liboqs 0.15** — direct C bindings for keygen, sign, verify
- **PQAccount.sol** — ERC-4337 compatible smart account with PQ signature verification
  - `validateUserOp()` for EntryPoint integration
  - `execute()` for direct PQ-signed transactions
  - Batch execution support
- **PQKeyRotation.sol** — key migration and rotation registry
  - ECDSA-to-PQ migration via `registerKey()`
  - PQ-to-PQ rotation with 24-hour timelock
  - Emergency revocation via ECDSA fallback
- **EIP-2718 Type 0x50** transaction envelope specification
- **Subnet-EVM module adapter** (`subnet-evm/module.go`)
- **Fuzz tests** for ABI decoder — `FuzzDecodeInput`, `FuzzPrecompileRun`
  - Found and fixed overflow bug in `decodeBytesAt()`
- **Test suite**: 36+ tests across 5 packages + 2 fuzz targets
- **Local subnet deployment** script with genesis configuration
- **On-chain demo** — valid + tampered signature verification on local subnet
- **Landing page** — bento-grid design with interactive expandable panels

### Security

- ABI decoder hardened with overflow-safe bounds checking (`safeUint64()`)
- All test data ephemeral — fresh keys generated per run
- Gas costs include 10x safety margin for validator hardware variance

### Infrastructure

- CI/CD via GitHub Actions (Go tests, benchmarks, lint, Solidity compilation)
- Automated local subnet deployment script
