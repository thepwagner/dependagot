package main

//
//import (
//	"context"
//	"net/http"
//	"time"
//
//	"github.com/github/dependabot/go/common/dependabot/v1"
//	"github.com/github/dependabot/go/runner/loaders"
//	"github.com/github/dependabot/go/runner/runner"
//	"github.com/sirupsen/logrus"
//)
//
//func main() {
//	logrus.SetLevel(logrus.DebugLevel)
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	updater := dependabot_v1.NewUpdateServiceJSONClient("http://localhost:9999", http.DefaultClient)
//	loader := loaders.NewFile(".")
//	r := runner.NewRunner(updater, loader)
//	if err := r.Run(ctx); err != nil {
//		panic(err)
//	}
//
//}
