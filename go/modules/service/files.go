package service

import (
	"context"

	dependabot_v1 "github.com/github/dependabot/go/common/dependabot/v1"
	"github.com/github/dependabot/go/modules/modules"
	"github.com/twitchtv/twirp"
)

// aggregator for incremental file loading:
type moduleFiles struct {
	gomod string
	gosum string
}

func (s *Update) Files(_ context.Context, req *dependabot_v1.FilesRequest) (*dependabot_v1.FilesResponse, error) {
	// Update incoming files:
	for path, data := range req.GetFiles() {
		switch path {
		case modules.GoMod:
			s.files.gomod = string(data)
		case modules.GoSum:
			s.files.gosum = string(data)
		}
	}

	// Determine missing files:
	var res dependabot_v1.FilesResponse
	if len(s.files.gomod) == 0 {
		res.RequiredPaths = append(res.RequiredPaths, modules.GoMod)
	}
	if len(s.files.gosum) == 0 {
		res.OptionalPaths = append(res.OptionalPaths, modules.GoSum)
	}
	return &res, nil
}

func (s *Update) modules() (*modules.Modules, twirp.Error) {
	if len(s.files.gomod) == 0 {
		return nil, twirp.NewError(twirp.FailedPrecondition, "Files()")
	}
	return modules.NewModules(s.files.gomod, s.files.gosum), nil
}
