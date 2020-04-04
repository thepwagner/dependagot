FROM rust:1.42.0-slim-buster AS dev

WORKDIR /src
RUN USER=root cargo new dependagot-rust-cargo

COPY common/ /src/common/

WORKDIR /src/dependagot-rust-cargo
COPY cargo/Cargo.toml cargo/Cargo.lock ./
RUN cargo build --release

COPY cargo/src ./src
RUN cargo install
ENTRYPOINT ["/usr/local/bin/dependagot-rust-cargo"]
