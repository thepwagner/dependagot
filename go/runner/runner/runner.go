package runner

import (
	"context"
	"fmt"

	"github.com/github/dependabot/go/common/dependabot/v1"
	"github.com/sirupsen/logrus"
)

type Runner struct {
	updater dependabot_v1.UpdateService
	loader  Loader
}

type Loader interface {
	Load(ctx context.Context, path string) ([]byte, bool, error)
}

func NewRunner(updater dependabot_v1.UpdateService, loader Loader) *Runner {
	return &Runner{
		updater: updater,
		loader:  loader,
	}
}

func (r *Runner) Run(ctx context.Context) error {
	if err := r.loadFiles(ctx); err != nil {
		return err
	}

	deps, err := r.updater.ListDependencies(ctx, &dependabot_v1.ListDependenciesRequest{})
	if err != nil {
		return err
	}
	logrus.WithField("deps", fmt.Sprintf("%+v", deps)).Debug("dep")
	return nil
}

func (r *Runner) loadFiles(ctx context.Context) error {
	var req dependabot_v1.FilesRequest
	for {
		res, err := r.updater.Files(ctx, &req)
		if err != nil {
			return err
		}

		required := res.GetRequiredPaths()
		optional := res.GetOptionalPaths()
		files := make(map[string][]byte, len(required)+len(optional))
		if err := r.loadPaths(ctx, required, true, files); err != nil {
			return err
		}
		if err := r.loadPaths(ctx, optional, false, files); err != nil {
			return err
		}

		// If we're adding no files this round, return
		if len(files) == 0 {
			return nil
		}
		req.Files = files
	}
}

func (r *Runner) loadPaths(ctx context.Context, paths []string, required bool, pathData map[string][]byte) error {
	for _, path := range paths {
		logrus.WithFields(logrus.Fields{
			"required": required,
			"path":     path,
		}).Debug("Loading path...")
		data, ok, err := r.loader.Load(ctx, path)
		if err != nil {
			return fmt.Errorf("loading %q: %w", path, err)
		}
		if required && !ok {
			return fmt.Errorf("required path not found: %q", path)
		}
		pathData[path] = data
	}
	return nil
}
