// Package pqverify implements the PQ_VERIFY precompile for Beam's Subnet-EVM.
//
// The precompile is a stateless contract deployed at address 0x0b00 that
// verifies post-quantum digital signatures. It supports ML-DSA-65 (Dilithium)
// and SLH-DSA-128s (SPHINCS+) via the liboqs library.
//
// ABI Encoding:
//
//	Input:  abi.encode(bytes pubkey, bytes signature, bytes message, uint8 algorithm)
//	Output: abi.encode(bool valid)
//
// Algorithm IDs:
//
//	0 = ML-DSA-65   (NIST FIPS 204, primary)
//	1 = SLH-DSA-128s (NIST FIPS 205, hash-based fallback)
package pqverify

import (
	"errors"
	"math/big"

	"github.com/pq-beam/verify-precompile/pkg/pqcrypto"
)

// Precompile address: 0x0000000000000000000000000000000000000b00
var PrecompileAddress = [20]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x0b, 0x00}

// Gas costs for each algorithm.
// Benchmarked on Apple M1 Pro:
//   ML-DSA-65 verify:    ~105μs (4.2x ecrecover at ~25μs)
//   SLH-DSA-128s verify: ~432μs (17.3x ecrecover)
//   ecrecover:           ~25μs = 3,000 gas
//
// Gas = (verify_time / ecrecover_time) * 3000 * safety_margin
// Safety margin = 10x (accounts for validator hardware variance, different
// liboqs builds, and ARM vs x86 performance gaps). Final values will be
// calibrated on Beam validator hardware before mainnet activation.
const (
	GasMLDSA65Verify    uint64 = 130_000 // 4.2x ecrecover * 3000 * 10x margin ≈ 126k, rounded up
	GasSLHDSA128sVerify uint64 = 520_000 // 17.3x ecrecover * 3000 * 10x margin ≈ 519k, rounded up
	GasBaseOverhead     uint64 = 3_600   // ABI decoding + dispatch
)

var (
	ErrInvalidABI        = errors.New("pqverify: invalid ABI input")
	ErrInputTooShort     = errors.New("pqverify: input too short")
	ErrInvalidAlgorithm  = errors.New("pqverify: unsupported algorithm")
	ErrInvalidDataLength = errors.New("pqverify: data length mismatch")
)

// PQVerifyPrecompile implements the stateless precompiled contract interface.
type PQVerifyPrecompile struct{}

// RequiredGas returns the gas cost for PQ signature verification.
// It reads the algorithm byte from the input to determine cost.
func (c *PQVerifyPrecompile) RequiredGas(input []byte) uint64 {
	if len(input) < 128 {
		return GasBaseOverhead // will fail in Run(), charge minimum
	}

	// Algorithm is the 4th ABI parameter (uint8), at offset 128+31 = 159
	// But we need to find it from the dynamic encoding.
	// In practice, algorithm is the last fixed-size param.
	// ABI layout: offset_pubkey(32) + offset_sig(32) + offset_msg(32) + algorithm(32)
	// = first 128 bytes are the head, algorithm is bytes [96:128]
	alg := input[127] // last byte of the 4th 32-byte word

	switch alg {
	case pqcrypto.AlgMLDSA65:
		return GasBaseOverhead + GasMLDSA65Verify
	case pqcrypto.AlgSLHDSA128s:
		return GasBaseOverhead + GasSLHDSA128sVerify
	default:
		return GasBaseOverhead
	}
}

// Run executes the PQ signature verification.
// Input is ABI-encoded: (bytes pubkey, bytes signature, bytes message, uint8 algorithm)
// Output is ABI-encoded: (bool valid)
func (c *PQVerifyPrecompile) Run(input []byte) ([]byte, error) {
	pubkey, signature, message, alg, err := decodeInput(input)
	if err != nil {
		return nil, err
	}

	valid, err := pqcrypto.Verify(alg, pubkey, signature, message)
	if err != nil {
		// Unsupported algorithm or internal error — return false, not error
		return encodeBool(false), nil
	}

	return encodeBool(valid), nil
}

