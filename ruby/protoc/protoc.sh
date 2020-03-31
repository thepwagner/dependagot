#!/bin/sh

for PROTO in "$@"; do
  grpc_tools_ruby_protoc $PROTO \
    -I src \
    --ruby_out=out/lib \
    --grpc_out=out/lib
done
