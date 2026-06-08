#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

mkdir -p dist

platforms=(
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
  "windows/arm64"
)

for platform in "${platforms[@]}"; do
  GOOS="${platform%%/*}"
  GOARCH="${platform##*/}"
  ext=""
  if [ "$GOOS" = "windows" ]; then
    ext=".exe"
  fi
  out="dist/share-${GOOS}-${GOARCH}${ext}"
  echo "Building $out"
  GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="-s -w" -o "$out" ./cmd/share
done

if [ "$(uname -s)" = "Darwin" ] && command -v lipo >/dev/null 2>&1; then
  echo "Building dist/share-darwin (universal amd64 + arm64)"
  lipo -create -output dist/share-darwin dist/share-darwin-amd64 dist/share-darwin-arm64
fi

echo "Built release binaries in dist/"
