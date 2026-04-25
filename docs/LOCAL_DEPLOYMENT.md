# Local Deployment Guide

How to build, deploy, and run the PQ_VERIFY precompile on a local Avalanche subnet. Covers the version constraints, known pitfalls, and the path toward production on Beam.

---

## Prerequisites

| Tool | Version | Install |
|------|---------|---------|
| Go | 1.23+ | `brew install go` or [golang.org](https://go.dev/dl/) |
| liboqs | 0.15.0+ | `brew install liboqs` (macOS) or [build from source](#linux-liboqs) |
| OpenSSL | 3.x | Usually pre-installed on macOS/Linux |
| Avalanche CLI | 1.9.x | `curl -sSfL https://raw.githubusercontent.com/ava-labs/avalanche-cli/main/scripts/install.sh \| sh -s` |
| Foundry | latest | `curl -L https://foundry.paradigm.xyz \| bash && foundryup` |

### Linux liboqs

```bash
sudo apt-get install -y cmake ninja-build libssl-dev
git clone --depth 1 --branch 0.15.0 https://github.com/open-quantum-safe/liboqs.git
cd liboqs && mkdir build && cd build
cmake -GNinja -DCMAKE_INSTALL_PREFIX=/usr/local ..
ninja && sudo ninja install && sudo ldconfig
```

---

## Quick Deploy (Automated)

```bash
./scripts/deploy_local.sh
```

This script handles everything: cloning subnet-evm, patching it, building both the custom VM and a matching AvalancheGo binary, creating the subnet, deploying it, and deploying the test contract.

If you prefer to understand each step or need to debug, read on.

---

## Manual Step-by-Step

### 1. Clone and patch subnet-evm

```bash
# Clone the archived subnet-evm repo (outside this project directory)
git clone --depth 1 https://github.com/ava-labs/subnet-evm.git ../subnet-evm
cd ../subnet-evm

# Copy pqcrypto CGo bindings
mkdir -p precompile/contracts/pqcrypto
cp ../pq-beam-verify-precompile/pkg/pqcrypto/pqcrypto.go precompile/contracts/pqcrypto/

# The deploy_local.sh script generates the adapted precompile files
# (contract.go, config.go, module.go) into precompile/contracts/pqverify/.
# See scripts/deploy_local.sh Step 2 for the exact file contents.

# Patch the registry
# In precompile/registry/registry.go, add:
#   _ "github.com/ava-labs/subnet-evm/precompile/contracts/pqverify"
```

### 2. Build the custom Subnet-EVM binary

```bash
cd ../subnet-evm

CGO_ENABLED=1 \
CGO_CFLAGS="-I/opt/homebrew/include" \
CGO_LDFLAGS="-L/opt/homebrew/lib -loqs -lcrypto" \
go build -o build/subnet-evm ./plugin/
```

**CGO_ENABLED=1 is mandatory.** The pqcrypto package uses CGo to call liboqs. A build with `CGO_ENABLED=0` will fail at compile time due to `import "C"`.

For Linux, adjust paths:
```bash
CGO_CFLAGS="-I/usr/local/include"
CGO_LDFLAGS="-L/usr/local/lib -loqs -lcrypto"
```

### 3. Build a matching AvalancheGo binary

> **This is critical. Skip this and deployment will fail.**

The archived subnet-evm master depends on `avalanchego v1.14.1-antithesis-docker-image-fix`. The Avalanche CLI ships AvalancheGo v1.14.0. These are **not compatible** -- see [Version Mismatch](#version-mismatch-the-critical-gotcha) below.

```bash
git clone --depth 1 --branch v1.14.1-antithesis-docker-image-fix \
    https://github.com/ava-labs/avalanchego.git /tmp/avalanchego-build
cd /tmp/avalanchego-build && ./scripts/build.sh

# Replace the CLI's AvalancheGo binary
cp /tmp/avalanchego-build/build/avalanchego \
    ~/.avalanche-cli/bin/avalanchego/avalanchego-v1.14.0/avalanchego

# Clean up
rm -rf /tmp/avalanchego-build
```

Verify both binaries report the same protocol:
```bash
~/.avalanche-cli/bin/avalanchego/avalanchego-v1.14.0/avalanchego --version
# avalanchego/1.14.0 [..., rpcchainvm=44, ...]

../subnet-evm/build/subnet-evm --version
# Subnet-EVM/v0.8.0 [AvalancheGo=avalanchego/1.14.0, rpcchainvm=44]
```

Both must show **rpcchainvm=44**.

### 4. Create and deploy the blockchain

```bash
# Delete any previous attempt
echo "y" | ~/bin/avalanche blockchain delete vanillatest 2>/dev/null

# Set up the blockchain config manually (non-interactive)
mkdir -p ~/.avalanche-cli/subnets/vanillatest ~/.avalanche-cli/vms
cp scripts/genesis.json ~/.avalanche-cli/subnets/vanillatest/genesis.json
cp ../subnet-evm/build/subnet-evm ~/.avalanche-cli/vms/vanillatest
chmod +x ~/.avalanche-cli/vms/vanillatest

# Create sidecar.json (see scripts/deploy_local.sh for the interactive alternative)
cat > ~/.avalanche-cli/subnets/vanillatest/sidecar.json << 'EOF'
{
    "Name": "vanillatest",
    "VM": "Custom",
    "RPCVersion": 44,
    "Subnet": "vanillatest",
    "TokenName": "PQ Token",
    "TokenSymbol": "PQ",
    "Version": "1.4.0",
    "Networks": {},
    "Sovereign": false
}
EOF

# Deploy with lowered disk space requirement (default 10GiB is too high for most dev laptops)
cat > /tmp/avago-wrapper.sh << 'WRAPPER'
#!/bin/bash
exec ~/.avalanche-cli/bin/avalanchego/avalanchego-v1.14.0/avalanchego \
    --system-tracker-disk-required-available-space=1073741824 "$@"
WRAPPER
chmod +x /tmp/avago-wrapper.sh

~/bin/avalanche blockchain deploy vanillatest --local --ewoq \
    --avalanchego-path /tmp/avago-wrapper.sh
```

### 5. Verify it works

```bash
# Get the RPC URL from the deployment output, then:
RPC_URL="http://127.0.0.1:9650/ext/bc/<blockchain-id>/rpc"

# Chain is alive?
cast chain-id --rpc-url "$RPC_URL"
# 13337

# Precompile is registered?
cast code 0x0300000000000000000000000000000000000000 --rpc-url "$RPC_URL"
# 0x01
```

### 6. Stop / restart

```bash
~/bin/avalanche network stop
~/bin/avalanche network start   # resumes from snapshot
```

---

## Version Mismatch: The Critical Gotcha

This is the single most important thing to understand when working with this codebase.

### The problem

The `subnet-evm` repo has been **archived** and moved into [AvalancheGo's graft directory](https://github.com/ava-labs/avalanchego/tree/master/graft/subnet-evm). The archived master branch's `go.mod` depends on:

```
github.com/ava-labs/avalanchego v1.14.1-antithesis-docker-image-fix
```

This pre-release tag added the `HeliconTime` field to the protobuf `NetworkUpgrades` message (`vm.proto` field 18). But the Avalanche CLI (v1.9.6) downloads **AvalancheGo v1.14.0**, which does **not** include `HeliconTime`.

### What happens

1. AvalancheGo v1.14.0 sends `InitializeRequest` to the VM plugin via gRPC
2. The `NetworkUpgrades` protobuf message omits `helicon_time` (field 18)
3. The VM binary (compiled against v1.14.1) deserializes it as `nil`
4. `vm_server.go:910` calls `grpcutils.TimestampAsTime(nil)` -- crashes with:
   ```
   proto: invalid nil Timestamp
   ```

### What it looks like

The blockchain "hangs" during bootstrap. The chain-specific log file has only 3 lines (handshake + gRPC setup). The real error is in `main.log`:

```
ERROR chains/manager.go:403 error creating chain
    "error": "error while creating new snowman vm rpc error: code = Unknown
              desc = invalid timestamp: proto: invalid nil Timestamp"
```

### The fix

Build AvalancheGo from the matching tag (`v1.14.1-antithesis-docker-image-fix`). Both binaries must use **rpcchainvm protocol 44**.

### Why not just upgrade to AvalancheGo v1.14.2?

AvalancheGo v1.14.2 bumped the protocol to **rpcchainvm=45**. The archived subnet-evm code is compiled against protocol 44 and has API incompatibilities with v1.14.2 (warp, firewood interfaces changed). You cannot simply bump `go.mod` -- the code won't compile.

### Version compatibility matrix

| AvalancheGo | rpcchainvm | HeliconTime | Works with archived subnet-evm? |
|-------------|-----------|-------------|--------------------------------|
| v1.14.0 (CLI default) | 44 | No | **No** -- nil timestamp crash |
| v1.14.1-antithesis (pre-release) | 44 | Yes | **Yes** -- must build from source |
| v1.14.2 | 45 | Yes | **No** -- protocol mismatch + API changes |

---

## Disk Space

AvalancheGo defaults to requiring **10 GiB** free disk space (`--system-tracker-disk-required-available-space`). Most dev laptops don't have this much free. If you see:

```
FATAL node/node.go:1464 low on disk space. Shutting down...
```

Lower the threshold with `--system-tracker-disk-required-available-space=1073741824` (1 GiB). The deploy script does this automatically via a wrapper.

---

## Precompile Address

The precompile is deployed at:

```
0x0300000000000000000000000000000000000000
```

**Note:** The README and some Solidity files reference `0x0b00`. The actual deployed address in the Subnet-EVM module is `0x0300...`. This is set in `subnet-evm/precompile/contracts/pqverify/module.go`. Update any client code accordingly.

---

## Genesis Configuration

The `pqVerifyConfig` section in `scripts/genesis.json` controls precompile activation:

```json
{
  "config": {
    "chainId": 13337,
    "feeConfig": { ... },
    "pqVerifyConfig": {
      "blockTimestamp": 0
    }
  }
}
```

- `"blockTimestamp": 0` -- active from genesis
- `"blockTimestamp": 1700000000` -- activate at that Unix timestamp
- Omit `pqVerifyConfig` entirely to deploy without the precompile (useful for debugging)

---

## Moving to Beam's Architecture (Future Steps)

When integrating PQ_VERIFY into Beam's production subnet, the approach changes significantly.

### What changes

1. **Source location**: Beam maintains its own Subnet-EVM fork (not the archived repo). The precompile code goes into Beam's fork directly.

2. **No version mismatch**: Beam's fork pins its own AvalancheGo version. The precompile is compiled as part of that fork, so the protobuf protocol always matches.

3. **Precompile address**: Must be coordinated with the Beam team. The address `0x0300...` is arbitrary for local dev. Production needs an address that doesn't conflict with existing Beam precompiles.

4. **CGo dependency on validators**: Every Beam validator must have `liboqs` installed. This is a non-trivial operational requirement:
   - Docker images need liboqs baked in
   - Different architectures (ARM vs x86) produce different verification timings
   - Static linking (`-static`) is strongly preferred for production to avoid shared library version drift

5. **Gas calibration**: The current gas costs use a 10x safety margin based on M1 Pro benchmarks. These must be re-calibrated on actual Beam validator hardware (likely x86 cloud instances).

### Steps for Beam integration

1. **Fork Beam's Subnet-EVM** (not the archived repo)
2. **Copy** `precompile/contracts/pqverify/` and `precompile/contracts/pqcrypto/` into the fork
3. **Add the registry import** in `precompile/registry/registry.go`
4. **Build with CGo** and liboqs linked
5. **Run Beam's existing test suite** to ensure no regressions
6. **Deploy to Beam's testnet** with `pqVerifyConfig` at a future timestamp
7. **Monitor** gas usage, verification times, and validator resource consumption
8. **Activate on mainnet** via governance once testnet is stable

### Things that will break or need attention

| Area | Risk | Mitigation |
|------|------|------------|
| liboqs on all validators | Validators without liboqs crash on startup | Static link liboqs; add to validator Docker image |
| Gas costs | 10x margin may be too high or too low on different hardware | Benchmark on Beam validator hardware; make costs configurable via genesis |
| Precompile address conflict | `0x0300...` may conflict with existing Beam precompiles | Coordinate address with Beam team before deployment |
| Subnet-EVM version drift | If Beam upgrades Subnet-EVM, the precompile interface may change | Pin to Subnet-EVM's `StatefulPrecompiledContract` interface; it's stable |
| ABI encoding | The precompile uses raw ABI encoding, not Solidity function selectors | Any future Solidity interface changes need matching precompile updates |
| SLH-DSA algorithm name | liboqs changed `"SLH-DSA-SHAKE-128s"` to `"SLH_DSA_PURE_SHA2_128S"` across versions | Pin liboqs version; test with the exact version deployed on validators |
| Key/signature sizes | ML-DSA-65 signatures are 3,309 bytes -- larger than standard EVM transaction data | Ensure gas limits and calldata pricing account for large inputs |
| Upgrade path | Activating at genesis is irreversible; bugs can't be patched by deactivating | Deploy with a future `blockTimestamp`; test thoroughly on testnet first |

### Production checklist

- [ ] liboqs statically linked in the validator binary
- [ ] Gas costs calibrated on Beam validator hardware
- [ ] Precompile address finalized and documented
- [ ] Security audit of CGo boundary (Phase 3)
- [ ] Fuzzing of ABI decoder with malformed inputs
- [ ] Testnet deployment with monitoring for 2+ weeks
- [ ] Governance proposal for mainnet activation timestamp

---

## Troubleshooting

### "proto: invalid nil Timestamp"
You are running AvalancheGo v1.14.0 with a subnet-evm binary compiled against v1.14.1. See [Version Mismatch](#version-mismatch-the-critical-gotcha).

### "low on disk space. Shutting down..."
Lower the threshold: `--system-tracker-disk-required-available-space=1073741824`. See [Disk Space](#disk-space).

### Chain hangs during bootstrap (no error visible)
The error is hidden. Check **main.log** (not the chain-specific log):
```bash
grep "error creating chain" ~/.avalanche-cli/runs/network_*/NodeID-*/logs/main.log
```

### "context deadline exceeded" during deploy
Same as above -- the chain failed to bootstrap within the timeout. Check main.log for the actual error.

### VM binary has no precompile
Verify with: `nm <binary> | grep pqverify`. If no output, the binary was built without the precompile patches. Rebuild from the patched source with `CGO_ENABLED=1`.

### Cast returns empty/errors on precompile call
Check that you're calling the correct address (`0x0300...`, not `0x0b00`). Verify the precompile is registered: `cast code 0x0300000000000000000000000000000000000000 --rpc-url $RPC_URL` should return `0x01`.
