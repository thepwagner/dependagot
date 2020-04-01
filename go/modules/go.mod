module github.com/github/dependabot/go/modules

go 1.14

replace github.com/github/dependabot/go/common => ../common

require (
	github.com/github/dependabot/go/common v1.0.0
	github.com/stretchr/testify v1.5.1
	google.golang.org/grpc v1.28.0
)