// decodeInput parses ABI-encoded input:
// (bytes pubkey, bytes signature, bytes message, uint8 algorithm)
//
// ABI encoding for mixed dynamic+static types:
// Head (128 bytes):
//
//	[0:32]   offset to pubkey data
//	[32:64]  offset to signature data
//	[64:96]  offset to message data
//	[96:128] algorithm (uint8, right-padded to 32 bytes)
//
// Tail: each bytes field is length(32) + data(padded to 32)
func decodeInput(input []byte) (pubkey, signature, message []byte, alg uint8, err error) {
	if len(input) < 128 {
		return nil, nil, nil, 0, ErrInputTooShort
	}

	// Parse head
	pubkeyOffset := new(big.Int).SetBytes(input[0:32]).Uint64()
	sigOffset := new(big.Int).SetBytes(input[32:64]).Uint64()
	msgOffset := new(big.Int).SetBytes(input[64:96]).Uint64()
	alg = input[127] // uint8 is in the last byte of the 32-byte word

	// Validate algorithm
	if _, err := pqcrypto.AlgorithmName(alg); err != nil {
		return nil, nil, nil, 0, ErrInvalidAlgorithm
	}

	// Parse dynamic fields
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

// decodeBytesAt reads an ABI-encoded bytes field at the given offset.
func decodeBytesAt(data []byte, offset uint64) ([]byte, error) {
	if offset+32 > uint64(len(data)) {
		return nil, ErrInvalidDataLength
	}

	length := new(big.Int).SetBytes(data[offset : offset+32]).Uint64()
	dataStart := offset + 32

	if dataStart+length > uint64(len(data)) {
		return nil, ErrInvalidDataLength
	}

	result := make([]byte, length)
	copy(result, data[dataStart:dataStart+length])
	return result, nil
}

// encodeBool ABI-encodes a boolean as a 32-byte word.
func encodeBool(v bool) []byte {
	result := make([]byte, 32)
	if v {
		result[31] = 1
	}
	return result
}

// EncodeInput ABI-encodes the precompile input for testing and SDK use.
func EncodeInput(pubkey, signature, message []byte, alg uint8) []byte {
	// Head: 4 x 32-byte words
	// Offsets for dynamic types are relative to start of encoding
	head := make([]byte, 128)

	// Calculate offsets: head is 128 bytes, then pubkey, sig, msg data follows
	pubkeyDataOffset := uint64(128)
	pubkeyPaddedLen := 32 + padTo32(len(pubkey))
	sigDataOffset := pubkeyDataOffset + uint64(pubkeyPaddedLen)
	sigPaddedLen := 32 + padTo32(len(signature))
	msgDataOffset := sigDataOffset + uint64(sigPaddedLen)

	// Write offsets
	writeUint256(head[0:32], pubkeyDataOffset)
	writeUint256(head[32:64], sigDataOffset)
	writeUint256(head[64:96], msgDataOffset)
	head[127] = alg // uint8 in last byte

	// Encode dynamic fields
	pubkeyEncoded := encodeBytes(pubkey)
	sigEncoded := encodeBytes(signature)
	msgEncoded := encodeBytes(message)

	result := make([]byte, 0, len(head)+len(pubkeyEncoded)+len(sigEncoded)+len(msgEncoded))
	result = append(result, head...)
	result = append(result, pubkeyEncoded...)
	result = append(result, sigEncoded...)
	result = append(result, msgEncoded...)

	return result
}

func encodeBytes(data []byte) []byte {
	paddedLen := padTo32(len(data))
	result := make([]byte, 32+paddedLen)
	writeUint256(result[0:32], uint64(len(data)))
	copy(result[32:], data)
	return result
}

func padTo32(n int) int {
	if n%32 == 0 {
		if n == 0 {
			return 32
		}
		return n
	}
	return n + (32 - n%32)
}

func writeUint256(dst []byte, val uint64) {
	b := new(big.Int).SetUint64(val)
	bytes := b.Bytes()
	copy(dst[32-len(bytes):], bytes)
}
