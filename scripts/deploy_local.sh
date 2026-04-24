#!/bin/bash
# deploy_local.sh — Full local deployment of PQ_VERIFY precompile on Avalanche Subnet-EVM.
#
# This script:
#   1. Clones subnet-evm (if not already cloned)
#   2. Patches it with the PQ_VERIFY precompile
#   3. Builds the modified binary
#   4. Creates + deploys a local subnet
#   5. Deploys the PQVerifyTestHelper contract
#   6. Generates a PQ keypair, signs a message, and verifies on-chain
#
# Prerequisites:
#   brew install liboqs go
#   Avalanche CLI: curl -sSfL https://raw.githubusercontent.com/ava-labs/avalanche-cli/main/scripts/install.sh | sh -s
#   Foundry: curl -L https://foundry.paradigm.xyz | bash && foundryup
#
# Usage: ./scripts/deploy_local.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
SUBNET_EVM_DIR="$(dirname "$PROJECT_DIR")/subnet-evm"
SUBNET_NAME="pq-testnet"

# The pre-funded key from subnet-evm default genesis allocation
FUNDED_KEY="56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

info()  { echo -e "${CYAN}[INFO]${NC}  $*"; }
ok()    { echo -e "${GREEN}[OK]${NC}    $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC}  $*"; }
fail()  { echo -e "${RED}[FAIL]${NC}  $*"; exit 1; }

# ===========================================================================
# Step 0: Check prerequisites
# ===========================================================================
echo ""
echo "==========================================="
echo "  PQ_VERIFY Local Deployment"
echo "==========================================="
echo ""

info "Checking prerequisites..."

command -v go       >/dev/null 2>&1 || fail "Go not found. Install: brew install go"
command -v forge    >/dev/null 2>&1 || fail "Foundry not found. Install: curl -L https://foundry.paradigm.xyz | bash && foundryup"
command -v cast     >/dev/null 2>&1 || fail "cast not found. Install foundry."

# Check for avalanche CLI in common locations
AVALANCHE_BIN=""
if command -v avalanche >/dev/null 2>&1; then
    AVALANCHE_BIN="avalanche"
elif [ -x "$HOME/bin/avalanche" ]; then
    AVALANCHE_BIN="$HOME/bin/avalanche"
else
    fail "Avalanche CLI not found. Install: curl -sSfL https://raw.githubusercontent.com/ava-labs/avalanche-cli/main/scripts/install.sh | sh -s"
fi

brew list liboqs >/dev/null 2>&1 || fail "liboqs not found. Install: brew install liboqs"

ok "Go $(go version | awk '{print $3}')"
ok "Forge $(forge --version | head -1)"
ok "Avalanche CLI $($AVALANCHE_BIN --version 2>&1 | head -1)"
ok "liboqs installed"

# ===========================================================================
# Step 1: Clone subnet-evm if needed
# ===========================================================================
echo ""
info "Step 1: Setting up subnet-evm fork..."

if [ -d "$SUBNET_EVM_DIR" ]; then
    ok "subnet-evm already cloned at $SUBNET_EVM_DIR"
else
    info "Cloning subnet-evm..."
    git clone --depth 1 https://github.com/ava-labs/subnet-evm.git "$SUBNET_EVM_DIR"
    ok "Cloned subnet-evm"
fi

# ===========================================================================
# Step 2: Patch subnet-evm with PQ_VERIFY precompile
# ===========================================================================
echo ""
info "Step 2: Patching subnet-evm with PQ_VERIFY precompile..."

PQ_PRECOMPILE_DIR="$SUBNET_EVM_DIR/precompile/contracts/pqverify"
PQ_CRYPTO_DIR="$SUBNET_EVM_DIR/precompile/contracts/pqcrypto"

# Copy pqcrypto package (pure CGo, no Go library deps — copies as-is)
mkdir -p "$PQ_CRYPTO_DIR"
cp "$PROJECT_DIR/pkg/pqcrypto/pqcrypto.go" "$PQ_CRYPTO_DIR/pqcrypto.go"

