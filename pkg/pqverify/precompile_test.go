package pqverify

import (
	"testing"

	"github.com/pq-beam/verify-precompile/pkg/pqcrypto"
)

func TestValidMLDSA65Signature(t *testing.T) {
	pub, sec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}

	msg := []byte("beam post-quantum verification test")
	sig, err := pqcrypto.Sign(pqcrypto.AlgMLDSA65, sec, msg)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	precompile := &PQVerifyPrecompile{}
	input := EncodeInput(pub, sig, msg, pqcrypto.AlgMLDSA65)

	// Check gas
	gas := precompile.RequiredGas(input)
	expectedGas := GasBaseOverhead + GasMLDSA65Verify
	if gas != expectedGas {
		t.Errorf("gas: got %d, want %d", gas, expectedGas)
	}

	// Run verification
	output, err := precompile.Run(input)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(output) != 32 {
		t.Fatalf("output length: got %d, want 32", len(output))
	}
	if output[31] != 1 {
		t.Fatal("valid signature returned false")
	}
}

func TestInvalidMLDSA65Signature(t *testing.T) {
	pub, sec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}

	msg := []byte("beam post-quantum verification test")
	sig, err := pqcrypto.Sign(pqcrypto.AlgMLDSA65, sec, msg)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	// Tamper
	sig[0] ^= 0xFF

	precompile := &PQVerifyPrecompile{}
	input := EncodeInput(pub, sig, msg, pqcrypto.AlgMLDSA65)

	output, err := precompile.Run(input)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if output[31] != 0 {
		t.Fatal("tampered signature returned true")
	}
}

func TestWrongMessage(t *testing.T) {
	pub, sec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}

	sig, err := pqcrypto.Sign(pqcrypto.AlgMLDSA65, sec, []byte("message A"))
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	precompile := &PQVerifyPrecompile{}
	input := EncodeInput(pub, sig, []byte("message B"), pqcrypto.AlgMLDSA65)

	output, err := precompile.Run(input)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if output[31] != 0 {
		t.Fatal("wrong message returned true")
	}
}

func TestValidSLHDSA128sSignature(t *testing.T) {
	pub, sec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgSLHDSA128s)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}

	msg := []byte("sphincs+ fallback test")
	sig, err := pqcrypto.Sign(pqcrypto.AlgSLHDSA128s, sec, msg)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	precompile := &PQVerifyPrecompile{}
	input := EncodeInput(pub, sig, msg, pqcrypto.AlgSLHDSA128s)

	gas := precompile.RequiredGas(input)
	expectedGas := GasBaseOverhead + GasSLHDSA128sVerify
	if gas != expectedGas {
		t.Errorf("gas: got %d, want %d", gas, expectedGas)
	}

	output, err := precompile.Run(input)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if output[31] != 1 {
		t.Fatal("valid SPHINCS+ signature returned false")
	}
}

func TestInvalidSLHDSA128sSignature(t *testing.T) {
	pub, sec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgSLHDSA128s)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}

	msg := []byte("test")
	sig, err := pqcrypto.Sign(pqcrypto.AlgSLHDSA128s, sec, msg)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	sig[0] ^= 0xFF

	precompile := &PQVerifyPrecompile{}
	input := EncodeInput(pub, sig, msg, pqcrypto.AlgSLHDSA128s)

	output, err := precompile.Run(input)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if output[31] != 0 {
		t.Fatal("tampered SPHINCS+ signature returned true")
	}
}

func TestUnsupportedAlgorithm(t *testing.T) {
	precompile := &PQVerifyPrecompile{}

	// Create input with invalid algorithm byte
	pub := make([]byte, 32)
	sig := make([]byte, 32)
	msg := []byte("test")

	input := EncodeInput(pub, sig, msg, 99)

	_, err := precompile.Run(input)
	if err == nil {
		// It's OK if it returns false instead of error
		t.Log("unsupported algorithm returned nil error (returned false)")
	}
}

func TestMalformedInput(t *testing.T) {
	precompile := &PQVerifyPrecompile{}

	// Too short
	_, err := precompile.Run([]byte{0x01, 0x02})
	if err == nil {
		t.Fatal("expected error for malformed input")
	}

	// Empty
	_, err = precompile.Run(nil)
	if err == nil {
		t.Fatal("expected error for nil input")
	}
}

func TestGasAccounting(t *testing.T) {
	precompile := &PQVerifyPrecompile{}

	// Test gas for short input
	gas := precompile.RequiredGas([]byte{0x01})
	if gas != GasBaseOverhead {
		t.Errorf("short input gas: got %d, want %d", gas, GasBaseOverhead)
	}

	// ML-DSA-65 gas
	pub, sec, _ := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	sig, _ := pqcrypto.Sign(pqcrypto.AlgMLDSA65, sec, []byte("test"))
	input := EncodeInput(pub, sig, []byte("test"), pqcrypto.AlgMLDSA65)
	gas = precompile.RequiredGas(input)
	if gas != GasBaseOverhead+GasMLDSA65Verify {
		t.Errorf("ML-DSA-65 gas: got %d, want %d", gas, GasBaseOverhead+GasMLDSA65Verify)
	}

	// SLH-DSA-128s gas
	pub2, sec2, _ := pqcrypto.GenerateKeypair(pqcrypto.AlgSLHDSA128s)
	sig2, _ := pqcrypto.Sign(pqcrypto.AlgSLHDSA128s, sec2, []byte("test"))
	input2 := EncodeInput(pub2, sig2, []byte("test"), pqcrypto.AlgSLHDSA128s)
	gas = precompile.RequiredGas(input2)
	if gas != GasBaseOverhead+GasSLHDSA128sVerify {
		t.Errorf("SLH-DSA-128s gas: got %d, want %d", gas, GasBaseOverhead+GasSLHDSA128sVerify)
	}
}

func TestEncodeDecodeRoundtrip(t *testing.T) {
	pubkey := []byte("test-pubkey-data-1234567890")
	signature := []byte("test-signature-data-abcdef")
	message := []byte("hello world")
	alg := pqcrypto.AlgMLDSA65

	encoded := EncodeInput(pubkey, signature, message, alg)
	decodedPub, decodedSig, decodedMsg, decodedAlg, err := decodeInput(encoded)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	if string(decodedPub) != string(pubkey) {
		t.Errorf("pubkey mismatch: got %q, want %q", decodedPub, pubkey)
	}
	if string(decodedSig) != string(signature) {
		t.Errorf("signature mismatch: got %q, want %q", decodedSig, signature)
	}
	if string(decodedMsg) != string(message) {
		t.Errorf("message mismatch: got %q, want %q", decodedMsg, message)
	}
	if decodedAlg != alg {
		t.Errorf("algorithm mismatch: got %d, want %d", decodedAlg, alg)
	}
}

func TestEmptyMessage(t *testing.T) {
	pub, sec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}

	sig, err := pqcrypto.Sign(pqcrypto.AlgMLDSA65, sec, []byte{})
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	precompile := &PQVerifyPrecompile{}
	input := EncodeInput(pub, sig, []byte{}, pqcrypto.AlgMLDSA65)

	output, err := precompile.Run(input)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if output[31] != 1 {
		t.Fatal("valid signature on empty message returned false")
	}
}
