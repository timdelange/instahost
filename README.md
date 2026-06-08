# InstaHost

Share static HTML pages via a single URL. A CLI minifies and compresses your HTML into a URL-safe base64 payload; a static decoder page renders it in the browser.

## Prerequisites

- Node.js 18+
- A modern browser with `DecompressionStream` support (recent Chrome, Firefox, or Safari)

## Setup

```bash
npm install
```

Optionally link the CLI globally:

```bash
npm link
```

## Quick start

1. Create or use an HTML file (see `example.html`).

2. Generate a share URL:

```bash
npm run share -- example.html
# or
node bin/share.js example.html
# or, if linked globally
share example.html
```

3. Open the printed URL in your browser.

The CLI prints the URL to stdout and size stats to stderr:

```
index.html?d=H4sIAAAAAAAAE02QvU7DQBCE...

Original:  423 bytes
Minified:  335 bytes
Compressed: 265 bytes
URL length: 367 chars
```

## Hosting the decoder page

The generated URLs point at `index.html` by default. Host that file somewhere static:

**Local testing**

```bash
npx serve .
# then generate with:
node bin/share.js example.html --base-url http://localhost:3000/
```

**GitHub Pages**

Deploy the repo (or just `index.html`) to GitHub Pages, then pass your site URL:

```bash
node bin/share.js mypage.html --base-url https://youruser.github.io/instahost/
```

## CLI reference

```
share <file> [--base-url <url>]
```

| Option | Description |
|--------|-------------|
| `<file>` | Path to the HTML file to share |
| `--base-url <url>` | Base URL of the hosted decoder page (default: `index.html`) |
| `--key <passphrase>` | XOR obfuscation key (default: built-in; adds `k` param when custom) |
| `-h`, `--help` | Show help |

The `--base-url` value should be the full path to `index.html` on your host, with or without a trailing slash. The CLI appends `?d=<payload>` automatically.

## How it works

```
HTML file → minify → gzip → XOR obfuscate → base64url → URL → browser decodes → render
```

1. **Minify** — whitespace, comments, and redundant attributes are stripped.
2. **Compress** — gzip reduces payload size for shorter URLs.
3. **Obfuscate** — a lightweight XOR cipher with a repeating key (zero byte overhead).
4. **Encode** — URL-safe base64 goes in the `d` query parameter.
5. **Render** — `index.html` decodes, deobfuscates, decompresses, and writes the HTML.

The default key is built into both the CLI and decoder page, so obfuscation adds no extra URL length. Use `--key` for a custom passphrase; the key is then included as a `k` query parameter.

```bash
node bin/share.js page.html --key my-secret
# index.html?d=...&k=bXktc2VjcmV0
```

This is weak obfuscation, not strong encryption — enough to hide content from casual inspection, not to protect secrets.

## Limitations

- **URL length** — Browsers cap URLs around 2,000–8,000 characters. Keep pages small; minification and compression help but very large documents will not fit.
- **Self-contained HTML** — External assets (images, CSS, JS from other URLs) still load from their original sources. Inline or data-URI assets work best for fully portable pages.
- **No server** — Content lives entirely in the URL. There is no backend storage or editing.

## Testing

```bash
npm test
```

Tests cover HTML minification, gzip/base64url encoding, URL building, round-trip decoding, argument parsing, and CLI integration (success, errors, and `--base-url`).

## Project layout

```
index.html      Static decoder page (host this)
bin/share.js    Share CLI
lib/            Encoding, decoding, and CLI helpers
test/           Test suite
example.html    Sample page for testing
```
