package service

import (
	"context"

	dependabot_v1 "github.com/github/dependabot/go/common/dependabot/v1"
)

type Update struct {
	files moduleFiles
}

var _ dependabot_v1.UpdateService = (*Update)(nil)

func NewUpdate() *Update {
	return &Update{}
}

func (s *Update) ListDependencies(context.Context, *dependabot_v1.ListDependenciesRequest) (*dependabot_v1.ListDependenciesResponse, error) {
	mod, twirpErr := s.modules()
	if twirpErr != nil {
		return nil, twirpErr
	}

	deps, err := mod.Dependencies()
	if err != nil {
		return nil, err
	}
	return &dependabot_v1.ListDependenciesResponse{
		Dependencies: deps,
	}, nil
}
