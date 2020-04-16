FROM dependagot-protoc

RUN mkdir /protoc-ruby
WORKDIR /protoc-ruby
COPY go.mod go.sum /protoc-ruby/
RUN go mod download
RUN go install github.com/twitchtv/twirp-ruby/protoc-gen-twirp_ruby

COPY protoc.sh /protoc.sh
ENTRYPOINT ["/protoc.sh"]
