package main

import (
	"context"
	"fmt"
	"time"

	"github.com/github/dependabot/go/common/dependabot/v1"
	"github.com/github/dependabot/go/runner/runner"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := runner.NewClient(ctx, "localhost:9999")
	if err != nil {
		panic(err)
	}

	res, err := client.Files(ctx, &dependabot_v1.FilesRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println("test", res)
}
