# Dependagot

## This is not endorsed by or associated with GitHub, Dependabot, etc.

It's not even endorsed by me, strictly for science and lulz.

## What is this?

This repository is a tire fire exploring the question:
What if every module in [dependabot-core](https://github.com/dependabot/dependabot-core) was implemented in the package manager's native language? (e.g. Javascript for NPM, Ruby for bundler).

The name is a portmanteau of "Dependabot" and "polygot".

Each module builds a Docker container that exposes a Twirp endpoint when started:
* `/go/cli` is a CLI that delegates to Docker containers
* `/go/modules` is an incomplete implementation supporting `go.mod` updates
* `/rust/cargo` is an incomplete implementation supporting `Cargo.toml`
  This is the first Rust project I've attempted, double-lulz. :heart:
* `/ruby/bundler` is a stubbed implementation in Ruby

Each language tends to have a `common` library that package managers in the language might share.
This pattern provides more test scenarios: `/go/common/go.mod` will have only external dependencies, while `/go/modules/go.mod` is a sample monorepo project/path dependencies.


### Build

All tools and builds happen in Docker containers.
[Gradle](https://github.com/gradle/gradle) is used to parallelize builds and manage dependencies between ecosystems.
