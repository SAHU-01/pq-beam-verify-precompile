package pqverify

// Implements PQ address derivation. PQ addresses are derived from:
//
//	keccak256(pubkey) -> last 20 bytes
//
// To distinguish from ECDSA addresses, PQ addresses use a type prefix
// in the EIP-2718 transaction envelope.

import (
	"golang.org/x/crypto/sha3"
)

// PQAddressPrefix is the type prefix used to identify PQ transactions
// in EIP-2718 typed envelopes. 0x50 = ASCII 'P' for post-quantum.
// The canonical constant is PQTxType in txtype.go.
const PQAddressPrefix byte = PQTxType

// DerivePQAddress derives a 20-byte Ethereum-style address from a
// post-quantum public key. The address is keccak256(pubkey)[12:32],
// identical to ECDSA address derivation but applied to a PQ public key.
func DerivePQAddress(pubkey []byte) [20]byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(pubkey)
	hash := h.Sum(nil)

	var addr [20]byte
	copy(addr[:], hash[12:32])
	return addr
}

// IsPQAddress is a heuristic check for whether an address was derived
// from a PQ key. Since keccak256 output is indistinguishable from random,
// there is no reliable way to determine this from the address alone.
// The real check happens during transaction validation, where the PQ
// public key is available and can be verified against the sender address.
//
// This function always returns true as a placeholder; actual PQ address
// identification relies on the transaction type (0x50) envelope.
func IsPQAddress(addr [20]byte) bool {
	// Addresses derived from PQ keys are structurally identical to ECDSA
	// addresses. The distinction is made at the transaction level (type 0x50),
	// not at the address level. This function exists as a hook for future
	// on-chain registry or bloom-filter based identification.
	_ = addr
	return true
}
