import { describe, it, expect } from "vitest";
import {
  PQAlgorithm,
  PQ_VERIFY_ADDRESS,
  MLDSA65_PUBKEY_SIZE,
  MLDSA65_SECKEY_SIZE,
  MLDSA65_SIG_SIZE,
  SLHDSA128S_PUBKEY_SIZE,
  SLHDSA128S_SECKEY_SIZE,
  SLHDSA128S_SIG_SIZE,
} from "../types.js";
import { isPQTransaction } from "../address.js";

describe("PQAlgorithm enum", () => {
  it("should assign MLDSA65 = 0", () => {
    expect(PQAlgorithm.MLDSA65).toBe(0);
  });

  it("should assign SLHDSA128s = 1", () => {
    expect(PQAlgorithm.SLHDSA128s).toBe(1);
  });
});

describe("PQ_VERIFY_ADDRESS", () => {
  it("should be a valid 20-byte hex address", () => {
    expect(PQ_VERIFY_ADDRESS).toMatch(/^0x[0-9a-fA-F]{40}$/);
  });

  it("should end with 0b00", () => {
    expect(PQ_VERIFY_ADDRESS.endsWith("0b00")).toBe(true);
  });
});

describe("ML-DSA-65 size constants", () => {
  it("should have correct public key size", () => {
    expect(MLDSA65_PUBKEY_SIZE).toBe(1952);
  });

  it("should have correct secret key size", () => {
    expect(MLDSA65_SECKEY_SIZE).toBe(4032);
  });

  it("should have correct signature size", () => {
    expect(MLDSA65_SIG_SIZE).toBe(3309);
  });
});

describe("SLH-DSA-128s size constants", () => {
  it("should have correct public key size", () => {
    expect(SLHDSA128S_PUBKEY_SIZE).toBe(32);
  });

  it("should have correct secret key size", () => {
    expect(SLHDSA128S_SECKEY_SIZE).toBe(64);
  });

  it("should have correct signature size", () => {
    expect(SLHDSA128S_SIG_SIZE).toBe(7856);
  });
});

describe("isPQTransaction", () => {
  it("should return true for 0x50-prefixed data", () => {
    expect(isPQTransaction("0x50abcdef")).toBe(true);
  });

  it("should return true without 0x prefix", () => {
    expect(isPQTransaction("50abcdef")).toBe(true);
  });

  it("should return false for non-PQ transaction data", () => {
    expect(isPQTransaction("0x02abcdef")).toBe(false);
  });

  it("should return false for empty data", () => {
    expect(isPQTransaction("")).toBe(false);
    expect(isPQTransaction("0x")).toBe(false);
  });
});
