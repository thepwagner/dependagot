package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/github/dependabot/go/cli/runner"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List dependencies",
	Long:  `List dependencies`,
	RunE:  LoadingUpdaterCommand(ListCommand),
}

// ListCommand parses module dependencies, and renders results as JSON
func ListCommand(ctx context.Context, _ *cobra.Command, lu *runner.LoadingUpdater) error {
	dependencies, err := lu.ListDependencies(ctx)
	if err != nil {
		return err
	}

	formatter := json.NewEncoder(os.Stdout)
	formatter.SetIndent("", "  ")
	output := map[string]interface{}{"dependencies": dependencies}
	return formatter.Encode(output)
}

func init() {
	rootCmd.AddCommand(listCmd)
}
