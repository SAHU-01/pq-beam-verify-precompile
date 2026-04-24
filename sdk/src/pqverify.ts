import { ethers } from "ethers";
import { PQAlgorithm, PQVerifyResult, PQ_VERIFY_ADDRESS } from "./types.js";

/**
 * Wraps calls to the PQ-verify precompile deployed on Beam.
 */
export class PQVerifier {
  private provider: ethers.Provider;

  constructor(provider: ethers.Provider) {
    this.provider = provider;
  }

  /**
   * ABI-encode a verification call to the precompile.
   *
   * Layout matches the precompile's Solidity ABI:
   *   abi.encode(bytes pubkey, bytes signature, bytes message, uint8 algorithm)
   *
   * The algorithm (uint8) is the LAST parameter, not the first.
   */
  encodeVerifyCall(
    pubkey: Uint8Array,
    signature: Uint8Array,
    message: Uint8Array,
    algorithm: PQAlgorithm,
  ): string {
    const coder = ethers.AbiCoder.defaultAbiCoder();
    return coder.encode(
      ["bytes", "bytes", "bytes", "uint8"],
      [pubkey, signature, message, algorithm],
    );
  }

  /**
   * Call the precompile to verify a post-quantum signature.
   */
  async verify(
    pubkey: Uint8Array,
    signature: Uint8Array,
    message: Uint8Array,
    algorithm: PQAlgorithm,
  ): Promise<PQVerifyResult> {
    const calldata = this.encodeVerifyCall(
      pubkey,
      signature,
      message,
      algorithm,
    );

    const result = await this.provider.call({
      to: PQ_VERIFY_ADDRESS,
      data: calldata,
    });

    // The precompile returns a single uint8: 1 = valid, 0 = invalid.
    const valid = result !== "0x" && BigInt(result) === 1n;

    // Estimate gas for informational purposes (optional).
    let gasUsed: bigint | undefined;
    try {
      gasUsed = await this.provider.estimateGas({
        to: PQ_VERIFY_ADDRESS,
        data: calldata,
      });
    } catch {
      // Gas estimation may not be available; leave undefined.
    }

    return { valid, gasUsed };
  }
}
