#!/bin/sh

chmod 755 /go/bin /go

for PROTO in "$@"; do
  protoc $PROTO \
    -I src \
    --ruby_out=out/lib \
    --twirp_ruby_out=out/lib
done
