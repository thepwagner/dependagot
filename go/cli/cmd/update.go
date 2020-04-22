package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thepwagner/dependagot/go/cli/diff"
	"github.com/thepwagner/dependagot/go/cli/runner"
	"github.com/thepwagner/dependagot/go/common/dependagot/v1"
)

// upgradeCmd represents the upgrade command
var upgradeCmd = &cobra.Command{
	Use:   "update",
	Short: "Update dependencies",
	RunE:  LoadingUpdaterCommand(Upgrade),
}

const (
	flagUpgradePackage = "package"
	flagTargetVersion  = "version"
)

func Upgrade(ctx context.Context, cmd *cobra.Command, lu *runner.LoadingUpdater) error {
	p, err := cmd.Flags().GetStringArray(flagUpgradePackage)
	if err != nil {
		return err
	}
	if len(p) == 0 {
		return errors.New("missing argument: --package/-p")
	}
	v, err := cmd.Flags().GetStringArray(flagTargetVersion)
	if err != nil {
		return err
	}
	if len(v) == 0 {
		return errors.New("missing argument: --version/-v")
	}
	if len(p) != len(v) {
		return fmt.Errorf("mismatch: %d packages, %d versions", len(p), len(v))
	}

	deps := make([]*dependagot_v1.Dependency, 0, len(p))
	for i, depPackage := range p {
		deps = append(deps, &dependagot_v1.Dependency{
			Package: depPackage,
			Version: v[i],
		})
	}

	newFiles, err := lu.UpdateDependencies(ctx, deps)
	if err != nil {
		return err
	}

	for path, newContent := range newFiles {
		old, _, err := lu.Loader.Load(ctx, path)
		if err != nil {
			return err
		}
		fmt.Println(path)
		fmt.Println(diff.FormatDiff(string(old), newContent))
	}
	return nil
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	upgradeCmd.Flags().StringArrayP(flagUpgradePackage, "p", nil, "Package to upgrade")
	upgradeCmd.Flags().StringArrayP(flagTargetVersion, "v", nil, "Version to upgrade to")
}
