module github.com/github/dependabot/go/modules

go 1.14

replace github.com/github/dependabot/go/common => ../common

require (
	github.com/dependabot/gomodules-extracted v1.1.0
	github.com/github/dependabot/go/common v1.0.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/sirupsen/logrus v1.5.0
	github.com/stretchr/testify v1.5.1
	github.com/twitchtv/twirp v5.10.1+incompatible
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.4 // indirect
)
