/**
 * Supported post-quantum signature algorithms.
 */
export enum PQAlgorithm {
  MLDSA65 = 0,
  SLHDSA128s = 1,
}

/**
 * A post-quantum keypair.
 */
export interface PQKeypair {
  publicKey: Uint8Array;
  secretKey: Uint8Array;
}

/**
 * A post-quantum signature with its algorithm identifier.
 */
export interface PQSignature {
  signature: Uint8Array;
  algorithm: PQAlgorithm;
}

/**
 * Result returned from the on-chain verification precompile.
 */
export interface PQVerifyResult {
  valid: boolean;
  gasUsed?: bigint;
}

/**
 * Address of the PQ-verify precompile on Beam.
 */
export const PQ_VERIFY_ADDRESS =
  "0x0000000000000000000000000000000000000b00";

// ---------- ML-DSA-65 constants ----------

/** ML-DSA-65 public key size in bytes. */
export const MLDSA65_PUBKEY_SIZE = 1952;

/** ML-DSA-65 secret key size in bytes. */
export const MLDSA65_SECKEY_SIZE = 4032;

/** ML-DSA-65 signature size in bytes. */
export const MLDSA65_SIG_SIZE = 3309;

// ---------- SLH-DSA-128s constants ----------

/** SLH-DSA-128s public key size in bytes. */
export const SLHDSA128S_PUBKEY_SIZE = 32;

/** SLH-DSA-128s secret key size in bytes. */
export const SLHDSA128S_SECKEY_SIZE = 64;

/** SLH-DSA-128s signature size in bytes. */
export const SLHDSA128S_SIG_SIZE = 7856;