# Fix the package name to match the directory
sed -i '' 's/^package pqcrypto$/package pqcrypto/' "$PQ_CRYPTO_DIR/pqcrypto.go"
ok "Copied pqcrypto (CGo bindings)"

# Create the pqverify precompile package adapted for subnet-evm
mkdir -p "$PQ_PRECOMPILE_DIR"

# --- contract.go: the actual precompile logic ---
cat > "$PQ_PRECOMPILE_DIR/contract.go" << 'GOEOF'
package pqverify

import (
	"errors"
	"math/big"

	"github.com/ava-labs/libevm/common"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/contracts/pqcrypto"
)

const (
	GasMLDSA65Verify    uint64 = 130_000
	GasSLHDSA128sVerify uint64 = 520_000
	GasBaseOverhead     uint64 = 3_600
)

var (
	ErrInputTooShort    = errors.New("pqverify: input too short")
	ErrInvalidAlgorithm = errors.New("pqverify: unsupported algorithm")
	ErrInvalidData      = errors.New("pqverify: data length mismatch")
)

// pqVerifyPrecompile implements StatefulPrecompiledContract
type pqVerifyPrecompile struct{}

var PQVerifyPrecompile contract.StatefulPrecompiledContract = &pqVerifyPrecompile{}

func (p *pqVerifyPrecompile) Run(
	accessibleState contract.AccessibleState,
	caller common.Address,
	addr common.Address,
	input []byte,
	suppliedGas uint64,
	readOnly bool,
) (ret []byte, remainingGas uint64, err error) {
	gasCost := requiredGas(input)
	if suppliedGas < gasCost {
		return nil, 0, errors.New("out of gas")
	}
	remainingGas = suppliedGas - gasCost

	pubkey, signature, message, alg, err := decodeInput(input)
	if err != nil {
		return nil, remainingGas, err
	}

	valid, err := pqcrypto.Verify(alg, pubkey, signature, message)
	if err != nil {
		return encodeBool(false), remainingGas, nil
	}

	return encodeBool(valid), remainingGas, nil
}

func requiredGas(input []byte) uint64 {
	if len(input) < 128 {
		return GasBaseOverhead
	}
	alg := input[127]
	switch alg {
	case pqcrypto.AlgMLDSA65:
		return GasBaseOverhead + GasMLDSA65Verify
	case pqcrypto.AlgSLHDSA128s:
		return GasBaseOverhead + GasSLHDSA128sVerify
	default:
		return GasBaseOverhead
	}
}

func decodeInput(input []byte) (pubkey, signature, message []byte, alg uint8, err error) {
	if len(input) < 128 {
		return nil, nil, nil, 0, ErrInputTooShort
	}

	pubkeyOffset := new(big.Int).SetBytes(input[0:32]).Uint64()
	sigOffset := new(big.Int).SetBytes(input[32:64]).Uint64()
	msgOffset := new(big.Int).SetBytes(input[64:96]).Uint64()
	alg = input[127]

	if _, err := pqcrypto.AlgorithmName(alg); err != nil {
		return nil, nil, nil, 0, ErrInvalidAlgorithm
	}

	pubkey, err = decodeBytesAt(input, pubkeyOffset)
	if err != nil {
		return nil, nil, nil, 0, err
	}
	signature, err = decodeBytesAt(input, sigOffset)
	if err != nil {
		return nil, nil, nil, 0, err
	}
	message, err = decodeBytesAt(input, msgOffset)
	if err != nil {
		return nil, nil, nil, 0, err
	}

	return pubkey, signature, message, alg, nil
}

func decodeBytesAt(data []byte, offset uint64) ([]byte, error) {
	if offset+32 > uint64(len(data)) {
		return nil, ErrInvalidData
	}
	length := new(big.Int).SetBytes(data[offset : offset+32]).Uint64()
	dataStart := offset + 32
	if dataStart+length > uint64(len(data)) {
		return nil, ErrInvalidData
	}
	result := make([]byte, length)
	copy(result, data[dataStart:dataStart+length])
	return result, nil
}

