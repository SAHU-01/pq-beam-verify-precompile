// Command keygen generates post-quantum keypairs using ML-DSA-65 or SLH-DSA-128s,
// writes the keys to files, and prints the derived PQ address to stdout.
package main

import (
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pq-beam/verify-precompile/pkg/pqcrypto"
	"golang.org/x/crypto/sha3"
)

func main() {
	algorithm := flag.String("algorithm", "ml-dsa-65", `PQ algorithm: "ml-dsa-65" or "slh-dsa-128s"`)
	flag.StringVar(algorithm, "a", "ml-dsa-65", `PQ algorithm (shorthand)`)

	output := flag.String("output", ".", "output directory for key files")
	flag.StringVar(output, "o", ".", "output directory (shorthand)")

	format := flag.String("format", "hex", `encoding format for key files: "hex" or "base64"`)

	verify := flag.String("verify", "", "path to a pubkey file to run a sign/verify round-trip test")

	flag.Parse()

	if *verify != "" {
		if err := runVerify(*verify, *algorithm, *format); err != nil {
			fmt.Fprintf(os.Stderr, "verify failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := runKeygen(*algorithm, *output, *format); err != nil {
		fmt.Fprintf(os.Stderr, "keygen failed: %v\n", err)
		os.Exit(1)
	}
}

// algID maps a CLI algorithm name to the pqcrypto constant.
func algID(name string) (uint8, error) {
	switch name {
	case "ml-dsa-65":
		return pqcrypto.AlgMLDSA65, nil
	case "slh-dsa-128s":
		return pqcrypto.AlgSLHDSA128s, nil
	default:
		return 0, fmt.Errorf("unknown algorithm %q (expected \"ml-dsa-65\" or \"slh-dsa-128s\")", name)
	}
}

// encode encodes raw bytes into the requested text format.
func encode(data []byte, format string) ([]byte, error) {
	switch format {
	case "hex":
		dst := make([]byte, hex.EncodedLen(len(data)))
		hex.Encode(dst, data)
		return dst, nil
	case "base64":
		dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
		base64.StdEncoding.Encode(dst, data)
		return dst, nil
	default:
		return nil, fmt.Errorf("unknown format %q (expected \"hex\" or \"base64\")", format)
	}
}

// decode decodes text-encoded bytes back to raw.
func decode(data []byte, format string) ([]byte, error) {
	switch format {
	case "hex":
		dst := make([]byte, hex.DecodedLen(len(data)))
		n, err := hex.Decode(dst, data)
		if err != nil {
			return nil, fmt.Errorf("hex decode: %w", err)
		}
		return dst[:n], nil
	case "base64":
		dst := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
		n, err := base64.StdEncoding.Decode(dst, data)
		if err != nil {
			return nil, fmt.Errorf("base64 decode: %w", err)
		}
		return dst[:n], nil
	default:
		return nil, fmt.Errorf("unknown format %q", format)
	}
}

// pqAddress computes the PQ address: last 20 bytes of keccak256(pubkey).
func pqAddress(pubkey []byte) []byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(pubkey)
	hash := h.Sum(nil)
	return hash[len(hash)-20:]
}

func runKeygen(algorithm, output, format string) error {
	alg, err := algID(algorithm)
	if err != nil {
		return err
	}

	pubkey, seckey, err := pqcrypto.GenerateKeypair(alg)
	if err != nil {
		return fmt.Errorf("GenerateKeypair: %w", err)
	}

	// Encode keys.
	pubEncoded, err := encode(pubkey, format)
	if err != nil {
		return err
	}
	secEncoded, err := encode(seckey, format)
	if err != nil {
		return err
	}

	// Ensure output directory exists.
	if err := os.MkdirAll(output, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", output, err)
	}

	pubPath := filepath.Join(output, algorithm+".pub")
	secPath := filepath.Join(output, algorithm+".sec")

	if err := os.WriteFile(pubPath, pubEncoded, 0o644); err != nil {
		return fmt.Errorf("write pubkey: %w", err)
	}
	if err := os.WriteFile(secPath, secEncoded, 0o600); err != nil {
		return fmt.Errorf("write seckey: %w", err)
	}

	addr := pqAddress(pubkey)

	fmt.Printf("Algorithm:  %s\n", algorithm)
	fmt.Printf("Public key: %s (%d bytes)\n", pubPath, len(pubkey))
	fmt.Printf("Secret key: %s (%d bytes)\n", secPath, len(seckey))
	fmt.Printf("PQ Address: 0x%x\n", addr)

	return nil
}

func runVerify(pubkeyPath, algorithm, format string) error {
	alg, err := algID(algorithm)
	if err != nil {
		return err
	}

	// Read and decode the public key.
	pubEncoded, err := os.ReadFile(pubkeyPath)
	if err != nil {
		return fmt.Errorf("read pubkey file: %w", err)
	}
	pubkey, err := decode(pubEncoded, format)
	if err != nil {
		return fmt.Errorf("decode pubkey: %w", err)
	}

	// Derive secret key path from pubkey path: replace .pub with .sec.
	secPath := pubkeyPath[:len(pubkeyPath)-len(filepath.Ext(pubkeyPath))] + ".sec"
	secEncoded, err := os.ReadFile(secPath)
	if err != nil {
		return fmt.Errorf("read seckey file %s: %w", secPath, err)
	}
	seckey, err := decode(secEncoded, format)
	if err != nil {
		return fmt.Errorf("decode seckey: %w", err)
	}

	// Sign a test message.
	testMessage := []byte("pq-beam keygen verification test")
	signature, err := pqcrypto.Sign(alg, seckey, testMessage)
	if err != nil {
		return fmt.Errorf("sign: %w", err)
	}

	// Verify the signature.
	ok, err := pqcrypto.Verify(alg, pubkey, signature, testMessage)
	if err != nil {
		return fmt.Errorf("verify: %w", err)
	}
	if !ok {
		return fmt.Errorf("signature verification failed")
	}

	addr := pqAddress(pubkey)

	fmt.Printf("Algorithm:   %s\n", algorithm)
	fmt.Printf("Public key:  %s (%d bytes)\n", pubkeyPath, len(pubkey))
	fmt.Printf("Secret key:  %s (%d bytes)\n", secPath, len(seckey))
	fmt.Printf("Signature:   %d bytes\n", len(signature))
	fmt.Printf("PQ Address:  0x%x\n", addr)
	fmt.Println("Verification: OK")

	return nil
}
