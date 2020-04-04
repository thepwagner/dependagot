package service

import (
	"context"

	dependabot_v1 "github.com/github/dependabot/go/common/dependabot/v1"
	"github.com/github/dependabot/go/modules/modules"
	"github.com/twitchtv/twirp"
)

func (s *Update) Files(_ context.Context, req *dependabot_v1.FilesRequest) (*dependabot_v1.FilesResponse, error) {
	// Update incoming files:
	for path, data := range req.GetFiles() {
		switch path {
		case modules.GoMod:
			s.modules.GoMod = string(data)
		case modules.GoSum:
			s.modules.GoSum = string(data)
		default:
			s.modules.Paths[path] = string(data)
		}
	}

	var res dependabot_v1.FilesResponse

	// If go.mod hasn't been provided, it's required:
	if len(s.modules.GoMod) == 0 {
		res.RequiredPaths = append(res.RequiredPaths, modules.GoMod)
	} else {
		parsedPaths, err := s.modules.AdditionalPaths()
		if err != nil {
			return nil, err
		}
		res.RequiredPaths = append(res.RequiredPaths, parsedPaths...)
	}

	// go.sum is always optional:
	if len(s.modules.GoSum) == 0 {
		res.OptionalPaths = append(res.OptionalPaths, modules.GoSum)
	}
	return &res, nil
}

func (s *Update) ensureModules() twirp.Error {
	if len(s.modules.GoMod) == 0 {
		return twirp.NewError(twirp.FailedPrecondition, "Files()")
	}
	return nil
}
