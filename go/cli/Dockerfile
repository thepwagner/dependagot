# Dev container holds sources and tools:
FROM golang:1.14.1 AS dev

RUN mkdir /app
COPY common /app/common/

WORKDIR /app/cli
COPY cli/go.mod cli/go.sum /app/cli/
RUN go mod download
COPY cli /app/cli/

# Builder assembles binary:
FROM dev AS builder
RUN go build -o /dependagot-cli .

# Final container is only the binary:
FROM debian:buster
COPY --from=builder /dependagot-cli /dependagot-cli
ENTRYPOINT ["/dependagot-cli"]
