package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/github/dependabot/go/cli/loaders"
	"github.com/github/dependabot/go/cli/modules"
	"github.com/github/dependabot/go/cli/runner"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List dependencies",
	Long:  `List dependencies`,
	RunE:  ListCommand,
}

func ListCommand(cmd *cobra.Command, _ []string) error {
	ctx := context.Background()
	container, err := moduleContainer(ctx, cmd)
	if err != nil {
		return err
	}
	defer container.Close()
	defer dumpContainerOutput(container)
	r := NewLoadingUpdater(container)

	// List dependencies:
	dependencies, err := r.ListDependencies(ctx)
	if err != nil {
		return err
	}

	// Render as JSON:
	formatter := json.NewEncoder(os.Stdout)
	formatter.SetIndent("", "  ")
	output := map[string]interface{}{"dependencies": dependencies}
	return formatter.Encode(output)
}

func NewLoadingUpdater(container *modules.Container) *runner.LoadingUpdater {
	updater := modules.NewUpdaterService(container)
	loader := loaders.NewFile(".")
	return runner.NewLoadingUpdater(updater, loader)
}

func init() {
	rootCmd.AddCommand(listCmd)
}
