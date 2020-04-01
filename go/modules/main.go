package main

import (
	"net"

	"github.com/github/dependabot/go/common/dependabot/v1"
	"github.com/github/dependabot/go/modules/service"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:9999")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	updater := service.NewUpdate()
	dependabot_v1.RegisterUpdateServiceServer(server, updater)

	if err := server.Serve(listener); err != nil && err != grpc.ErrServerStopped {
		panic(err)
	}
}
