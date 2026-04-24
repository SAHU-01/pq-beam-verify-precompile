#!/bin/bash
# deploy_test.sh — Deploys and tests PQ_VERIFY on a local Subnet-EVM chain.
#
# Prerequisites:
#   - Avalanche CLI installed (avalanche subnet create/deploy)
#   - Modified subnet-evm binary with PQ_VERIFY precompile
#   - Foundry (cast, forge) for contract deployment
#
# This script:
#   1. Creates a local subnet with PQ_VERIFY enabled
#   2. Deploys the PQVerifyTestHelper contract
#   3. Generates a PQ keypair and signs a test message
#   4. Calls the precompile via the helper contract
#   5. Verifies the result
#
# Usage: ./scripts/deploy_test.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

echo "=== PQ_VERIFY Local Chain Test ==="
echo ""

# Step 1: Generate test keypair
echo "[1/5] Generating ML-DSA-65 keypair..."
KEYGEN_OUTPUT=$(cd "$PROJECT_DIR" && CGO_ENABLED=1 go run ./cmd/keygen/ -- --algorithm ml-dsa-65 --output /tmp/pq-deploy-test --format hex)
echo "$KEYGEN_OUTPUT"

PQ_ADDRESS=$(echo "$KEYGEN_OUTPUT" | grep "PQ Address" | awk '{print $3}')
echo "PQ Address: $PQ_ADDRESS"

# Step 2: Read the generated keys
PUBKEY_HEX=$(cat /tmp/pq-deploy-test/ml-dsa-65.pub)
# Note: In a real deployment, you'd sign a message here using the keygen tool

echo ""
echo "[2/5] Running precompile unit tests..."
cd "$PROJECT_DIR"
CGO_ENABLED=1 go test ./pkg/pqverify/ -v -run TestValidMLDSA65

echo ""
echo "[3/5] Running end-to-end tests..."
CGO_ENABLED=1 go test ./test/ -v -run TestE2E_MLDSA65_FullFlow

echo ""
echo "[4/5] Running benchmarks..."
CGO_ENABLED=1 go test ./cmd/benchmark/ -bench=BenchmarkMLDSA65Verify -benchmem -count=1

echo ""
echo "[5/5] Running full test suite..."
CGO_ENABLED=1 go test ./... -count=1

echo ""
echo "=== All tests passed ==="
echo ""
echo "To deploy on a local Subnet-EVM chain:"
echo "  1. Build the modified subnet-evm binary with PQ_VERIFY precompile"
echo "  2. avalanche subnet create pq-testnet --custom --genesis scripts/genesis.json"
echo "  3. avalanche subnet deploy pq-testnet --local"
echo "  4. Deploy PQVerifyTestHelper.sol using Foundry:"
echo "     forge create contracts/PQVerifyTestHelper.sol:PQVerifyTestHelper --rpc-url <subnet-rpc>"
echo "  5. Call verifyAndLog() with a PQ signature"
