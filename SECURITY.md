# Security Policy

## Scope

This project implements post-quantum cryptographic signature verification as a Subnet-EVM precompile. Security is critical — a vulnerability here could compromise on-chain signature verification for the entire Beam network.

### In Scope

- `pkg/pqcrypto/` — CGo bindings to liboqs (key generation, signing, verification)
- `pkg/pqverify/` — ABI decoder, gas metering, precompile logic
- `subnet-evm/` — Subnet-EVM module adapter
- `contracts/` — Solidity smart contracts (PQAccount, PQKeyRotation, IPQVerify)
- CGo memory management (buffer allocation, deallocation, pointer safety)
- ABI decoder handling of untrusted input
- Gas cost correctness (underpriced operations enable DoS)

### Out of Scope

- liboqs library internals (report to [Open Quantum Safe](https://github.com/open-quantum-safe/liboqs/security))
- Subnet-EVM / AvalancheGo vulnerabilities (report to [Ava Labs](https://github.com/ava-labs/subnet-evm/security))
- The landing page (`site/`)

## Reporting a Vulnerability

**Do not open a public GitHub issue for security vulnerabilities.**

Email: **ankita.sahu.dev@gmail.com**

Include:
- Description of the vulnerability
- Steps to reproduce
- Impact assessment
- Suggested fix (if any)

You will receive an acknowledgment within 48 hours. We aim to provide a fix or mitigation plan within 7 days for critical issues.

## Security Model

### Precompile Properties

| Property | Status |
|----------|--------|
| Stateless | Yes — no storage reads or writes |
| Read-only | Yes — called via `staticcall` only |
| Constant-time | Delegated to liboqs (designed for constant-time) |
| Gas-metered | Yes — proportional to computation with 10x safety margin |
| Reentrancy safe | Yes — no external calls, no state changes |
| Non-malleable | Yes — ML-DSA-65 and SLH-DSA-128s produce canonical, deterministic signatures (no EIP-2 equivalent needed) |

### Known Limitations

1. **CGo boundary timing**: The Go/C transition may introduce timing variations not present in pure-C liboqs. This is in scope for the Phase 3 security audit.

2. **Gas cost approximation**: Gas costs are benchmarked on Apple M1 Pro with a 10x safety margin. Production validators on different hardware may see different performance. Gas costs are configurable per-chain via genesis `gasOverrides`.

3. **liboqs version sensitivity**: Algorithm names changed between liboqs versions (e.g., `SLH-DSA-SHAKE-128s` to `SLH_DSA_PURE_SHA2_128S`). The precompile is tested against liboqs 0.15.0. Other versions may fail silently.

4. **Memory expansion costs**: SLH-DSA-128s signatures are ~8KB. The gas formula accounts for CPU compute time with a 10x safety margin but does not separately meter memory I/O for large signature buffers. This must be validated on production hardware.

5. **CGo in consensus path**: Injecting C code into a Go-based consensus client has operational risks (cross-compilation complexity, validator build fragility). This has precedent (Ethereum's bn256 precompile used C), but the long-term goal is a pure Go implementation of FIPS 204/205.

## Audit Status

| Phase | Status | Scope |
|-------|--------|-------|
| Phase 1 (current) | Fuzz tested, not audited | ABI decoder fuzz tests found and fixed overflow bug |
| Phase 3 (planned) | Third-party audit | CGo boundary, ABI decoder, gas accounting, side-channel analysis |

## Dependencies

| Dependency | Version | License | Security Contact |
|------------|---------|---------|-----------------|
| liboqs | 0.15.0 | MIT | [OQS Security](https://github.com/open-quantum-safe/liboqs/security) |
| go-ethereum | v1.17.2 | LGPL-3.0 | [Geth Security](https://github.com/ethereum/go-ethereum/security) |
| OpenSSL | 3.x | Apache-2.0 | [OpenSSL Security](https://www.openssl.org/policies/secpolicy.html) |
