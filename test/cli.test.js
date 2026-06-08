import { describe, it } from "node:test";
import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import { fileURLToPath } from "node:url";
import { dirname, join } from "node:path";
import { parseArgs, DEFAULT_BASE_URL } from "../lib/cli.js";
import { decodePayload, extractPayloadFromUrl } from "../lib/decode.js";
import { DEFAULT_KEY } from "../lib/cipher.js";

const __dirname = dirname(fileURLToPath(import.meta.url));
const root = join(__dirname, "..");
const shareBin = join(root, "bin/share.js");
const exampleFile = join(root, "example.html");

function runShare(args) {
  return spawnSync(process.execPath, [shareBin, ...args], {
    cwd: root,
    encoding: "utf8",
  });
}

describe("parseArgs", () => {
  it("parses file and default base URL", () => {
    const opts = parseArgs(["node", "share", "page.html"]);

    assert.deepEqual(opts, { file: "page.html", baseUrl: DEFAULT_BASE_URL, key: DEFAULT_KEY });
  });

  it("parses --base-url", () => {
    const opts = parseArgs(["node", "share", "page.html", "--base-url", "https://x.test/"]);

    assert.deepEqual(opts, { file: "page.html", baseUrl: "https://x.test/", key: DEFAULT_KEY });
  });

  it("parses --key", () => {
    const opts = parseArgs(["node", "share", "page.html", "--key", "secret"]);

    assert.deepEqual(opts, { file: "page.html", baseUrl: DEFAULT_BASE_URL, key: "secret" });
  });

  it("returns help flag", () => {
    assert.deepEqual(parseArgs(["node", "share", "--help"]), { help: true });
  });

  it("returns error when file is missing", () => {
    const opts = parseArgs(["node", "share"]);

    assert.equal(opts.error, "Error: exactly one file argument is required");
    assert.equal(opts.showUsage, true);
  });

  it("returns error when --base-url has no value", () => {
    const opts = parseArgs(["node", "share", "page.html", "--base-url"]);

    assert.equal(opts.error, "Error: --base-url requires a value");
  });
});

describe("share CLI", () => {
  it("prints a share URL on success", () => {
    const result = runShare([exampleFile]);

    assert.equal(result.status, 0);
    assert.match(result.stdout.trim(), /^index\.html\?d=[A-Za-z0-9_-]+$/);
    assert.match(result.stderr, /Original:/);
    assert.match(result.stderr, /Minified:/);
    assert.match(result.stderr, /Compressed:/);
    assert.match(result.stderr, /URL length:/);
  });

  it("honors --base-url", () => {
    const result = runShare([exampleFile, "--base-url", "https://example.com/host/"]);

    assert.equal(result.status, 0);
    assert.match(result.stdout.trim(), /^https:\/\/example\.com\/host\/\?d=[A-Za-z0-9_-]+$/);
  });

  it("includes k param when --key is used", () => {
    const result = runShare([exampleFile, "--key", "secret"]);

    assert.equal(result.status, 0);
    assert.match(result.stdout.trim(), /^index\.html\?d=[A-Za-z0-9_-]+&k=c2VjcmV0$/);
  });

  it("CLI output round-trips to original minified HTML", () => {
    const result = runShare([exampleFile]);

    assert.equal(result.status, 0);
    const url = result.stdout.trim();
    const { payload, key } = extractPayloadFromUrl(url);
    const html = decodePayload(payload, key);

    assert.ok(html.includes("Hello from InstaHost"));
    assert.ok(html.includes("<!doctype html>") || html.includes("<!DOCTYPE html>"));
  });

  it("exits with error for missing file", () => {
    const result = runShare(["does-not-exist.html"]);

    assert.notEqual(result.status, 0);
    assert.match(result.stderr, /cannot read/);
  });

  it("exits with error when no file argument is given", () => {
    const result = runShare([]);

    assert.notEqual(result.status, 0);
    assert.match(result.stderr, /exactly one file argument is required/);
  });

  it("prints help and exits 0", () => {
    const result = runShare(["--help"]);

    assert.equal(result.status, 0);
    assert.match(result.stderr, /Usage: share/);
  });
});
