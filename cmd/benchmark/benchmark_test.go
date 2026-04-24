package main

import (
	"testing"

	"github.com/pq-beam/verify-precompile/pkg/pqcrypto"
)

var benchMsg = []byte("benchmark message for PQ signature verification")

// ---------------------------------------------------------------------------
// ML-DSA-65
// ---------------------------------------------------------------------------

func BenchmarkMLDSA65Verify(b *testing.B) {
	pub, sec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	if err != nil {
		b.Fatal(err)
	}
	sig, err := pqcrypto.Sign(pqcrypto.AlgMLDSA65, sec, benchMsg)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ok, err := pqcrypto.Verify(pqcrypto.AlgMLDSA65, pub, sig, benchMsg)
		if err != nil || !ok {
			b.Fatal("verification failed")
		}
	}
}

func BenchmarkMLDSA65Sign(b *testing.B) {
	_, sec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := pqcrypto.Sign(pqcrypto.AlgMLDSA65, sec, benchMsg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ---------------------------------------------------------------------------
// SLH-DSA-128s
// ---------------------------------------------------------------------------

func BenchmarkSLHDSA128sVerify(b *testing.B) {
	pub, sec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgSLHDSA128s)
	if err != nil {
		b.Fatal(err)
	}
	sig, err := pqcrypto.Sign(pqcrypto.AlgSLHDSA128s, sec, benchMsg)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ok, err := pqcrypto.Verify(pqcrypto.AlgSLHDSA128s, pub, sig, benchMsg)
		if err != nil || !ok {
			b.Fatal("verification failed")
		}
	}
}

func BenchmarkSLHDSA128sSign(b *testing.B) {
	_, sec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgSLHDSA128s)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := pqcrypto.Sign(pqcrypto.AlgSLHDSA128s, sec, benchMsg)
		if err != nil {
			b.Fatal(err)
		}
	}
}
