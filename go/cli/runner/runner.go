package runner

import (
	"context"
	"fmt"

	"github.com/github/dependabot/go/common/dependabot/v1"
	"github.com/sirupsen/logrus"
)

type LoadingUpdater struct {
	updater dependabot_v1.UpdateService
	loader  Loader
}

type Loader interface {
	Load(ctx context.Context, path string) ([]byte, bool, error)
}

func NewLoadingUpdater(updater dependabot_v1.UpdateService, loader Loader) *LoadingUpdater {
	return &LoadingUpdater{
		updater: updater,
		loader:  loader,
	}
}

func (r *LoadingUpdater) ListDependencies(ctx context.Context) ([]*dependabot_v1.Dependency, error) {
	if err := r.loadFiles(ctx); err != nil {
		return nil, err
	}
	deps, err := r.updater.ListDependencies(ctx, &dependabot_v1.ListDependenciesRequest{})
	if err != nil {
		return nil, err
	}
	return deps.Dependencies, nil
}

func (r *LoadingUpdater) loadFiles(ctx context.Context) error {
	var fileCount, requestCount int
	var byteCount int64
	logrus.Debug("Starting file loading...")
	defer func() {
		logrus.WithFields(logrus.Fields{
			"files":    fileCount,
			"requests": requestCount,
			"bytes":    byteCount,
		}).Info("Finished file loading")
	}()

	var req dependabot_v1.FilesRequest
	for {
		logrus.WithFields(logrus.Fields{
			"included_files": len(req.Files),
		}).Debug("Requesting Files()")

		requestCount++
		fileCount += len(req.Files)
		for _, f := range req.Files {
			byteCount += int64(len(f))
		}
		res, err := r.updater.Files(ctx, &req)
		if err != nil {
			return err
		}

		required := res.GetRequiredPaths()
		optional := res.GetOptionalPaths()
		if len(required) == 0 && len(optional) == 0 {
			logrus.Debug("API requested no files, loading finished.")
			return nil
		} else {
			logrus.WithFields(logrus.Fields{
				"required": required,
				"optional": optional,
			}).Debug("API requested files")
		}

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

func (r *LoadingUpdater) loadPaths(ctx context.Context, paths []string, required bool, pathData map[string][]byte) error {
	for _, path := range paths {
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
