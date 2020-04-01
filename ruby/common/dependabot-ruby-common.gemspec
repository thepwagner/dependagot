Gem::Specification.new do |s|
  s.name        = 'dependabot-ruby-common'
  s.version     = '0.1.0'
  s.licenses    = ['Nonstandard']
  s.summary     = ''
  s.authors     = ['thepwagner@github.com']

  s.files       = Dir["**/*.rb"]
  s.require_paths = ["lib"]

  s.add_dependency 'grpc', '~> 1.0'
end