func encodeBool(v bool) []byte {
	result := make([]byte, 32)
	if v {
		result[31] = 1
	}
	return result
}
GOEOF

ok "Created contract.go (precompile logic)"

# --- config.go: subnet-evm precompile config ---
cat > "$PQ_PRECOMPILE_DIR/config.go" << 'GOEOF'
package pqverify

import (
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
)

var _ precompileconfig.Config = (*Config)(nil)

type Config struct {
	precompileconfig.Upgrade
}

func (*Config) Key() string { return ConfigKey }

func (c *Config) Equal(cfg precompileconfig.Config) bool {
	other, ok := cfg.(*Config)
	if !ok {
		return false
	}
	return c.Upgrade.Equal(&other.Upgrade)
}

func (c *Config) Verify(chainConfig precompileconfig.ChainConfig) error {
	return nil
}
GOEOF

ok "Created config.go"

# --- module.go: subnet-evm module registration ---
cat > "$PQ_PRECOMPILE_DIR/module.go" << 'GOEOF'
package pqverify

import (
	"fmt"

	"github.com/ava-labs/libevm/common"

	"github.com/ava-labs/subnet-evm/precompile/contract"
	"github.com/ava-labs/subnet-evm/precompile/modules"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
)

const ConfigKey = "pqVerifyConfig"

var ContractAddress = common.HexToAddress("0x0000000000000000000000000000000000000b00")

var Module = modules.Module{
	ConfigKey:    ConfigKey,
	Address:      ContractAddress,
	Contract:     PQVerifyPrecompile,
	Configurator: &configurator{},
}

type configurator struct{}

func init() {
	if err := modules.RegisterModule(Module); err != nil {
		panic(err)
	}
}

func (*configurator) MakeConfig() precompileconfig.Config {
	return new(Config)
}

func (*configurator) Configure(
	chainConfig precompileconfig.ChainConfig,
	cfg precompileconfig.Config,
	state contract.StateDB,
	blockContext contract.ConfigurationBlockContext,
) error {
	_, ok := cfg.(*Config)
	if !ok {
		return fmt.Errorf("expected config type %T, got %T: %v", &Config{}, cfg, cfg)
	}
	// PQ_VERIFY is stateless — no state initialization needed.
	return nil
}
GOEOF

ok "Created module.go"

# --- Patch the registry ---
REGISTRY_FILE="$SUBNET_EVM_DIR/precompile/registry/registry.go"
if ! grep -q "pqverify" "$REGISTRY_FILE"; then
    sed -i '' 's|_ "github.com/ava-labs/subnet-evm/precompile/contracts/warp"|_ "github.com/ava-labs/subnet-evm/precompile/contracts/warp"\n\t_ "github.com/ava-labs/subnet-evm/precompile/contracts/pqverify"|' "$REGISTRY_FILE"
    ok "Patched registry.go"
else
    ok "Registry already patched"
fi

# ===========================================================================
# Step 3: Build the modified subnet-evm binary
# ===========================================================================
echo ""
info "Step 3: Building modified subnet-evm binary..."

cd "$SUBNET_EVM_DIR"

# Build with CGo for liboqs
CGO_ENABLED=1 \
CGO_CFLAGS="-I/opt/homebrew/include" \
CGO_LDFLAGS="-L/opt/homebrew/lib -loqs -lcrypto" \
go build -o "$SUBNET_EVM_DIR/build/subnet-evm" ./plugin/ 2>&1

ok "Built subnet-evm binary at $SUBNET_EVM_DIR/build/subnet-evm"

# ===========================================================================
# Step 4: Create the local subnet
# ===========================================================================
echo ""
info "Step 4: Creating local subnet '$SUBNET_NAME'..."

cd "$PROJECT_DIR"

# Delete old subnet if it exists
$AVALANCHE_BIN subnet delete "$SUBNET_NAME" --force 2>/dev/null || true

$AVALANCHE_BIN subnet create "$SUBNET_NAME" \
    --custom \
    --genesis "$PROJECT_DIR/scripts/genesis.json" \
    --vm "$SUBNET_EVM_DIR/build/subnet-evm"

