#!/bin/sh

OUT=out
for PROTO in "$@"; do
  protoc $PROTO \
    -I src \
    --go_out=$OUT \
    --twirp_out=$OUT
done
