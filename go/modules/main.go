package main

import (
	"net"
	"net/http"

	"github.com/github/dependabot/go/common/dependabot/v1"
	"github.com/github/dependabot/go/modules/service"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	listener, err := net.Listen("tcp", "localhost:9999")
	if err != nil {
		panic(err)
	}

	updater := service.NewUpdate()

	handler := dependabot_v1.NewUpdateServiceServer(updater, nil)

	mux := http.NewServeMux()
	mux.Handle(handler.PathPrefix(), handler)

	if err := http.Serve(listener, mux); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
