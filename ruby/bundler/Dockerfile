FROM ruby:2.7-slim-buster

RUN mkdir /app
WORKDIR /app

COPY Gemfile Gemfile.lock dependabot-ruby-bundler.gemspec /app/
#RUN bundle install

COPY . /app