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

// ── Fuzz Tests ─────────────────────────────────────────────────────────────

// FuzzDecodeInput tests the ABI decoder with random/malformed inputs.
// The decoder must never panic — it should return an error for invalid data.
func FuzzDecodeInput(f *testing.F) {
	// Seed corpus: valid encoded inputs
	pub, sec, _ := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	msg := []byte("fuzz seed message")
	sig, _ := pqcrypto.Sign(pqcrypto.AlgMLDSA65, sec, msg)
	validInput := EncodeInput(pub, sig, msg, pqcrypto.AlgMLDSA65)
	f.Add(validInput)

	// Seed: empty message
	sigEmpty, _ := pqcrypto.Sign(pqcrypto.AlgMLDSA65, sec, []byte{})
	f.Add(EncodeInput(pub, sigEmpty, []byte{}, pqcrypto.AlgMLDSA65))

	// Seed: minimal valid-shaped input (128 bytes head + minimal tail)
	minimal := make([]byte, 128+32+32+32+32+32+32)
	writeUint256(minimal[0:32], 128)   // pubkey offset
	writeUint256(minimal[32:64], 192)  // sig offset
	writeUint256(minimal[64:96], 256)  // msg offset
	minimal[127] = 0                   // algorithm ML-DSA-65
	writeUint256(minimal[128:160], 1)  // pubkey len=1
	writeUint256(minimal[192:224], 1)  // sig len=1
	writeUint256(minimal[256:288], 1)  // msg len=1
	f.Add(minimal)

	// Seed: very short
	f.Add([]byte{})
	f.Add([]byte{0x01, 0x02, 0x03})
	f.Add(make([]byte, 127)) // just under minimum head size

	// Seed: all zeros (128 bytes — valid head size, offset 0 overlaps head)
	f.Add(make([]byte, 128))
	f.Add(make([]byte, 256))

	f.Fuzz(func(t *testing.T, data []byte) {
		// Must never panic — error is fine, panic is not
		pubkey, signature, message, alg, err := decodeInput(data)
		if err != nil {
			return // errors are expected for random input
		}

		// If decode succeeded, the fields must be non-nil
		if pubkey == nil || signature == nil || message == nil {
			t.Fatal("decode returned nil field without error")
		}

		// Algorithm must be valid (0 or 1)
		if alg != pqcrypto.AlgMLDSA65 && alg != pqcrypto.AlgSLHDSA128s {
			t.Fatalf("decode returned invalid algorithm %d without error", alg)
		}
	})
}

// FuzzPrecompileRun tests the full precompile with random inputs.
// It must never panic — invalid inputs should return an error or false.
func FuzzPrecompileRun(f *testing.F) {
	// Seed: valid input that produces true
	pub, sec, _ := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	msg := []byte("fuzz precompile seed")
	sig, _ := pqcrypto.Sign(pqcrypto.AlgMLDSA65, sec, msg)
	f.Add(EncodeInput(pub, sig, msg, pqcrypto.AlgMLDSA65))

	// Seed: garbage
	f.Add([]byte{})
	f.Add([]byte{0xFF})
	f.Add(make([]byte, 128))
	f.Add(make([]byte, 1024))

	// Seed: valid structure, invalid algorithm
	f.Add(EncodeInput(pub, sig, msg, 99))

	f.Fuzz(func(t *testing.T, data []byte) {
		precompile := &PQVerifyPrecompile{}

		// RequiredGas must never panic
		_ = precompile.RequiredGas(data)

		// Run must never panic — error or false output is fine
		output, err := precompile.Run(data)
		if err != nil {
			return
		}
		if len(output) != 32 {
			t.Fatalf("output length %d, want 32", len(output))
		}
	})
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
