package update

import (
	"context"

	dependabot_v1 "github.com/github/dependabot/go/common/dependabot/v1"
)

type Service struct {
}

var _ dependabot_v1.UpdateServiceServer = (*Service)(nil)

func NewService() *Service {
	return &Service{}
}

func (s *Service) Files(ctx context.Context, req *dependabot_v1.FilesRequest) (*dependabot_v1.FilesResponse, error) {
	return &dependabot_v1.FilesResponse{
		Paths: []string{"go", "server"},
	}, nil
}
