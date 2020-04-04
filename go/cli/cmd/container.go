package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/github/dependabot/go/cli/modules"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func moduleContainer(ctx context.Context, cmd *cobra.Command) (*modules.Container, error) {
	mod, err := cmdModule(cmd)
	if err != nil {
		return nil, err
	}

	factory, err := modules.NewFactoryFromEnv()
	if err != nil {
		return nil, err
	}
	return factory.NewModuleContainer(ctx, mod)
}

func cmdModule(cmd *cobra.Command) (modules.DependencyModule, error) {
	s := cmd.Flag(flagModule).Value.String()
	if s == "" {
		// TODO: autodetect?
		return "", errors.New("must provide module: --module/-m")
	}
	for _, m := range modules.Modules {
		if string(m) == s {
			return m, nil
		}
	}
	return "", fmt.Errorf("unsupported module: %q", s)
}

func dumpContainerOutput(container *modules.Container) {
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		const (
			prefix = 8
			header = "Container Output"
			width  = 80
		)
		fmt.Println("")
		fmt.Printf("%s %s %s\n", strings.Repeat("-", prefix), header, strings.Repeat("-", width-len(header)-prefix-2))
		fmt.Println(container.Output())
		fmt.Println(strings.Repeat("-", width))
	}
}
