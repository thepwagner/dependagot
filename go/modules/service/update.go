package service

import (
	"context"

	"github.com/thepwagner/dependagot/go/common/dependagot/v1"
	"github.com/thepwagner/dependagot/go/modules/modules"
)

type Update struct {
	modules *modules.Modules
}

var _ dependagot_v1.UpdateService = (*Update)(nil)

func NewUpdate() *Update {
	return &Update{
		modules: modules.NewModules("", "", nil),
	}
}

func (s *Update) ListDependencies(_ context.Context, _ *dependagot_v1.ListDependenciesRequest) (*dependagot_v1.ListDependenciesResponse, error) {
	if twirpErr := s.ensureModules(); twirpErr != nil {
		return nil, twirpErr
	}

	deps, err := s.modules.Dependencies()
	if err != nil {
		return nil, err
	}
	return &dependagot_v1.ListDependenciesResponse{
		Dependencies: deps,
	}, nil
}

func (s *Update) UpdateDependencies(_ context.Context, req *dependagot_v1.UpdateDependenciesRequest) (*dependagot_v1.UpdateDependenciesResponse, error) {
	if twirpErr := s.ensureModules(); twirpErr != nil {
		return nil, twirpErr
	}

	files, err := s.modules.UpdateDependencies(req.Dependencies)
	if err != nil {
		return nil, err
	}
	return &dependagot_v1.UpdateDependenciesResponse{
		NewFiles: files,
	}, nil
}
