module github.com/github/dependabot/go/modules

go 1.14

replace github.com/github/dependabot/go/common => ../common

require (
	github.com/github/dependabot/go/common v1.0.0
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.5.1
	github.com/twitchtv/twirp v5.10.1+incompatible // indirect
)
