import { ethers } from "ethers";

/**
 * Derive an Ethereum-style address from a post-quantum public key.
 *
 * Computes keccak256 of the public key and takes the last 20 bytes,
 * matching the standard Ethereum address derivation scheme.
 */
export function derivePQAddress(publicKey: Uint8Array): string {
  const hash = ethers.keccak256(publicKey);
  // Last 20 bytes = last 40 hex characters
  return ethers.getAddress("0x" + hash.slice(-40));
}

/**
 * Check whether raw transaction data represents a PQ-signed transaction.
 *
 * PQ transactions use a 0x50 type prefix to distinguish them from
 * standard EIP-2718 typed transactions.
 */
export function isPQTransaction(txData: string): boolean {
  const data = txData.startsWith("0x") ? txData.slice(2) : txData;
  return data.length >= 2 && data.slice(0, 2).toLowerCase() === "50";
}
