package modules

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/client"
)

type Container struct {
	docker  *client.Client
	id      string
	logs    *bytes.Buffer
	apiAddr string
}

func (c *Container) Wait(ctx context.Context) (int64, error) {
	return c.docker.ContainerWait(ctx, c.id)
}

func (c *Container) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := c.docker.ContainerKill(ctx, c.id, "SIGKILL"); err != nil {
		// A missing container is already dead and was autoremoved, that's fine
		if strings.Contains(err.Error(), "No such container: ") {
			return nil
		}
		return fmt.Errorf("killing container: %w", err)
	}
	return nil
}

func (c *Container) Output() string {
	return c.logs.String()
}
