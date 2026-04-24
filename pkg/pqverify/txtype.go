package pqverify

// PQ transaction type for EIP-2718 typed envelopes.
// Type 0x50 ("P" for post-quantum) carries PQ signature fields
// instead of ECDSA v, r, s.

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"

	"github.com/pq-beam/verify-precompile/pkg/pqcrypto"
)

// PQTxType is the EIP-2718 transaction type for post-quantum signed transactions.
const PQTxType byte = 0x50

var (
	ErrMissingPQSignature = errors.New("pqtx: missing PQ signature")
	ErrMissingPQPublicKey = errors.New("pqtx: missing PQ public key")
	ErrInvalidTxType      = errors.New("pqtx: invalid transaction type prefix")
	ErrTxTooShort         = errors.New("pqtx: encoded transaction too short")
)

// PQTransaction represents a post-quantum signed Ethereum transaction.
// It replaces the ECDSA v/r/s fields with PQ algorithm ID, public key,
// and signature.
type PQTransaction struct {
	ChainID     *big.Int `json:"chainId"`
	Nonce       uint64   `json:"nonce"`
	GasPrice    *big.Int `json:"gasPrice"`
	GasLimit    uint64   `json:"gasLimit"`
	To          *[20]byte `json:"to"` // nil for contract creation
	Value       *big.Int `json:"value"`
	Data        []byte   `json:"data"`
	PQAlgorithm uint8    `json:"pqAlgorithm"`
	PQPublicKey []byte   `json:"pqPublicKey"`
	PQSignature []byte   `json:"pqSignature"`
}

// rlpUnsignedTx is the RLP-serializable form of the unsigned transaction
// fields used for hash computation. The signature is computed over:
//
//	keccak256(0x50 || rlp([chainId, nonce, gasPrice, gasLimit, to, value, data, pqAlgorithm]))
type rlpUnsignedTx struct {
	ChainID     *big.Int
	Nonce       uint64
	GasPrice    *big.Int
	GasLimit    uint64
	To          []byte // empty for contract creation, 20 bytes otherwise
	Value       *big.Int
	Data        []byte
	PQAlgorithm uint8
}

// rlpSignedTx is the RLP-serializable form of the full signed transaction.
type rlpSignedTx struct {
	ChainID     *big.Int
	Nonce       uint64
	GasPrice    *big.Int
	GasLimit    uint64
	To          []byte // empty for contract creation, 20 bytes otherwise
	Value       *big.Int
	Data        []byte
	PQAlgorithm uint8
	PQPublicKey []byte
	PQSignature []byte
}

// TxHash computes the keccak256 hash of the unsigned transaction fields.
// This is the message that gets signed by the PQ signature scheme.
// The hash covers: 0x50 || rlp([chainId, nonce, gasPrice, gasLimit, to, value, data, pqAlgorithm])
func TxHash(tx *PQTransaction) []byte {
	unsigned := rlpUnsignedTx{
		ChainID:     tx.ChainID,
		Nonce:       tx.Nonce,
		GasPrice:    tx.GasPrice,
		GasLimit:    tx.GasLimit,
		To:          addrToBytes(tx.To),
		Value:       tx.Value,
		Data:        tx.Data,
		PQAlgorithm: tx.PQAlgorithm,
	}

	encoded, err := rlp.EncodeToBytes(unsigned)
	if err != nil {
		// This should never happen with valid Go types.
		panic("pqtx: failed to RLP encode unsigned tx: " + err.Error())
	}

	// Prepend the type byte, then hash.
	payload := make([]byte, 1+len(encoded))
	payload[0] = PQTxType
	copy(payload[1:], encoded)

	h := sha3.NewLegacyKeccak256()
	h.Write(payload)
	return h.Sum(nil)
}

// SignPQTransaction signs the transaction hash with the given secret key
// using the specified PQ algorithm. It populates tx.PQSignature and
// tx.PQAlgorithm.
func SignPQTransaction(tx *PQTransaction, seckey []byte, alg uint8) error {
	tx.PQAlgorithm = alg
	hash := TxHash(tx)

	sig, err := pqcrypto.Sign(alg, seckey, hash)
	if err != nil {
		return err
	}

	tx.PQSignature = sig
	return nil
}

// VerifyPQTransaction verifies the PQ signature on a transaction using the
// embedded public key and algorithm. Returns (true, nil) if the signature
// is valid.
func VerifyPQTransaction(tx *PQTransaction) (bool, error) {
	if len(tx.PQPublicKey) == 0 {
		return false, ErrMissingPQPublicKey
	}
	if len(tx.PQSignature) == 0 {
		return false, ErrMissingPQSignature
	}

	hash := TxHash(tx)
	return pqcrypto.Verify(tx.PQAlgorithm, tx.PQPublicKey, tx.PQSignature, hash)
}

// EncodePQTransaction RLP-encodes a signed PQ transaction with the 0x50
// type prefix: 0x50 || rlp([chainId, nonce, gasPrice, gasLimit, to, value, data, pqAlgorithm, pqPublicKey, pqSignature])
func EncodePQTransaction(tx *PQTransaction) ([]byte, error) {
	signed := rlpSignedTx{
		ChainID:     tx.ChainID,
		Nonce:       tx.Nonce,
		GasPrice:    tx.GasPrice,
		GasLimit:    tx.GasLimit,
		To:          addrToBytes(tx.To),
		Value:       tx.Value,
		Data:        tx.Data,
		PQAlgorithm: tx.PQAlgorithm,
		PQPublicKey: tx.PQPublicKey,
		PQSignature: tx.PQSignature,
	}

	encoded, err := rlp.EncodeToBytes(signed)
	if err != nil {
		return nil, err
	}

	// Prepend the type byte.
	result := make([]byte, 1+len(encoded))
	result[0] = PQTxType
	copy(result[1:], encoded)

	return result, nil
}

// DecodePQTransaction decodes a type-prefixed RLP-encoded PQ transaction.
// The input must begin with the 0x50 type byte.
func DecodePQTransaction(data []byte) (*PQTransaction, error) {
	if len(data) < 2 {
		return nil, ErrTxTooShort
	}
	if data[0] != PQTxType {
		return nil, ErrInvalidTxType
	}

	var signed rlpSignedTx
	if err := rlp.DecodeBytes(data[1:], &signed); err != nil {
		return nil, err
	}

	return &PQTransaction{
		ChainID:     signed.ChainID,
		Nonce:       signed.Nonce,
		GasPrice:    signed.GasPrice,
		GasLimit:    signed.GasLimit,
		To:          bytesToAddr(signed.To),
		Value:       signed.Value,
		Data:        signed.Data,
		PQAlgorithm: signed.PQAlgorithm,
		PQPublicKey: signed.PQPublicKey,
		PQSignature: signed.PQSignature,
	}, nil
}

// Sender derives the sender address from the PQ public key embedded in the
// transaction. This is equivalent to ecrecover for ECDSA transactions.
func Sender(tx *PQTransaction) [20]byte {
	return DerivePQAddress(tx.PQPublicKey)
}

// addrToBytes converts *[20]byte to []byte for RLP encoding.
// nil (contract creation) becomes empty slice.
func addrToBytes(addr *[20]byte) []byte {
	if addr == nil {
		return []byte{}
	}
	return addr[:]
}

// bytesToAddr converts []byte back to *[20]byte after RLP decoding.
// Empty slice becomes nil (contract creation).
func bytesToAddr(b []byte) *[20]byte {
	if len(b) == 0 {
		return nil
	}
	var addr [20]byte
	copy(addr[:], b)
	return &addr
}
