package service

import (
	"context"

	dependabot_v1 "github.com/github/dependabot/go/common/dependabot/v1"
)

type Update struct {
}

var _ dependabot_v1.UpdateServiceServer = (*Update)(nil)

func NewUpdate() *Update {
	return &Update{}
}

func (s *Update) Files(ctx context.Context, req *dependabot_v1.FilesRequest) (*dependabot_v1.FilesResponse, error) {
	return &dependabot_v1.FilesResponse{
		Paths: []string{"go", "server"},
	}, nil
}
