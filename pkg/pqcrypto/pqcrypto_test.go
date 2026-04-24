package pqcrypto

import (
	"testing"
)

func TestMLDSA65KeygenSignVerify(t *testing.T) {
	pub, sec, err := GenerateKeypair(AlgMLDSA65)
	if err != nil {
		t.Fatalf("keygen failed: %v", err)
	}
	t.Logf("ML-DSA-65 pubkey: %d bytes, seckey: %d bytes", len(pub), len(sec))

	msg := []byte("hello beam post-quantum world")
	sig, err := Sign(AlgMLDSA65, sec, msg)
	if err != nil {
		t.Fatalf("sign failed: %v", err)
	}
	t.Logf("ML-DSA-65 signature: %d bytes", len(sig))

	valid, err := Verify(AlgMLDSA65, pub, sig, msg)
	if err != nil {
		t.Fatalf("verify error: %v", err)
	}
	if !valid {
		t.Fatal("valid signature rejected")
	}
}

func TestMLDSA65InvalidSignature(t *testing.T) {
	pub, sec, err := GenerateKeypair(AlgMLDSA65)
	if err != nil {
		t.Fatalf("keygen failed: %v", err)
	}

	msg := []byte("original message")
	sig, err := Sign(AlgMLDSA65, sec, msg)
	if err != nil {
		t.Fatalf("sign failed: %v", err)
	}

	// Tamper with signature
	sig[0] ^= 0xFF

	valid, err := Verify(AlgMLDSA65, pub, sig, msg)
	if err != nil {
		t.Fatalf("verify error: %v", err)
	}
	if valid {
		t.Fatal("tampered signature accepted")
	}
}

func TestMLDSA65WrongMessage(t *testing.T) {
	pub, sec, err := GenerateKeypair(AlgMLDSA65)
	if err != nil {
		t.Fatalf("keygen failed: %v", err)
	}

	sig, err := Sign(AlgMLDSA65, sec, []byte("message A"))
	if err != nil {
		t.Fatalf("sign failed: %v", err)
	}

	valid, err := Verify(AlgMLDSA65, pub, sig, []byte("message B"))
	if err != nil {
		t.Fatalf("verify error: %v", err)
	}
	if valid {
		t.Fatal("signature verified against wrong message")
	}
}

func TestSLHDSA128sKeygenSignVerify(t *testing.T) {
	pub, sec, err := GenerateKeypair(AlgSLHDSA128s)
	if err != nil {
		t.Fatalf("keygen failed: %v", err)
	}
	t.Logf("SLH-DSA-128s pubkey: %d bytes, seckey: %d bytes", len(pub), len(sec))

	msg := []byte("sphincs+ fallback test")
	sig, err := Sign(AlgSLHDSA128s, sec, msg)
	if err != nil {
		t.Fatalf("sign failed: %v", err)
	}
	t.Logf("SLH-DSA-128s signature: %d bytes", len(sig))

	valid, err := Verify(AlgSLHDSA128s, pub, sig, msg)
	if err != nil {
		t.Fatalf("verify error: %v", err)
	}
	if !valid {
		t.Fatal("valid signature rejected")
	}
}

func TestSLHDSA128sInvalidSignature(t *testing.T) {
	pub, sec, err := GenerateKeypair(AlgSLHDSA128s)
	if err != nil {
		t.Fatalf("keygen failed: %v", err)
	}

	msg := []byte("test message")
	sig, err := Sign(AlgSLHDSA128s, sec, msg)
	if err != nil {
		t.Fatalf("sign failed: %v", err)
	}

	sig[0] ^= 0xFF

	valid, err := Verify(AlgSLHDSA128s, pub, sig, msg)
	if err != nil {
		t.Fatalf("verify error: %v", err)
	}
	if valid {
		t.Fatal("tampered signature accepted")
	}
}

func TestUnsupportedAlgorithm(t *testing.T) {
	_, err := AlgorithmName(99)
	if err == nil {
		t.Fatal("expected error for unsupported algorithm")
	}
}

func TestEmptyInputs(t *testing.T) {
	valid, err := Verify(AlgMLDSA65, nil, []byte{1}, []byte{1})
	if err == nil && valid {
		t.Fatal("expected error or false for nil pubkey")
	}

	valid, err = Verify(AlgMLDSA65, []byte{1}, nil, []byte{1})
	if err == nil && valid {
		t.Fatal("expected error or false for nil signature")
	}
}

func TestEmptyMessage(t *testing.T) {
	pub, sec, err := GenerateKeypair(AlgMLDSA65)
	if err != nil {
		t.Fatalf("keygen failed: %v", err)
	}

	// Sign empty message
	sig, err := Sign(AlgMLDSA65, sec, []byte{})
	if err != nil {
		t.Fatalf("sign failed: %v", err)
	}

	valid, err := Verify(AlgMLDSA65, pub, sig, []byte{})
	if err != nil {
		t.Fatalf("verify error: %v", err)
	}
	if !valid {
		t.Fatal("valid signature on empty message rejected")
	}
}

func TestSigLengths(t *testing.T) {
	pubLen, secLen, sigLen, err := SigLengths(AlgMLDSA65)
	if err != nil {
		t.Fatalf("SigLengths failed: %v", err)
	}
	t.Logf("ML-DSA-65: pub=%d sec=%d sig=%d", pubLen, secLen, sigLen)

	if pubLen == 0 || secLen == 0 || sigLen == 0 {
		t.Fatal("expected non-zero lengths")
	}

	pubLen, secLen, sigLen, err = SigLengths(AlgSLHDSA128s)
	if err != nil {
		t.Fatalf("SigLengths failed: %v", err)
	}
	t.Logf("SLH-DSA-128s: pub=%d sec=%d sig=%d", pubLen, secLen, sigLen)
}

func TestCrossAlgorithmRejection(t *testing.T) {
	// Generate ML-DSA key and try to verify with SLH-DSA algorithm
	pub, sec, err := GenerateKeypair(AlgMLDSA65)
	if err != nil {
		t.Fatalf("keygen failed: %v", err)
	}

	msg := []byte("cross-algorithm test")
	sig, err := Sign(AlgMLDSA65, sec, msg)
	if err != nil {
		t.Fatalf("sign failed: %v", err)
	}

	// Try verifying ML-DSA sig with SLH-DSA algorithm — should fail
	valid, err := Verify(AlgSLHDSA128s, pub, sig, msg)
	if valid {
		t.Fatal("cross-algorithm verification should fail")
	}
	_ = err // error is acceptable here
}
