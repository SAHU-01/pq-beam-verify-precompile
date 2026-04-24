// Package main implements a benchmark tool for post-quantum signature
// verification, comparing ML-DSA-65 and SLH-DSA-128s against a SHA-256/Keccak
// hash baseline (as an ecrecover proxy).
package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/pq-beam/verify-precompile/pkg/pqcrypto"
	"golang.org/x/crypto/sha3"
)

const iterations = 1000

// BenchmarkResult holds the timing and memory results for a single algorithm.
type BenchmarkResult struct {
	Algorithm    string  `json:"algorithm"`
	Iterations   int     `json:"iterations"`
	TotalTimeMs  float64 `json:"total_time_ms"`
	AvgTimeUs    float64 `json:"avg_time_us"`
	OpsPerSecond float64 `json:"ops_per_second"`
	AvgAllocBytes uint64 `json:"avg_alloc_bytes"`
}

// BenchmarkReport is the top-level structure written to benchmarks.json.
type BenchmarkReport struct {
	Timestamp string            `json:"timestamp"`
	GoVersion string            `json:"go_version"`
	OS        string            `json:"os"`
	Arch      string            `json:"arch"`
	Results   []BenchmarkResult `json:"results"`
}

func main() {
	msg := []byte("benchmark message for PQ signature verification")

	// --- Prepare ML-DSA-65 ---
	fmt.Println("Generating ML-DSA-65 keypair...")
	mlPub, mlSec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgMLDSA65)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ML-DSA-65 keygen failed: %v\n", err)
		os.Exit(1)
	}
	mlSig, err := pqcrypto.Sign(pqcrypto.AlgMLDSA65, mlSec, msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ML-DSA-65 sign failed: %v\n", err)
		os.Exit(1)
	}

	// --- Prepare SLH-DSA-128s ---
	fmt.Println("Generating SLH-DSA-128s keypair...")
	slhPub, slhSec, err := pqcrypto.GenerateKeypair(pqcrypto.AlgSLHDSA128s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "SLH-DSA-128s keygen failed: %v\n", err)
		os.Exit(1)
	}
	slhSig, err := pqcrypto.Sign(pqcrypto.AlgSLHDSA128s, slhSec, msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "SLH-DSA-128s sign failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nRunning benchmarks (%d iterations each)...\n\n", iterations)

	// --- Benchmark ML-DSA-65 Verify ---
	mlResult := benchmarkVerify("ML-DSA-65", pqcrypto.AlgMLDSA65, mlPub, mlSig, msg)

	// --- Benchmark SLH-DSA-128s Verify ---
	slhResult := benchmarkVerify("SLH-DSA-128s", pqcrypto.AlgSLHDSA128s, slhPub, slhSig, msg)

	// --- Benchmark SHA-256 + Keccak-256 (ecrecover baseline proxy) ---
	hashResult := benchmarkHash("SHA256+Keccak256", msg)

	results := []BenchmarkResult{mlResult, slhResult, hashResult}

	// --- Print formatted table ---
	printTable(results)

	// --- Write benchmarks.json ---
	report := BenchmarkReport{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		Results:   results,
	}

	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal JSON: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("benchmarks.json", jsonData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write benchmarks.json: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\nResults written to benchmarks.json")
}

// benchmarkVerify runs iterations of PQ signature verification and measures
// CPU time and memory allocation.
func benchmarkVerify(name string, alg uint8, pub, sig, msg []byte) BenchmarkResult {
	// Force GC before measuring.
	runtime.GC()

	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	start := time.Now()
	for i := 0; i < iterations; i++ {
		ok, err := pqcrypto.Verify(alg, pub, sig, msg)
		if err != nil || !ok {
			fmt.Fprintf(os.Stderr, "%s verification failed at iteration %d\n", name, i)
			os.Exit(1)
		}
	}
	elapsed := time.Since(start)

	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	totalAllocBytes := memAfter.TotalAlloc - memBefore.TotalAlloc
	avgAllocBytes := totalAllocBytes / uint64(iterations)

	totalMs := float64(elapsed.Microseconds()) / 1000.0
	avgUs := float64(elapsed.Microseconds()) / float64(iterations)
	opsPerSec := float64(iterations) / elapsed.Seconds()

	return BenchmarkResult{
		Algorithm:     name,
		Iterations:    iterations,
		TotalTimeMs:   totalMs,
		AvgTimeUs:     avgUs,
		OpsPerSecond:  opsPerSec,
		AvgAllocBytes: avgAllocBytes,
	}
}

// benchmarkHash runs iterations of SHA-256 followed by Keccak-256, serving as
// a rough baseline proxy for ecrecover cost.
func benchmarkHash(name string, msg []byte) BenchmarkResult {
	runtime.GC()

	var memBefore runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	start := time.Now()
	for i := 0; i < iterations; i++ {
		h := sha256.Sum256(msg)
		k := sha3.NewLegacyKeccak256()
		k.Write(h[:])
		_ = k.Sum(nil)
	}
	elapsed := time.Since(start)

	var memAfter runtime.MemStats
	runtime.ReadMemStats(&memAfter)

	totalAllocBytes := memAfter.TotalAlloc - memBefore.TotalAlloc
	avgAllocBytes := totalAllocBytes / uint64(iterations)

	totalMs := float64(elapsed.Microseconds()) / 1000.0
	avgUs := float64(elapsed.Microseconds()) / float64(iterations)
	opsPerSec := float64(iterations) / elapsed.Seconds()

	return BenchmarkResult{
		Algorithm:     name,
		Iterations:    iterations,
		TotalTimeMs:   totalMs,
		AvgTimeUs:     avgUs,
		OpsPerSecond:  opsPerSec,
		AvgAllocBytes: avgAllocBytes,
	}
}

// printTable outputs a formatted comparison table to stdout.
func printTable(results []BenchmarkResult) {
	sep := "+----------------------+------------+----------------+----------------+------------------+"
	fmt.Println(sep)
	fmt.Printf("| %-20s | %10s | %14s | %14s | %16s |\n",
		"Algorithm", "Iterations", "Avg Time (us)", "Ops/sec", "Avg Alloc (B)")
	fmt.Println(sep)
	for _, r := range results {
		fmt.Printf("| %-20s | %10d | %14.2f | %14.2f | %16d |\n",
			r.Algorithm, r.Iterations, r.AvgTimeUs, r.OpsPerSecond, r.AvgAllocBytes)
	}
	fmt.Println(sep)
}
