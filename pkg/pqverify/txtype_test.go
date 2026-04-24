package pqverify

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/pq-beam/verify-precompile/pkg/pqcrypto"
)

// newTestTx creates a sample PQ transaction for testing.
func newTestTx(chainID int64) *PQTransaction {
	to := [20]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
		0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14}
	return &PQTransaction{
		ChainID:  big.NewInt(chainID),
		Nonce:    42,
		GasPrice: big.NewInt(25_000_000_000), // 25 gwei
		GasLimit: 21_000,
		To:       &to,
		Value:    big.NewInt(1_000_000_000_000_000_000), // 1 ETH
		Data:     nil,
	}
}

func TestPQTransactionSignAndVerify(t *testing.T) {
	pub, sec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}

	tx := newTestTx(999)
	tx.PQPublicKey = pub

	err = SignPQTransaction(tx, sec, pqcrypto.AlgMLDSA65)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	if len(tx.PQSignature) == 0 {
		t.Fatal("signature is empty after signing")
	}

	valid, err := VerifyPQTransaction(tx)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !valid {
		t.Fatal("valid PQ transaction signature returned false")
	}

	// Verify sender address derivation
	sender := Sender(tx)
	expectedAddr := DerivePQAddress(pub)
	if sender != expectedAddr {
		t.Errorf("sender mismatch: got %x, want %x", sender, expectedAddr)
	}
}

func TestPQTransactionInvalidSignature(t *testing.T) {
	pub, sec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}

	tx := newTestTx(999)
	tx.PQPublicKey = pub

	err = SignPQTransaction(tx, sec, pqcrypto.AlgMLDSA65)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	// Tamper with the signature
	tx.PQSignature[0] ^= 0xFF

	valid, err := VerifyPQTransaction(tx)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if valid {
		t.Fatal("tampered signature should not verify")
	}

	// Also test with missing signature
	tx.PQSignature = nil
	_, err = VerifyPQTransaction(tx)
	if err != ErrMissingPQSignature {
		t.Errorf("expected ErrMissingPQSignature, got %v", err)
	}

	// Test with missing public key
	tx.PQPublicKey = nil
	_, err = VerifyPQTransaction(tx)
	if err != ErrMissingPQPublicKey {
		t.Errorf("expected ErrMissingPQPublicKey, got %v", err)
	}
}

func TestPQTransactionReplayProtection(t *testing.T) {
	pub, sec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}

	// Sign on chain 999
	tx1 := newTestTx(999)
	tx1.PQPublicKey = pub
	err = SignPQTransaction(tx1, sec, pqcrypto.AlgMLDSA65)
	if err != nil {
		t.Fatalf("sign tx1: %v", err)
	}

	// Same tx fields but on chain 1000
	tx2 := newTestTx(1000)
	tx2.PQPublicKey = pub
	tx2.PQSignature = tx1.PQSignature // reuse signature from chain 999
	tx2.PQAlgorithm = tx1.PQAlgorithm

	valid, err := VerifyPQTransaction(tx2)
	if err != nil {
		t.Fatalf("verify tx2: %v", err)
	}
	if valid {
		t.Fatal("replayed signature on different chain ID should not verify")
	}

	// Verify hashes are actually different
	hash1 := TxHash(tx1)
	hash2 := TxHash(tx2)
	if bytes.Equal(hash1, hash2) {
		t.Fatal("tx hashes should differ for different chain IDs")
	}

	// The original tx should still verify
	valid, err = VerifyPQTransaction(tx1)
	if err != nil {
		t.Fatalf("verify tx1: %v", err)
	}
	if !valid {
		t.Fatal("original tx should still verify after replay attempt")
	}
}

