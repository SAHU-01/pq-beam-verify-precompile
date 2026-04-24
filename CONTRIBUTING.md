# Contributing to PQ_VERIFY

## Development Setup

### macOS

```bash
brew install liboqs go
git clone https://github.com/pq-beam/verify-precompile.git
cd verify-precompile
CGO_ENABLED=1 go test ./... -v
```

### Linux

```bash
# Install liboqs from source (see README.md)
# Then:
CGO_ENABLED=1 \
  CGO_CFLAGS="-I/usr/local/include" \
  CGO_LDFLAGS="-L/usr/local/lib -loqs -lcrypto" \
  go test ./... -v
```

## Code Organization

- `pkg/pqcrypto/` — Low-level CGo bindings to liboqs. Changes here require understanding of C memory management and liboqs API.
- `pkg/pqverify/` — Precompile logic, ABI encoding, tx types, gas schedule. Most feature work happens here.
- `contracts/` — Solidity interfaces and example contracts. Compile with `solc ^0.8.20`.
- `cmd/` — CLI tools (benchmark, keygen). Self-contained entry points.
- `sdk/` — TypeScript SDK for Beam SDK integration (M2).
- `test/` — End-to-end tests that exercise the full stack.

## Testing

```bash
# Unit tests
CGO_ENABLED=1 go test ./pkg/... -v

# E2E tests
CGO_ENABLED=1 go test ./test/ -v

# Benchmarks
CGO_ENABLED=1 go test ./cmd/benchmark/ -bench=. -benchmem
```

## Adding a New PQ Algorithm

1. Add the algorithm ID constant in `pkg/pqcrypto/pqcrypto.go`
2. Add the OQS algorithm name mapping in `AlgorithmName()`
3. Add gas cost constant in `pkg/pqverify/precompile.go`
4. Add the algorithm case in `RequiredGas()`
5. Add tests in both `pkg/pqcrypto/` and `pkg/pqverify/`
6. Run benchmarks and set the gas cost based on results
7. Update `docs/TECHNICAL_SPEC.md`

## Security

- Never commit secret keys or test keys to the repo
- The precompile must remain stateless (no storage reads/writes)
- All CGo calls must properly free allocated memory
- Gas costs must prevent DoS — benchmark any changes

## For Auditors

The security-critical paths are:
1. `pkg/pqcrypto/pqcrypto.go` — CGo boundary, memory safety
2. `pkg/pqverify/precompile.go` — ABI decoding, gas accounting
3. `pkg/pqverify/txtype.go` — Transaction signing/verification
4. `contracts/PQAccount.sol` — Smart account authorization

## Contact

Open an issue or reach out to @asahudev on GitHub.
