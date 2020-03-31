#!/bin/sh

for PROTO in "$@"; do
  protoc $PROTO \
    -I . \
    --go_out=go \
    --ruby_out=ruby
done
