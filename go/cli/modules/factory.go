package modules

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"
)

type Factory struct {
	docker *client.Client
}

func NewFactory(client *client.Client) *Factory {
	return &Factory{
		docker: client,
	}
}

func NewFactoryFromEnv() (*Factory, error) {
	docker, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return NewFactory(docker), nil
}

func (f *Factory) NewModuleContainer(ctx context.Context, module DependencyModule) (*Container, error) {
	image, ok := moduleImage(module)
	if !ok {
		return nil, fmt.Errorf("unsupported module: %q", module)
	}
	logrus.WithField("image", image).Debug("Resolved image")

	const apiPort = nat.Port("9999/tcp")
	c := container.Config{
		Image: image,
		ExposedPorts: nat.PortSet{
			apiPort: {},
		},
		Env: []string{
			"DEPENDAGOT_PORT=9999",
		},
	}
	h := container.HostConfig{
		AutoRemove: true,
		PortBindings: nat.PortMap{
			apiPort: []nat.PortBinding{{HostIP: "127.0.0.1", HostPort: "0"}},
		},
	}
	n := network.NetworkingConfig{
		//EndpointsConfig: map[string]*network.EndpointSettings{
		//	"bridge": {},
		//},
	}
	create, err := f.docker.ContainerCreate(ctx, &c, &h, &n, "")
	if err != nil {
		return nil, fmt.Errorf("creating container: %w", err)
	}
	containerID := create.ID
	logrus.WithField("container_id", containerID).Debug("Created container")

	if err := f.docker.ContainerStart(ctx, containerID, types.ContainerStartOptions{}); err != nil {
		return nil, fmt.Errorf("starting container: %w", err)
	}

	inspect, err := f.docker.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("inspecting container: %w", err)
	}
	apiPortBindings, ok := inspect.NetworkSettings.Ports[apiPort]
	if !ok || len(apiPortBindings) == 0 {
		return nil, errors.New("api port not bound")
	}
	apiPortBinding := apiPortBindings[0]
	apiAddr := fmt.Sprintf("http://%s:%s", apiPortBinding.HostIP, apiPortBinding.HostPort)
	logrus.WithFields(logrus.Fields{
		"host": apiPortBinding.HostIP,
		"port": apiPortBinding.HostPort,
	}).Debug("Bound container to port")

	var buf bytes.Buffer
	logs, err := f.docker.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		Follow:     true,
		ShowStderr: true,
		ShowStdout: true,
	})
	if err != nil {
		return nil, fmt.Errorf("tailing container logs: %w", err)
	}
	go func() {
		defer logs.Close()
		_, _ = stdcopy.StdCopy(&buf, &buf, logs)
	}()

	return &Container{
		docker:  f.docker,
		id:      containerID,
		logs:    &buf,
		apiAddr: apiAddr,
	}, nil
}

// TODO: versions
func moduleImage(module DependencyModule) (string, bool) {
	switch module {
	case GoModules:
		return "dependagot-go-modules:latest", true
	case RubyBundler:
		return "dependagot-ruby-bundler:latest", true
	default:
		return "", false
	}
}
