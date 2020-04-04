package modules_test

import (
	"context"
	"testing"
	"time"

	"github.com/docker/docker/client"
	"github.com/github/dependabot/go/cli/modules"
	"github.com/stretchr/testify/require"
)

func TestNewFactory(t *testing.T) {
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
