module github.com/thepwagner/dependagot/go/cli

go 1.14

replace github.com/thepwagner/dependagot/go/common => ../common

require (
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0 // indirect
	github.com/thepwagner/dependagot/go/common v1.0.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/mitchellh/go-homedir v1.1.0
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/sergi/go-diff v1.1.0
	github.com/sirupsen/logrus v1.5.0
	github.com/spf13/cobra v0.0.7
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
)
