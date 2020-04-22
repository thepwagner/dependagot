package runner

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/thepwagner/dependagot/go/common/dependagot/v1"
)

type LoadingUpdater struct {
	Updater dependagot_v1.UpdateService
	Loader  Loader
}

type Loader interface {
	Load(ctx context.Context, path string) ([]byte, bool, error)
}

func NewLoadingUpdater(updater dependagot_v1.UpdateService, loader Loader) *LoadingUpdater {
	return &LoadingUpdater{
		Updater: updater,
		Loader:  loader,
	}
}

func (r *LoadingUpdater) ListDependencies(ctx context.Context) ([]*dependagot_v1.Dependency, error) {
	if err := r.loadFiles(ctx); err != nil {
		return nil, err
	}
	res, err := r.Updater.ListDependencies(ctx, &dependagot_v1.ListDependenciesRequest{})
	if err != nil {
		return nil, err
	}
	return res.Dependencies, nil
}

func (r *LoadingUpdater) UpdateDependencies(ctx context.Context, deps []*dependagot_v1.Dependency) (map[string]string, error) {
	if err := r.loadFiles(ctx); err != nil {
		return nil, err
	}
	res, err := r.Updater.UpdateDependencies(ctx, &dependagot_v1.UpdateDependenciesRequest{
		Dependencies: deps,
	})
	if err != nil {
		return nil, err
	}
	return res.NewFiles, nil
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

	var req dependagot_v1.FilesRequest
	for {
		reqFiles := len(req.Files)
		var reqBytes int64
		for _, f := range req.Files {
			reqBytes += int64(len(f))
		}
		requestCount++
		fileCount += reqFiles
		byteCount += reqBytes
		logrus.WithFields(logrus.Fields{
			"included_files": reqFiles,
			"included_bytes": reqBytes,
		}).Debug("Requesting Files()")

		res, err := r.Updater.Files(ctx, &req)
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
		data, ok, err := r.Loader.Load(ctx, path)
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
