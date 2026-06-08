import { gzipSync } from "node:zlib";
import { minify } from "html-minifier-terser";
import { DEFAULT_KEY, sealWithChecksum, xorBytes } from "./cipher.js";

export const MINIFY_OPTIONS = {
  collapseWhitespace: true,
  removeComments: true,
  removeRedundantAttributes: true,
  removeScriptTypeAttributes: true,
  removeStyleLinkTypeAttributes: true,
  useShortDoctype: true,
  minifyCSS: true,
  minifyJS: true,
};

export const MAX_SHARE_URL_LENGTH = 8192;

export function urlLengthWarning(url) {
  if (url.length <= MAX_SHARE_URL_LENGTH) {
    return "";
  }

  return `\n*** WARNING: URL length is ${url.length} chars and exceeds ${MAX_SHARE_URL_LENGTH} characters. ***\n*** Some apps (including Slack) may truncate or refuse to send this link. ***\n\n`;
}

export async function encodeHtml(html, key = DEFAULT_KEY) {
  const minified = await minify(html, MINIFY_OPTIONS);
  const compressed = gzipSync(Buffer.from(minified, "utf8"));
  const sealed = sealWithChecksum(compressed);
  const obfuscated = xorBytes(sealed, key);
  const encoded = obfuscated.toString("base64url");

  return { minified, compressed, encoded, key };
}

export function buildShareUrl(baseUrl, encoded, key = DEFAULT_KEY) {
  let url = `${baseUrl}#d=${encoded}`;

  if (key !== DEFAULT_KEY) {
    url += `&k=${Buffer.from(key, "utf8").toString("base64url")}`;
  }

  return url;
}
