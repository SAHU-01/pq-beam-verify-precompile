# Subnet-EVM Integration Guide

Step-by-step instructions for integrating PQ_VERIFY into a Subnet-EVM fork.

> **Important:** The `subnet-evm` repo has been [archived and moved into AvalancheGo](https://github.com/ava-labs/avalanchego/tree/master/graft/subnet-evm). The archived master depends on `avalanchego v1.14.1-antithesis-docker-image-fix`, which requires building AvalancheGo from source to get a matching binary. See [`docs/LOCAL_DEPLOYMENT.md`](../docs/LOCAL_DEPLOYMENT.md) for the full version compatibility story.

## Prerequisites

- Go 1.23+
- liboqs 0.15.0+ installed system-wide
- Avalanche CLI (`avalanche blockchain create/deploy`)
- Foundry (for contract testing: `cast`, `forge`)

## Step 1: Fork Subnet-EVM

```bash
git clone https://github.com/ava-labs/subnet-evm.git
cd subnet-evm
git checkout -b pq-verify-integration
```

## Step 2: Add PQ_VERIFY precompile

```bash
# Copy the precompile package
mkdir -p precompile/contracts/pqverify
cp -r <this-repo>/pkg/pqverify/* precompile/contracts/pqverify/
cp -r <this-repo>/pkg/pqcrypto/ precompile/contracts/pqcrypto/
cp <this-repo>/subnet-evm/module.go precompile/contracts/pqverify/module.go
```

Update the module.go imports to use Subnet-EVM's internal paths:
```go
import (
    "github.com/ava-labs/subnet-evm/precompile/contract"
    "github.com/ava-labs/subnet-evm/precompile/modules"
)
```

## Step 3: Register in the precompile registry

Edit `precompile/registry/registry.go` and add the import:

```go
import (
    // existing imports...
    _ "github.com/ava-labs/subnet-evm/precompile/contracts/pqverify"
)
```

Create `precompile/contracts/pqverify/init.go`:

```go
package pqverify

import "github.com/ava-labs/subnet-evm/precompile/modules"

func init() {
    modules.RegisterModule(Module)
}
```

## Step 4: Add CGo flags to the build

In the repository root `Makefile` or build script, ensure CGo is enabled:

```bash
export CGO_ENABLED=1
export CGO_CFLAGS="-I/opt/homebrew/include"    # macOS
export CGO_LDFLAGS="-L/opt/homebrew/lib -loqs -lcrypto"
```

For Linux:
```bash
export CGO_CFLAGS="-I/usr/local/include"
export CGO_LDFLAGS="-L/usr/local/lib -loqs -lcrypto"
```

## Step 5: Build the modified binary

```bash
CGO_ENABLED=1 go build -o build/subnet-evm ./plugin/
```

## Step 5b: Build matching AvalancheGo

The Avalanche CLI's default AvalancheGo v1.14.0 is **incompatible** with the archived subnet-evm (missing `HeliconTime` in protobuf). You must build from the matching tag:

```bash
git clone --depth 1 --branch v1.14.1-antithesis-docker-image-fix \
    https://github.com/ava-labs/avalanchego.git /tmp/avago-build
cd /tmp/avago-build && ./scripts/build.sh
cp /tmp/avago-build/build/avalanchego ~/.avalanche-cli/bin/avalanchego/avalanchego-v1.14.0/avalanchego
```

See [`docs/LOCAL_DEPLOYMENT.md`](../docs/LOCAL_DEPLOYMENT.md) for why this is necessary.

## Step 6: Create and deploy the subnet

```bash
# Create subnet with PQ genesis config
avalanche blockchain create pq-testnet \
    --custom \
    --genesis <this-repo>/scripts/genesis.json \
    --custom-vm-path build/subnet-evm

# Deploy locally (use wrapper to lower disk space requirement)
avalanche blockchain deploy pq-testnet --local --ewoq
```

## Step 7: Verify the precompile

```bash
# Deploy the test helper contract
forge create contracts/PQVerifyTestHelper.sol:PQVerifyTestHelper \
    --rpc-url <subnet-rpc-url> \
    --private-key <deployer-key>

# Generate a test keypair and verify
go run ./cmd/keygen/ -- --algorithm ml-dsa-65 --verify
```

## Genesis Configuration

The `pqVerifyConfig` section in genesis.json controls precompile activation:

```json
{
  "pqVerifyConfig": {
    "blockTimestamp": 0,
    "gasOverrides": {
      "mlDsa65Gas": 130000,
      "slhDsa128sGas": 520000,
      "baseGas": 3600
    }
  }
}
```

Set `blockTimestamp` to a future timestamp to activate the precompile at a specific time instead of genesis.