func TestEncodeDecode(t *testing.T) {
	pub, sec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}

	tx := newTestTx(999)
	tx.PQPublicKey = pub
	err = SignPQTransaction(tx, sec, pqcrypto.AlgMLDSA65)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	// Encode
	encoded, err := EncodePQTransaction(tx)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	// Check type prefix
	if encoded[0] != PQTxType {
		t.Errorf("type prefix: got 0x%02x, want 0x%02x", encoded[0], PQTxType)
	}

	// Decode
	decoded, err := DecodePQTransaction(encoded)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	// Verify all fields match
	if decoded.ChainID.Cmp(tx.ChainID) != 0 {
		t.Errorf("chainID: got %s, want %s", decoded.ChainID, tx.ChainID)
	}
	if decoded.Nonce != tx.Nonce {
		t.Errorf("nonce: got %d, want %d", decoded.Nonce, tx.Nonce)
	}
	if decoded.GasPrice.Cmp(tx.GasPrice) != 0 {
		t.Errorf("gasPrice: got %s, want %s", decoded.GasPrice, tx.GasPrice)
	}
	if decoded.GasLimit != tx.GasLimit {
		t.Errorf("gasLimit: got %d, want %d", decoded.GasLimit, tx.GasLimit)
	}
	if *decoded.To != *tx.To {
		t.Errorf("to: got %x, want %x", decoded.To, tx.To)
	}
	if decoded.Value.Cmp(tx.Value) != 0 {
		t.Errorf("value: got %s, want %s", decoded.Value, tx.Value)
	}
	if !bytes.Equal(decoded.Data, tx.Data) {
		t.Errorf("data: got %x, want %x", decoded.Data, tx.Data)
	}
	if decoded.PQAlgorithm != tx.PQAlgorithm {
		t.Errorf("algorithm: got %d, want %d", decoded.PQAlgorithm, tx.PQAlgorithm)
	}
	if !bytes.Equal(decoded.PQPublicKey, tx.PQPublicKey) {
		t.Error("public key mismatch after decode")
	}
	if !bytes.Equal(decoded.PQSignature, tx.PQSignature) {
		t.Error("signature mismatch after decode")
	}

	// Decoded tx should still verify
	valid, err := VerifyPQTransaction(decoded)
	if err != nil {
		t.Fatalf("verify decoded: %v", err)
	}
	if !valid {
		t.Fatal("decoded transaction should still verify")
	}
}

func TestDecodeInvalidData(t *testing.T) {
	// Too short
	_, err := DecodePQTransaction([]byte{PQTxType})
	if err != ErrTxTooShort {
		t.Errorf("expected ErrTxTooShort, got %v", err)
	}

	// Wrong type prefix
	_, err = DecodePQTransaction([]byte{0x01, 0x02, 0x03})
	if err != ErrInvalidTxType {
		t.Errorf("expected ErrInvalidTxType, got %v", err)
	}

	// Empty
	_, err = DecodePQTransaction(nil)
	if err != ErrTxTooShort {
		t.Errorf("expected ErrTxTooShort for nil, got %v", err)
	}
}

func TestTxHashDeterministic(t *testing.T) {
	tx := newTestTx(999)
	tx.PQAlgorithm = pqcrypto.AlgMLDSA65

	hash1 := TxHash(tx)
	hash2 := TxHash(tx)

	if !bytes.Equal(hash1, hash2) {
		t.Fatal("TxHash should be deterministic")
	}

	if len(hash1) != 32 {
		t.Errorf("hash length: got %d, want 32", len(hash1))
	}
}

func TestContractCreationTx(t *testing.T) {
	pub, sec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	if err != nil {
		t.Fatalf("keygen: %v", err)
	}

	// Contract creation: To is nil
	tx := &PQTransaction{
		ChainID:     big.NewInt(999),
		Nonce:       0,
		GasPrice:    big.NewInt(25_000_000_000),
		GasLimit:    100_000,
		To:          nil,
		Value:       big.NewInt(0),
		Data:        []byte{0x60, 0x80, 0x60, 0x40, 0x52}, // sample bytecode
		PQPublicKey: pub,
	}

	err = SignPQTransaction(tx, sec, pqcrypto.AlgMLDSA65)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	valid, err := VerifyPQTransaction(tx)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if !valid {
		t.Fatal("contract creation tx should verify")
	}

	// Encode/decode roundtrip for nil To
	encoded, err := EncodePQTransaction(tx)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	decoded, err := DecodePQTransaction(encoded)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	if decoded.To != nil {
		t.Error("decoded To should be nil for contract creation")
	}
}
