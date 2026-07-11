#!/bin/bash
set -e

usage="Usage: ./publish.sh <target>"
target=${1:?"Missing target. ${usage}"}

root="$(dirname ${BASH_SOURCE[0]})/../src"
ext_name="bicep-ext-http"

# prefer bicep from $PATH, fall back to ~/.azure/bin/bicep
bicep_cmd=$(command -v bicep || echo "$HOME/.azure/bin/bicep")

# build various flavors
env GOOS=darwin GOARCH=arm64 go build -C $root -o "build/http-osx-arm64"
env GOOS=linux GOARCH=amd64 go build -C $root -o "build/http-linux-x64"
env GOOS=linux GOARCH=arm64 go build -C $root -o "build/http-linux-arm64"
env GOOS=windows GOARCH=amd64 go build -C $root -o "build/http-win-x64.exe"
env GOOS=windows GOARCH=arm64 go build -C $root -o "build/http-win-arm64.exe"

# publish to the registry
"$bicep_cmd" publish-extension \
  --bin-osx-arm64 "$root/build/http-osx-arm64" \
  --bin-linux-x64 "$root/build/http-linux-x64" \
  --bin-linux-arm64 "$root/build/http-linux-arm64" \
  --bin-win-x64 "$root/build/http-win-x64.exe" \
  --bin-win-arm64 "$root/build/http-win-arm64.exe" \
  --target "$target" \
  --force