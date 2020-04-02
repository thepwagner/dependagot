package service

import (
	"context"

	dependabot_v1 "github.com/github/dependabot/go/common/dependabot/v1"
	"github.com/github/dependabot/go/modules/sandbox"
)

type Update struct {
	sandbox *sandbox.Sandbox
}

var _ dependabot_v1.UpdateService = (*Update)(nil)

func NewUpdate() *Update {
	return &Update{
		sandbox: &sandbox.Sandbox{},
	}
}
func (s *Update) Files(_ context.Context, req *dependabot_v1.FilesRequest) (*dependabot_v1.FilesResponse, error) {
	for path, data := range req.GetFiles() {
		s.sandbox.File(path, data)
	}
	return &dependabot_v1.FilesResponse{
		RequiredPaths: s.sandbox.RequiredPaths(),
		OptionalPaths: s.sandbox.OptionalPaths(),
	}, nil
}

func (s *Update) ListDependencies(context.Context, *dependabot_v1.ListDependenciesRequest) (*dependabot_v1.ListDependenciesResponse, error) {
	deps, err := s.sandbox.Dependencies()
	if err != nil {
		return nil, err
	}
	return &dependabot_v1.ListDependenciesResponse{
		Dependencies: deps,
	}, nil
}
