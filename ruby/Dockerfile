FROM ruby:2.6-slim-buster

RUN mkdir /app
COPY common app/common
COPY bundler/Gemfile bundler/Gemfile.lock bundler/dependabot-ruby-bundler.gemspec /app/bundler/

WORKDIR /app/bundler
RUN bundle install

COPY . /app

ENTRYPOINT [ "bundle", "exec", "ruby", "run.rb" ]
