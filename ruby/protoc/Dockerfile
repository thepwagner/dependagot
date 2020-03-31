FROM dependagot-protoc

RUN apt-get update \
  && apt-get -y install ruby-bundler \
  && rm -rf /var/lib/apt/lists/*

COPY Gemfile Gemfile.lock /app/
RUN bundle install

COPY protoc.sh /protoc.sh
ENTRYPOINT ["/protoc.sh"]
