// Package test contains end-to-end integration tests for the PQ_VERIFY precompile.
//
// These tests exercise the full flow:
//   1. Generate PQ keypair
//   2. Derive PQ address
//   3. Sign a message
//   4. ABI-encode and call the precompile
//   5. Verify the result
//
// These run without a live chain — they test the precompile logic directly.
// For on-chain tests, see scripts/deploy_test.sh (requires local Subnet-EVM).
package test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/pq-beam/verify-precompile/pkg/pqcrypto"
	"github.com/pq-beam/verify-precompile/pkg/pqverify"
)

// TestE2E_MLDSA65_FullFlow tests the complete flow for ML-DSA-65.
func TestE2E_MLDSA65_FullFlow(t *testing.T) {
	// Step 1: Generate keypair
	pub, sec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}
	t.Logf("Generated ML-DSA-65 keypair: pub=%d bytes, sec=%d bytes", len(pub), len(sec))

	// Step 2: Derive PQ address
	addr := pqverify.DerivePQAddress(pub)
	t.Logf("PQ address: 0x%s", hex.EncodeToString(addr[:]))

	// Step 3: Sign a transaction-like message
	message := []byte("beam-pq-tx:transfer:100BEAM:to:0xdeadbeef")
	sig, err := pqcrypto.Sign(pqcrypto.AlgMLDSA65, sec, message)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	t.Logf("Signature: %d bytes", len(sig))

	// Step 4: Verify via precompile
	precompile := &pqverify.PQVerifyPrecompile{}
	input := pqverify.EncodeInput(pub, sig, message, pqcrypto.AlgMLDSA65)
	t.Logf("Precompile input: %d bytes", len(input))

	gas := precompile.RequiredGas(input)
	t.Logf("Gas required: %d", gas)

	output, err := precompile.Run(input)
	if err != nil {
		t.Fatalf("precompile Run: %v", err)
	}

	if output[31] != 1 {
		t.Fatal("FAIL: valid ML-DSA-65 signature rejected by precompile")
	}
	t.Log("PASS: ML-DSA-65 signature verified successfully")

	// Step 5: Verify rejection of tampered signature
	tamperedSig := make([]byte, len(sig))
	copy(tamperedSig, sig)
	tamperedSig[0] ^= 0xFF

	input2 := pqverify.EncodeInput(pub, tamperedSig, message, pqcrypto.AlgMLDSA65)
	output2, err := precompile.Run(input2)
	if err != nil {
		t.Fatalf("precompile Run (tampered): %v", err)
	}
	if output2[31] != 0 {
		t.Fatal("FAIL: tampered signature accepted by precompile")
	}
	t.Log("PASS: tampered signature correctly rejected")
}

// TestE2E_SLHDSA128s_FullFlow tests the complete flow for SLH-DSA-128s.
func TestE2E_SLHDSA128s_FullFlow(t *testing.T) {
	pub, sec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgSLHDSA128s)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}
	t.Logf("Generated SLH-DSA-128s keypair: pub=%d bytes, sec=%d bytes", len(pub), len(sec))

	addr := pqverify.DerivePQAddress(pub)
	t.Logf("PQ address: 0x%s", hex.EncodeToString(addr[:]))

	message := []byte("sphincs-plus-fallback-verification")
	sig, err := pqcrypto.Sign(pqcrypto.AlgSLHDSA128s, sec, message)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	precompile := &pqverify.PQVerifyPrecompile{}
	input := pqverify.EncodeInput(pub, sig, message, pqcrypto.AlgSLHDSA128s)
	gas := precompile.RequiredGas(input)
	t.Logf("Gas required: %d", gas)

	output, err := precompile.Run(input)
	if err != nil {
		t.Fatalf("precompile Run: %v", err)
	}
	if output[31] != 1 {
		t.Fatal("FAIL: valid SLH-DSA-128s signature rejected")
	}
	t.Log("PASS: SLH-DSA-128s signature verified successfully")
}

