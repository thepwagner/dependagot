FROM golang:1.14.1

RUN apt-get update \
  && apt-get -fy install unzip ruby-bundler \
  && rm -rf /var/lib/apt/lists/*

# Protoc binary:
ARG PROTOC_VERSION=3.9.1
ARG PROTOC_SHASUM=77410d08e9a3c1ebb68afc13ee0c0fb4272c01c20bfd289adfb51b1c622bab07
RUN curl -o /tmp/protoc.zip -L https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip \
  && (echo "${PROTOC_SHASUM} /tmp/protoc.zip" | sha256sum -c) \
  && unzip -d /usr/local /tmp/protoc.zip \
  && rm /tmp/protoc.zip

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download
RUN go install github.com/golang/protobuf/protoc-gen-go

COPY protoc.sh /protoc.sh
ENTRYPOINT ["/protoc.sh"]
