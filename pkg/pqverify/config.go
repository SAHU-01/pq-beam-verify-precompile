package pqverify

// Config holds the configuration for the PQ_VERIFY precompile.
// In Subnet-EVM, precompiles can be enabled/disabled at specific timestamps.
type Config struct {
	// Timestamp at which the precompile becomes active (nil = genesis)
	Timestamp *uint64 `json:"timestamp,omitempty"`

	// AllowList restricts who can call the precompile (nil = anyone)
	// For PQ_VERIFY, this should be nil (public verification)
	AllowList []string `json:"allowList,omitempty"`

	// GasOverrides allows chain operators to customize gas costs
	GasOverrides *GasConfig `json:"gasOverrides,omitempty"`
}

// GasConfig allows chain operators to override the default gas costs
// for each supported PQ algorithm. Nil fields use the compiled defaults.
type GasConfig struct {
	MLDSA65Gas    *uint64 `json:"mlDsa65Gas,omitempty"`
	SLHDSA128sGas *uint64 `json:"slhDsa128sGas,omitempty"`
	BaseGas       *uint64 `json:"baseGas,omitempty"`
}

// IsEnabled returns true if the precompile should be active at the given block
// timestamp. If Timestamp is nil, the precompile is active from genesis.
func (c *Config) IsEnabled(blockTimestamp uint64) bool {
	if c == nil {
		return false
	}
	if c.Timestamp == nil {
		return true // active from genesis
	}
	return blockTimestamp >= *c.Timestamp
}

// EffectiveGas returns the gas cost for the given algorithm, applying any
// overrides from the chain configuration.
func (c *Config) EffectiveGas(alg uint8) uint64 {
	base := GasBaseOverhead
	if c != nil && c.GasOverrides != nil && c.GasOverrides.BaseGas != nil {
		base = *c.GasOverrides.BaseGas
	}

	switch alg {
	case 0: // ML-DSA-65
		algGas := GasMLDSA65Verify
		if c != nil && c.GasOverrides != nil && c.GasOverrides.MLDSA65Gas != nil {
			algGas = *c.GasOverrides.MLDSA65Gas
		}
		return base + algGas
	case 1: // SLH-DSA-128s
		algGas := GasSLHDSA128sVerify
		if c != nil && c.GasOverrides != nil && c.GasOverrides.SLHDSA128sGas != nil {
			algGas = *c.GasOverrides.SLHDSA128sGas
		}
		return base + algGas
	default:
		return base
	}
}

// Verify checks the configuration for internal consistency.
func (c *Config) Verify() error {
	// Currently no constraints to check; this is a hook for future validation.
	return nil
}
