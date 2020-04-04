package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/github/dependabot/go/common/dependabot/v1"
	"github.com/github/dependabot/go/modules/service"
	"github.com/sirupsen/logrus"
	"github.com/twitchtv/twirp"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "15:04:05.000Z07:00",
	})
	port := os.Getenv("DEPENDAGOT_PORT")
	if port == "" {
		port = "9999"
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		panic(err)
	}

	updater := service.NewUpdate()
	handler := dependabot_v1.NewUpdateServiceServer(updater, &twirp.ServerHooks{
		RequestRouted: func(ctx context.Context) (context.Context, error) {
			twirpFields(ctx).Debug("Starting API request")
			return ctx, nil
		},
		ResponseSent: func(ctx context.Context) {
			twirpFields(ctx).Debug("Handled API request")
		},
	})

	mux := http.NewServeMux()
	mux.Handle(handler.PathPrefix(), handler)

	logrus.WithField("addr", listener.Addr()).Info("API server starting")
	if err := http.Serve(listener, mux); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

func twirpFields(ctx context.Context) logrus.FieldLogger {
	f := logrus.Fields{}
	if method, ok := twirp.MethodName(ctx); ok {
		f["method"] = method
	}
	if status, ok := twirp.StatusCode(ctx); ok {
		f["status"] = status
	}
	return logrus.WithFields(f)
}
