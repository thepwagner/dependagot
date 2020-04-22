package modules_test

import (
	"context"
	"testing"
	"time"

	"github.com/docker/docker/client"
	"github.com/stretchr/testify/require"
	"github.com/thepwagner/dependagot/go/cli/modules"
)

func TestNewFactory(t *testing.T) {
	t.Skip("docker")

	docker, err := client.NewEnvClient()
	require.NoError(t, err)

	f := modules.NewFactory(docker)

	ctx := context.Background()
	c, err := f.NewModuleContainer(ctx, modules.GoModules)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)
	err = c.Close()
	require.NoError(t, err)
	t.Log(c.Output())
}
