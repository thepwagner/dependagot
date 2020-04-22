FROM ruby:2.6-slim-buster

RUN mkdir /app
COPY common /app/common

WORKDIR /app/bundler
COPY bundler/Gemfile bundler/Gemfile.lock bundler/dependagot-ruby-bundler.gemspec /app/bundler/
RUN bundle install --without development test

COPY bundler /app/bundler/
ENTRYPOINT [ "bundle", "exec", "ruby", "run.rb" ]
