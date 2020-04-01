package runner

import (
	"context"
	"fmt"

	"github.com/github/dependabot/go/common/dependabot/v1"
	"google.golang.org/grpc"
)

type Client struct {
	dependabot_v1.UpdateServiceClient
	conn *grpc.ClientConn
}

func NewClient(ctx context.Context, addr string) (*Client, error) {
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("dialing: %w", err)
	}
	return &Client{
		conn:                conn,
		UpdateServiceClient: dependabot_v1.NewUpdateServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}