// TestE2E_DualAlgorithm verifies both algorithms work through the same precompile.
func TestE2E_DualAlgorithm(t *testing.T) {
	algorithms := []struct {
		name string
		alg  uint8
	}{
		{"ML-DSA-65", pqcrypto.AlgMLDSA65},
		{"SLH-DSA-128s", pqcrypto.AlgSLHDSA128s},
	}

	precompile := &pqverify.PQVerifyPrecompile{}

	for _, tc := range algorithms {
		t.Run(tc.name, func(t *testing.T) {
			pub, sec, err := pqcrypto.GenerateKeypair(tc.alg)
			if err != nil {
				t.Fatalf("keygen: %v", err)
			}

			msg := []byte(fmt.Sprintf("dual-alg-test-%s", tc.name))
			sig, err := pqcrypto.Sign(tc.alg, sec, msg)
			if err != nil {
				t.Fatalf("sign: %v", err)
			}

			input := pqverify.EncodeInput(pub, sig, msg, tc.alg)
			output, err := precompile.Run(input)
			if err != nil {
				t.Fatalf("Run: %v", err)
			}
			if output[31] != 1 {
				t.Fatalf("FAIL: %s verification failed", tc.name)
			}
			t.Logf("PASS: %s verified", tc.name)
		})
	}
}

// TestE2E_CrossAlgorithmRejection ensures algorithm mismatch is rejected.
func TestE2E_CrossAlgorithmRejection(t *testing.T) {
	pub, sec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}

	msg := []byte("cross-algorithm-test")
	sig, err := pqcrypto.Sign(pqcrypto.AlgMLDSA65, sec, msg)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	// Try verifying ML-DSA sig with SLH-DSA algorithm — should fail
	precompile := &pqverify.PQVerifyPrecompile{}
	input := pqverify.EncodeInput(pub, sig, msg, pqcrypto.AlgSLHDSA128s)
	output, err := precompile.Run(input)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if output[31] != 0 {
		t.Fatal("FAIL: cross-algorithm verification should fail")
	}
	t.Log("PASS: cross-algorithm correctly rejected")
}

// TestE2E_LargeMessage tests verification with a large random message.
func TestE2E_LargeMessage(t *testing.T) {
	pub, sec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}

	// 1MB message
	message := make([]byte, 1024*1024)
	_, err = rand.Read(message)
	if err != nil {
		t.Fatalf("rand: %v", err)
	}

	sig, err := pqcrypto.Sign(pqcrypto.AlgMLDSA65, sec, message)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	precompile := &pqverify.PQVerifyPrecompile{}
	input := pqverify.EncodeInput(pub, sig, message, pqcrypto.AlgMLDSA65)

	output, err := precompile.Run(input)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if output[31] != 1 {
		t.Fatal("FAIL: large message verification failed")
	}
	t.Log("PASS: 1MB message verified")
}

// TestE2E_MultipleSignatures tests rapid sequential verifications.
func TestE2E_MultipleSignatures(t *testing.T) {
	pub, sec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}

	precompile := &pqverify.PQVerifyPrecompile{}

	for i := 0; i < 50; i++ {
		msg := []byte(fmt.Sprintf("transaction-%d", i))
		sig, err := pqcrypto.Sign(pqcrypto.AlgMLDSA65, sec, msg)
		if err != nil {
			t.Fatalf("sign %d: %v", i, err)
		}

		input := pqverify.EncodeInput(pub, sig, msg, pqcrypto.AlgMLDSA65)
		output, err := precompile.Run(input)
		if err != nil {
			t.Fatalf("Run %d: %v", i, err)
		}
		if output[31] != 1 {
			t.Fatalf("FAIL: signature %d rejected", i)
		}
	}
	t.Log("PASS: 50 sequential verifications succeeded")
}

// TestE2E_AddressDerivationDeterministic verifies same pubkey always gives same address.
func TestE2E_AddressDerivationDeterministic(t *testing.T) {
	pub, _, err := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}

	addr1 := pqverify.DerivePQAddress(pub)
	addr2 := pqverify.DerivePQAddress(pub)

	if addr1 != addr2 {
		t.Fatal("FAIL: same pubkey produced different addresses")
	}
	t.Logf("PASS: deterministic address: 0x%s", hex.EncodeToString(addr1[:]))

	// Different keypair should produce different address
	pub2, _, _ := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	addr3 := pqverify.DerivePQAddress(pub2)
	if addr1 == addr3 {
		t.Fatal("FAIL: different pubkeys produced same address (collision)")
	}
	t.Log("PASS: different pubkeys produce different addresses")
}
