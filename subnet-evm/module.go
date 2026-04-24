// Package pqverifyvm provides the Subnet-EVM precompile module adapter for
// PQ_VERIFY. This code is designed to be copied into a Subnet-EVM fork at:
//
//	precompile/contracts/pqverify/
//
// It implements the StatefulPrecompiledContract interface required by Subnet-EVM
// and delegates all cryptographic work to the pqverify package.
//
// Integration steps:
//  1. Copy this package into your Subnet-EVM fork under precompile/contracts/pqverify/
//  2. Add the import in precompile/registry/registry.go
//  3. Update genesis.json with pqVerifyConfig
//  4. Build the modified subnet-evm binary
//
// See INTEGRATION.md in this directory for full instructions.
package pqverifyvm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	// In the actual Subnet-EVM fork, these imports become:
	//   "github.com/ava-labs/subnet-evm/precompile/allowlist"
	//   "github.com/ava-labs/subnet-evm/precompile/contract"
	//   "github.com/ava-labs/subnet-evm/precompile/modules"
	// For now, we define the interfaces inline for compilation outside the fork.

	"github.com/pq-beam/verify-precompile/pkg/pqverify"
)

// ContractAddress is the precompile address for PQ_VERIFY.
var ContractAddress = common.HexToAddress("0x0000000000000000000000000000000000000b00")

// Module implements the Subnet-EVM precompile Module interface.
// When integrated into a Subnet-EVM fork, this struct is registered in
// precompile/registry/registry.go via an init() function.
var Module = &module{}

type module struct{}

// ConfigKey returns the genesis JSON key for this precompile.
// In genesis.json, this appears as "pqVerifyConfig": { ... }
func (m *module) ConfigKey() string {
	return "pqVerifyConfig"
}

// Address returns the precompile contract address.
func (m *module) Address() common.Address {
	return ContractAddress
}

// --- Subnet-EVM Contract Interface ---

// PQVerifyPrecompileAdapter wraps our pqverify.PQVerifyPrecompile to implement
// Subnet-EVM's StatefulPrecompiledContract interface.
//
// In Subnet-EVM, the interface is:
//
//	type StatefulPrecompiledContract interface {
//	    Run(accessibleState AccessibleState, caller common.Address,
//	        addr common.Address, input []byte, suppliedGas uint64,
//	        readOnly bool) (ret []byte, remainingGas uint64, err error)
//	}
type PQVerifyPrecompileAdapter struct {
	inner *pqverify.PQVerifyPrecompile
	cfg   *pqverify.Config
}

// NewPQVerifyPrecompile creates the adapter with optional gas overrides from
// the genesis config.
func NewPQVerifyPrecompile(cfg *pqverify.Config) *PQVerifyPrecompileAdapter {
	return &PQVerifyPrecompileAdapter{
		inner: &pqverify.PQVerifyPrecompile{},
		cfg:   cfg,
	}
}

// Run implements StatefulPrecompiledContract. It verifies gas, calls the
// underlying precompile, and returns remaining gas.
//
// This is the method Subnet-EVM calls when a transaction or staticcall
// targets address 0x0b00.
func (p *PQVerifyPrecompileAdapter) Run(
	caller common.Address,
	addr common.Address,
	input []byte,
	suppliedGas uint64,
	readOnly bool,
) (ret []byte, remainingGas uint64, err error) {

	// Calculate required gas (respecting any genesis overrides).
	gasCost := p.requiredGas(input)
	if suppliedGas < gasCost {
		return nil, 0, fmt.Errorf("out of gas: have %d, need %d", suppliedGas, gasCost)
	}
	remainingGas = suppliedGas - gasCost

	// PQ_VERIFY is stateless and read-only. It does not access state,
	// modify storage, or transfer value. Safe to call in any context.
	ret, err = p.inner.Run(input)
	if err != nil {
		// Malformed input — consume gas but return error.
		return nil, remainingGas, err
	}

	return ret, remainingGas, nil
}

// requiredGas returns the gas cost, applying genesis overrides if configured.
func (p *PQVerifyPrecompileAdapter) requiredGas(input []byte) uint64 {
	if p.cfg != nil && p.cfg.GasOverrides != nil {
		// Use config-based gas calculation.
		if len(input) >= 128 {
			alg := input[127]
			return p.cfg.EffectiveGas(alg)
		}
		base := pqverify.GasBaseOverhead
		if p.cfg.GasOverrides.BaseGas != nil {
			base = *p.cfg.GasOverrides.BaseGas
		}
		return base
	}
	return p.inner.RequiredGas(input)
}
