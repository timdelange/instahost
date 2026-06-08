import { describe, it } from "node:test";
import assert from "node:assert/strict";
import { gzipSync } from "node:zlib";
import { CHECKSUM_BYTES, crc32, DEFAULT_KEY, openWithChecksum, sealWithChecksum, xorBytes } from "../lib/cipher.js";

describe("checksum", () => {
  it("seals and opens a payload", () => {
    const input = Buffer.from("gzip-like bytes");
    const sealed = sealWithChecksum(input);
    const opened = openWithChecksum(sealed);

    assert.equal(sealed.length, input.length + CHECKSUM_BYTES);
    assert.deepEqual(opened, input);
  });

  it("returns null when checksum does not match", () => {
    const sealed = sealWithChecksum(Buffer.from("payload"));
    sealed[0] ^= 0xff;

    assert.equal(openWithChecksum(sealed), null);
  });

  it("computes a stable crc32", () => {
    assert.equal(crc32(Buffer.from("test")), 0xd87f7e0c);
  });
});

describe("xorBytes", () => {
  it("preserves byte length", () => {
    const input = gzipSync(Buffer.from("hello world"));
    const output = xorBytes(input, DEFAULT_KEY);

    assert.equal(output.length, input.length);
  });

  it("is reversible with the same key", () => {
    const input = Buffer.from("test payload");
    const encrypted = xorBytes(input, "secret");
    const decrypted = xorBytes(encrypted, "secret");

    assert.deepEqual(decrypted, input);
  });

  it("changes output without changing size", () => {
    const input = gzipSync(Buffer.from("<p>hello</p>"));
    const obfuscated = xorBytes(input, DEFAULT_KEY);

    assert.notDeepEqual(obfuscated, input);
    assert.equal(obfuscated.length, input.length);
  });
});
