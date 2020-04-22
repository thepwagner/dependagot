Gem::Specification.new do |s|
  s.name        = 'dependagot-ruby-common'
  s.version     = '0.1.0'
  s.licenses    = ['Nonstandard']
  s.summary     = ''
  s.authors     = ['thepwagner@github.com']

  s.files       = Dir["**/*.rb"]
  s.require_paths = ["lib"]

  s.add_dependency 'twirp', '~> 1.4'
end
