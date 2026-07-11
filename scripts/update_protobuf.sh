#!/bin/bash
set -e

root="$(dirname ${BASH_SOURCE[0]})/../src"

# protoc invokes the protoc-gen-go[-grpc] plugins, which `go install` places in
# GOPATH/bin. That's not on PATH by default, so add it here.
export PATH="$PATH:$(go env GOPATH)/bin"

tmpdir=$(mktemp -d)
curl -fsSL https://raw.githubusercontent.com/Azure/bicep/refs/heads/main/src/Bicep.Local.Rpc/extension.proto -o "$tmpdir/extension.proto"

mapping="Mextension.proto=bicep.azure.com/protos/extension"

protoc -I "$tmpdir" \
  --go_out="$root" --go_opt="$mapping" \
  --go-grpc_out="$root" --go-grpc_opt="$mapping" \
  extension.proto