import { describe, it } from "node:test";
import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { dirname, join } from "node:path";
import { gzipSync } from "node:zlib";
import { encodeHtml, buildShareUrl } from "../lib/encode.js";
import { decodePayload, extractPayloadFromUrl } from "../lib/decode.js";
import { CHECKSUM_BYTES, DEFAULT_KEY } from "../lib/cipher.js";

const __dirname = dirname(fileURLToPath(import.meta.url));
const exampleHtml = readFileSync(join(__dirname, "../example.html"), "utf8");

describe("encodeHtml", () => {
  it("minifies HTML and reduces size", async () => {
    const html = "<!DOCTYPE html><html><body>  <!-- comment -->  <p>Hi</p>  </body></html>";
    const { minified, compressed, encoded } = await encodeHtml(html);

    assert.ok(minified.length < html.length);
    assert.ok(!minified.includes("comment"));
    assert.ok(compressed.length > 0);
    assert.ok(encoded.length > 0);
  });

  it("produces URL-safe base64 without padding", async () => {
    const { encoded } = await encodeHtml(exampleHtml);

    assert.match(encoded, /^[A-Za-z0-9_-]+$/);
    assert.ok(!encoded.includes("+"));
    assert.ok(!encoded.includes("/"));
    assert.ok(!encoded.endsWith("="));
  });

  it("obfuscates sealed payload without extra URL overhead beyond checksum", async () => {
    const { compressed, encoded } = await encodeHtml(exampleHtml);
    const plainEncoded = compressed.toString("base64url");

    assert.notEqual(encoded, plainEncoded);
    assert.equal(Buffer.from(encoded, "base64url").length, compressed.length + CHECKSUM_BYTES);
  });

  it("round-trips through xor and gzip decode", async () => {
    const { minified, encoded } = await encodeHtml(exampleHtml);
    const decoded = decodePayload(encoded);

    assert.equal(decoded, minified);
    assert.ok(decoded.includes("Hello from InstaHost"));
  });

  it("round-trips with a custom key", async () => {
    const { minified, encoded } = await encodeHtml(exampleHtml, "custom-key");
    const decoded = decodePayload(encoded, "custom-key");

    assert.equal(decoded, minified);
  });

  it("decodes legacy unobfuscated payloads", async () => {
    const { minified } = await encodeHtml(exampleHtml);
    const legacy = gzipSync(Buffer.from(minified, "utf8")).toString("base64url");
    const decoded = decodePayload(legacy);

    assert.equal(decoded, minified);
  });

  it("rejects tampered payloads", async () => {
    const { encoded } = await encodeHtml(exampleHtml);
    const bytes = Buffer.from(encoded, "base64url");
    bytes[0] ^= 0xff;

    assert.throws(() => decodePayload(bytes.toString("base64url")));
  });

  it("produces stable output for the same input", async () => {
    const first = await encodeHtml(exampleHtml);
    const second = await encodeHtml(exampleHtml);

    assert.equal(first.encoded, second.encoded);
    assert.equal(first.minified, second.minified);
  });
});

describe("buildShareUrl", () => {
  it("appends d query param to a plain base URL", () => {
    const url = buildShareUrl("index.html", "abc123");

    assert.equal(url, "index.html?d=abc123");
  });

  it("uses & when base URL already has query params", () => {
    const url = buildShareUrl("https://example.com/page?foo=bar", "abc123");

    assert.equal(url, "https://example.com/page?foo=bar&d=abc123");
  });

  it("omits k param for the default key", () => {
    const url = buildShareUrl("index.html", "abc123", DEFAULT_KEY);

    assert.equal(url, "index.html?d=abc123");
    assert.ok(!url.includes("k="));
  });

  it("includes k param for a custom key", () => {
    const url = buildShareUrl("index.html", "abc123", "secret");

    assert.match(url, /^index\.html\?d=abc123&k=c2VjcmV0$/);
  });
});

describe("extractPayloadFromUrl", () => {
  it("extracts encoded payload from a share URL", async () => {
    const { encoded } = await encodeHtml(exampleHtml);
    const url = buildShareUrl("https://example.com/", encoded);

    const { payload, key } = extractPayloadFromUrl(url);

    assert.equal(payload, encoded);
    assert.equal(key, DEFAULT_KEY);
  });
});