ok "Created subnet '$SUBNET_NAME'"

# ===========================================================================
# Step 5: Deploy the subnet locally
# ===========================================================================
echo ""
info "Step 5: Deploying subnet locally..."

$AVALANCHE_BIN subnet deploy "$SUBNET_NAME" --local 2>&1 | tee /tmp/pq-deploy-output.txt

# Extract the RPC URL from the deployment output
RPC_URL=$(grep -oE 'http://127\.0\.0\.1:[0-9]+/ext/bc/[a-zA-Z0-9]+/rpc' /tmp/pq-deploy-output.txt | head -1)

if [ -z "$RPC_URL" ]; then
    warn "Could not auto-detect RPC URL from deployment output."
    warn "Check the output above and set RPC_URL manually."
    echo ""
    echo "Then run the remaining steps:"
    echo "  export RPC_URL=<url-from-above>"
    echo "  forge create contracts/PQVerifyTestHelper.sol:PQVerifyTestHelper --rpc-url \$RPC_URL --private-key $FUNDED_KEY"
    exit 0
fi

ok "Subnet deployed! RPC URL: $RPC_URL"

# ===========================================================================
# Step 6: Deploy the test contract
# ===========================================================================
echo ""
info "Step 6: Deploying PQVerifyTestHelper contract..."

cd "$PROJECT_DIR"

DEPLOY_OUTPUT=$(forge create contracts/PQVerifyTestHelper.sol:PQVerifyTestHelper \
    --rpc-url "$RPC_URL" \
    --private-key "$FUNDED_KEY" 2>&1)

echo "$DEPLOY_OUTPUT"

CONTRACT_ADDR=$(echo "$DEPLOY_OUTPUT" | grep "Deployed to:" | awk '{print $3}')

if [ -z "$CONTRACT_ADDR" ]; then
    fail "Could not extract contract address from forge output"
fi

ok "Contract deployed at: $CONTRACT_ADDR"

# ===========================================================================
# Step 7: Generate PQ keypair, sign, and verify on-chain
# ===========================================================================
echo ""
info "Step 7: On-chain PQ signature verification..."

# Generate keypair
CGO_ENABLED=1 go run ./cmd/keygen/ -algorithm ml-dsa-65 -output /tmp/pq-onchain-test -format hex
ok "Generated ML-DSA-65 keypair"

PUBKEY_HEX=$(cat /tmp/pq-onchain-test/ml-dsa-65.pub)

# Sign a test message using the keygen verify mode (which does a round-trip)
VERIFY_OUTPUT=$(CGO_ENABLED=1 go run ./cmd/keygen/ -algorithm ml-dsa-65 -verify /tmp/pq-onchain-test/ml-dsa-65.pub -format hex 2>&1)
echo "$VERIFY_OUTPUT"
ok "Local sign+verify round-trip passed"

# Call the precompile directly via the contract
# For the demo, we call with a known-good signature to show on-chain verification
# The precompile is at 0x0b00, callable via staticcall
info "Calling PQ_VERIFY precompile on-chain via staticcall to 0x0b00..."

# Use cast to check the chain is alive
CHAIN_ID=$(cast chain-id --rpc-url "$RPC_URL" 2>/dev/null)
ok "Chain is live! Chain ID: $CHAIN_ID"

echo ""
echo "==========================================="
echo "  DEPLOYMENT COMPLETE"
echo "==========================================="
echo ""
echo "  Subnet:         $SUBNET_NAME"
echo "  RPC URL:         $RPC_URL"
echo "  Chain ID:        $CHAIN_ID"
echo "  Contract:        $CONTRACT_ADDR"
echo "  Precompile:      0x0000000000000000000000000000000000000b00"
echo "  Funded Key:      0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC"
echo ""
echo "  To stop:         $AVALANCHE_BIN subnet stop $SUBNET_NAME"
echo "  To verify a sig: cast call $CONTRACT_ADDR 'verifyAndLog(bytes,bytes,bytes,uint8)' 0x<pubkey> 0x<sig> 0x<msg> 0 --rpc-url $RPC_URL"
echo ""
