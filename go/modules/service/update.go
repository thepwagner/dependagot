package service

import (
	"context"

	dependabot_v1 "github.com/github/dependabot/go/common/dependabot/v1"
)

type Update struct {
}

var _ dependabot_v1.UpdateService = (*Update)(nil)

func NewUpdate() *Update {
	return &Update{}
}

func (s *Update) Files(ctx context.Context, req *dependabot_v1.FilesRequest) (*dependabot_v1.FilesResponse, error) {
	return &dependabot_v1.FilesResponse{
		RequiredPaths: []string{"go", "server"},
	}, nil
}

func (s *Update) Setup(ctx context.Context, req *dependabot_v1.SetupRequest) (*dependabot_v1.SetupResponse, error) {
	sandboxID := "test"
	return &dependabot_v1.SetupResponse{
		SandboxId: sandboxID,
	}, nil
}
