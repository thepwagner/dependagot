# Dev container holds sources and tools:
FROM golang:1.14.1 AS dev

RUN mkdir /app
COPY common /app/common/

WORKDIR /app/modules
COPY modules/go.mod modules/go.sum /app/modules/
RUN go mod download
COPY modules /app/modules/

# Builder assembles binary:
FROM dev AS builder
RUN go build -o /go-modules .

# Final container is only the binary:
FROM golang:1.14.1
RUN apt-get update \
  && apt-get -fy install git \
  && rm -rf /var/lib/apt/lists/*
ENV PATH=/usr/bin:/usr/local/go/bin
ENV DEPENDAGOT_PORT=9999
COPY --from=builder /go-modules /go-modules
ENTRYPOINT ["/go-modules"]
