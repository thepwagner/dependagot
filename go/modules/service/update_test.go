package service_test

import (
	"context"
	"testing"

	"github.com/github/dependabot/go/common/dependabot/v1"
	"github.com/github/dependabot/go/modules/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdate_Files(t *testing.T) {
	u := service.NewUpdate()
	ctx := context.Background()
	files, err := u.Files(ctx, &dependabot_v1.FilesRequest{})
	require.NoError(t, err)
	assert.Equal(t, []string{"go", "server"}, files.GetPaths())
}
