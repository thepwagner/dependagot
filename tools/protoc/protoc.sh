#!/bin/sh

for PROTO in "$@"; do
  protoc $PROTO \
    -I . \
    --go_out=plugins=grpc:go

  grpc_tools_ruby_protoc $PROTO \
    -I . \
    --ruby_out=ruby/lib \
    --grpc_out=ruby/lib
done
