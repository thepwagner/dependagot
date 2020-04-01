#!/bin/sh

chmod 755 /go/bin /go

for PROTO in "$@"; do
  grpc_tools_ruby_protoc $PROTO \
    -I src \
    --ruby_out=out/lib \
    --grpc_out=out/lib
done
