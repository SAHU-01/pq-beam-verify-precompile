// registry_snippet.go — This file shows exactly what to add to
// precompile/registry/registry.go in your Subnet-EVM fork.
//
// DO NOT compile this file directly. It is a reference snippet.
// Copy the relevant lines into the Subnet-EVM registry.

//go:build ignore

package registry

// Add this import to precompile/registry/registry.go:
//
//   import (
//       _ "github.com/ava-labs/subnet-evm/precompile/contracts/pqverify"
//   )
//
// The blank import triggers the init() function in the pqverify package,
// which registers the module with Subnet-EVM's precompile framework.

// In precompile/contracts/pqverify/init.go (create this file in the fork):
//
//   package pqverify
//
//   import (
//       "github.com/ava-labs/subnet-evm/precompile/modules"
//   )
//
//   func init() {
//       modules.RegisterModule(Module)
//   }
