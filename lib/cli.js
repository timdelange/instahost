import { DEFAULT_KEY } from "./cipher.js";

export const DEFAULT_BASE_URL = "index.html";

export function usage() {
  return `Usage: share <file> [--base-url <url>] [--key <passphrase>]

  Minify, compress, obfuscate, and encode an HTML file into a shareable URL.

Options:
  --base-url <url>   Base URL for the static page (default: index.html)
  --key <passphrase> XOR obfuscation key (default: built-in key)
  -h, --help         Show this help`;
}

export function parseArgs(argv) {
  const args = argv.slice(2);
  let baseUrl = DEFAULT_BASE_URL;
  let key = DEFAULT_KEY;
  const positional = [];

  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    if (arg === "-h" || arg === "--help") {
      return { help: true };
    }
    if (arg === "--base-url") {
      baseUrl = args[++i];
      if (!baseUrl) {
        return { error: "Error: --base-url requires a value" };
      }
      continue;
    }
    if (arg === "--key") {
      key = args[++i];
      if (!key) {
        return { error: "Error: --key requires a value" };
      }
      continue;
    }
    positional.push(arg);
  }

  if (positional.length !== 1) {
    return { error: "Error: exactly one file argument is required", showUsage: true };
  }

  return { file: positional[0], baseUrl, key };
}
