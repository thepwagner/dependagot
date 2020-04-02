FROM dependagot-protoc

COPY go.mod go.sum /app/
WORKDIR /app
RUN go mod download
RUN go install github.com/twitchtv/twirp/protoc-gen-twirp

COPY protoc.sh /protoc.sh
ENTRYPOINT ["/protoc.sh"]
