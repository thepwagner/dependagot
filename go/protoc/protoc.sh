#!/bin/sh

#!/bin/sh

for PROTO in "$@"; do
  protoc $PROTO \
    -I src \
    --go_out=plugins=grpc:out
done
