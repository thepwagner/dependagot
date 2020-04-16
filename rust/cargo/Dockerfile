FROM rust:1.42.0-buster AS dev

WORKDIR /src
RUN USER=root cargo new dependagot-rust-cargo

COPY common/ /src/common/

WORKDIR /src/dependagot-rust-cargo
COPY cargo/Cargo.toml cargo/Cargo.lock ./
RUN mkdir .cargo
RUN cargo vendor > .cargo/config

COPY cargo/src ./src
RUN cargo build --release

FROM debian:buster-slim
RUN apt-get update \
  && apt-get -fy install libcurl4 openssl \
  && rm -rf /var/lib/apt/lists/*
COPY --from=dev /src/dependagot-rust-cargo/target/release/dependagot-rust-cargo /dependagot-rust-cargo
ENTRYPOINT ["/dependagot-rust-cargo"]
