import { gunzipSync } from "node:zlib";
import { DEFAULT_KEY, openWithChecksum, xorBytes } from "./cipher.js";

export function decodePayload(encoded, key = DEFAULT_KEY) {
  const bytes = Buffer.from(encoded, "base64url");

  try {
    const deobfuscated = xorBytes(bytes, key);
    const payload = openWithChecksum(deobfuscated);
    if (payload) {
      return gunzipSync(payload).toString("utf8");
    }
    return gunzipSync(deobfuscated).toString("utf8");
  } catch {
    return gunzipSync(bytes).toString("utf8");
  }
}

export function extractPayloadFromUrl(url) {
  const hashIndex = url.indexOf("#");
  const fragment = hashIndex >= 0 ? url.slice(hashIndex + 1) : "";
  const params = new URLSearchParams(fragment);
  const payload = params.get("d");
  if (!payload) {
    throw new Error("missing d in URL fragment");
  }

  const keyParam = params.get("k");
  const key = keyParam
    ? Buffer.from(keyParam, "base64url").toString("utf8")
    : DEFAULT_KEY;

  return { payload, key };
}
