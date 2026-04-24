// Package pqcrypto provides Go bindings to liboqs for post-quantum signature
// verification. It wraps ML-DSA-65 (Dilithium) and SLH-DSA-128s (SPHINCS+)
// via CGo calls to the liboqs C library.
package pqcrypto

/*
#cgo CFLAGS: -I/opt/homebrew/include
#cgo LDFLAGS: -L/opt/homebrew/lib -loqs -lcrypto
#include <oqs/oqs.h>
#include <stdlib.h>
#include <string.h>

// pq_verify wraps OQS_SIG verification. Returns 0 on success, -1 on failure.
int pq_verify(const char *alg_name,
              const uint8_t *message, size_t message_len,
              const uint8_t *signature, size_t signature_len,
              const uint8_t *public_key) {
    OQS_SIG *sig = OQS_SIG_new(alg_name);
    if (sig == NULL) {
        return -1;
    }
    OQS_STATUS rc = OQS_SIG_verify(sig, message, message_len,
                                    signature, signature_len, public_key);
    OQS_SIG_free(sig);
    return (rc == OQS_SUCCESS) ? 0 : -1;
}

// pq_keygen generates a keypair. Returns 0 on success.
int pq_keygen(const char *alg_name,
              uint8_t *public_key, uint8_t *secret_key) {
    OQS_SIG *sig = OQS_SIG_new(alg_name);
    if (sig == NULL) {
        return -1;
    }
    OQS_STATUS rc = OQS_SIG_keypair(sig, public_key, secret_key);
    OQS_SIG_free(sig);
    return (rc == OQS_SUCCESS) ? 0 : -1;
}

// pq_sign signs a message. Returns 0 on success.
int pq_sign(const char *alg_name,
            uint8_t *signature, size_t *signature_len,
            const uint8_t *message, size_t message_len,
            const uint8_t *secret_key) {
    OQS_SIG *sig = OQS_SIG_new(alg_name);
    if (sig == NULL) {
        return -1;
    }
    OQS_STATUS rc = OQS_SIG_sign(sig, signature, signature_len,
                                  message, message_len, secret_key);
    OQS_SIG_free(sig);
    return (rc == OQS_SUCCESS) ? 0 : -1;
}

// pq_sig_lengths returns the lengths of pubkey, seckey, and sig for an algorithm.
int pq_sig_lengths(const char *alg_name,
                   size_t *pubkey_len, size_t *seckey_len, size_t *sig_len) {
    OQS_SIG *sig = OQS_SIG_new(alg_name);
    if (sig == NULL) {
        return -1;
    }
    *pubkey_len = sig->length_public_key;
    *seckey_len = sig->length_secret_key;
    *sig_len = sig->length_signature;
    OQS_SIG_free(sig);
    return 0;
}
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

// Algorithm identifiers matching the precompile's uint8 algorithm field.
const (
	AlgMLDSA65    uint8 = 0
	AlgSLHDSA128s uint8 = 1
)

// OQS algorithm name strings.
const (
	OQSNameMLDSA65    = "ML-DSA-65"
	OQSNameSLHDSA128s = "SLH_DSA_PURE_SHA2_128S"
)

var (
	ErrUnsupportedAlgorithm = errors.New("unsupported PQ algorithm")
	ErrVerificationFailed   = errors.New("PQ signature verification failed")
	ErrKeygenFailed         = errors.New("PQ key generation failed")
	ErrSignFailed           = errors.New("PQ signing failed")
	ErrInvalidInput         = errors.New("invalid input parameters")
)

// AlgorithmName maps the uint8 algorithm ID to the OQS algorithm name.
func AlgorithmName(alg uint8) (string, error) {
	switch alg {
	case AlgMLDSA65:
		return OQSNameMLDSA65, nil
	case AlgSLHDSA128s:
		return OQSNameSLHDSA128s, nil
	default:
		return "", fmt.Errorf("%w: %d", ErrUnsupportedAlgorithm, alg)
	}
}

// SigLengths returns (pubkeyLen, seckeyLen, sigLen) for an algorithm.
func SigLengths(alg uint8) (pubkeyLen, seckeyLen, sigLen int, err error) {
	algName, err := AlgorithmName(alg)
	if err != nil {
		return 0, 0, 0, err
	}
	cName := C.CString(algName)
	defer C.free(unsafe.Pointer(cName))

	var cPub, cSec, cSig C.size_t
	rc := C.pq_sig_lengths(cName, &cPub, &cSec, &cSig)
	if rc != 0 {
		return 0, 0, 0, fmt.Errorf("%w: failed to get lengths for %s", ErrUnsupportedAlgorithm, algName)
	}
	return int(cPub), int(cSec), int(cSig), nil
}

// Verify verifies a post-quantum signature.
func Verify(alg uint8, pubkey, signature, message []byte) (bool, error) {
	algName, err := AlgorithmName(alg)
	if err != nil {
		return false, err
	}
	if len(pubkey) == 0 || len(signature) == 0 {
		return false, ErrInvalidInput
	}

	cName := C.CString(algName)
	defer C.free(unsafe.Pointer(cName))

	var msgPtr *C.uint8_t
	if len(message) > 0 {
		msgPtr = (*C.uint8_t)(unsafe.Pointer(&message[0]))
	}

	rc := C.pq_verify(
		cName,
		msgPtr, C.size_t(len(message)),
		(*C.uint8_t)(unsafe.Pointer(&signature[0])), C.size_t(len(signature)),
		(*C.uint8_t)(unsafe.Pointer(&pubkey[0])),
	)

	if rc != 0 {
		return false, nil // verification failed, not an error
	}
	return true, nil
}

// GenerateKeypair generates a new PQ keypair.
func GenerateKeypair(alg uint8) (pubkey, seckey []byte, err error) {
	pubLen, secLen, _, err := SigLengths(alg)
	if err != nil {
		return nil, nil, err
	}

	algName, _ := AlgorithmName(alg)
	cName := C.CString(algName)
	defer C.free(unsafe.Pointer(cName))

	pubkey = make([]byte, pubLen)
	seckey = make([]byte, secLen)

	rc := C.pq_keygen(
		cName,
		(*C.uint8_t)(unsafe.Pointer(&pubkey[0])),
		(*C.uint8_t)(unsafe.Pointer(&seckey[0])),
	)
	if rc != 0 {
		return nil, nil, ErrKeygenFailed
	}
	return pubkey, seckey, nil
}

// Sign signs a message with a PQ secret key.
func Sign(alg uint8, seckey, message []byte) ([]byte, error) {
	_, _, sigLen, err := SigLengths(alg)
	if err != nil {
		return nil, err
	}
	if len(seckey) == 0 {
		return nil, ErrInvalidInput
	}

	algName, _ := AlgorithmName(alg)
	cName := C.CString(algName)
	defer C.free(unsafe.Pointer(cName))

	signature := make([]byte, sigLen)
	var actualLen C.size_t

	var msgPtr *C.uint8_t
	if len(message) > 0 {
		msgPtr = (*C.uint8_t)(unsafe.Pointer(&message[0]))
	}

	rc := C.pq_sign(
		cName,
		(*C.uint8_t)(unsafe.Pointer(&signature[0])), &actualLen,
		msgPtr, C.size_t(len(message)),
		(*C.uint8_t)(unsafe.Pointer(&seckey[0])),
	)
	if rc != 0 {
		return nil, ErrSignFailed
	}
	return signature[:int(actualLen)], nil
}
