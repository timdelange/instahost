#!/usr/bin/env node

import { readFileSync } from "node:fs";
import { resolve, basename } from "node:path";
import { parseArgs, usage } from "../lib/cli.js";
import { encodeHtml, buildShareUrl } from "../lib/encode.js";

async function main() {
  const opts = parseArgs(process.argv);
  if (opts.help) {
    console.error(usage());
    process.exit(0);
  }
  if (opts.error) {
    console.error(opts.error);
    if (opts.showUsage) console.error(usage());
    process.exit(1);
  }

  const filePath = resolve(opts.file);
  let html;
  try {
    html = readFileSync(filePath, "utf8");
  } catch (err) {
    console.error(`Error: cannot read ${basename(filePath)}: ${err.message}`);
    process.exit(1);
  }

  const { minified, compressed, encoded } = await encodeHtml(html, opts.key);
  const url = buildShareUrl(opts.baseUrl, encoded, opts.key);

  console.log(url);
  console.error(`\nOriginal:  ${html.length} bytes`);
  console.error(`Minified:  ${minified.length} bytes`);
  console.error(`Compressed: ${compressed.length} bytes`);
  console.error(`URL length: ${url.length} chars`);
}

main();
