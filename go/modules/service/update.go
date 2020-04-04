package service

import (
	"context"

	dependabot_v1 "github.com/github/dependabot/go/common/dependabot/v1"
	"github.com/github/dependabot/go/modules/modules"
)

type Update struct {
	modules *modules.Modules
}

var _ dependabot_v1.UpdateService = (*Update)(nil)

func NewUpdate() *Update {
	return &Update{
		modules: modules.NewModules("", "", nil),
	}
}

func (s *Update) ListDependencies(_ context.Context, _ *dependabot_v1.ListDependenciesRequest) (*dependabot_v1.ListDependenciesResponse, error) {
	if twirpErr := s.ensureModules(); twirpErr != nil {
		return nil, twirpErr
	}

	deps, err := s.modules.Dependencies()
	if err != nil {
		return nil, err
	}
	return &dependabot_v1.ListDependenciesResponse{
		Dependencies: deps,
	}, nil
}

func (s *Update) UpdateDependencies(_ context.Context, req *dependabot_v1.UpdateDependenciesRequest) (*dependabot_v1.UpdateDependenciesResponse, error) {
	if twirpErr := s.ensureModules(); twirpErr != nil {
		return nil, twirpErr
	}

	files, err := s.modules.UpdateDependencies(req.Dependencies)
	if err != nil {
		return nil, err
	}
	return &dependabot_v1.UpdateDependenciesResponse{
		NewFiles: files,
	}, nil
}
