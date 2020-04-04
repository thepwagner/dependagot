package service_test

import (
	"context"
	"testing"

	"github.com/github/dependabot/go/common/dependabot/v1"
	"github.com/github/dependabot/go/modules/modules"
	"github.com/github/dependabot/go/modules/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdate_Files(t *testing.T) {
	ctx := context.Background()

	t.Run("empty", func(t *testing.T) {
		u := service.NewUpdate()
		files, err := u.Files(ctx, &dependabot_v1.FilesRequest{})
		require.NoError(t, err)
		assert.Equal(t, []string{modules.GoMod}, files.GetRequiredPaths())
		assert.Equal(t, []string{modules.GoSum}, files.GetOptionalPaths())
	})

	t.Run("with go.sum", func(t *testing.T) {
		u := service.NewUpdate()
		files, err := u.Files(ctx, &dependabot_v1.FilesRequest{
			Files: map[string][]byte{
				modules.GoSum: make([]byte, 1),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, []string{modules.GoMod}, files.GetRequiredPaths())
		assert.Empty(t, files.GetOptionalPaths())
	})

	t.Run("with go.mod and go.sum", func(t *testing.T) {
		u := service.NewUpdate()
		files, err := u.Files(ctx, &dependabot_v1.FilesRequest{
			Files: map[string][]byte{
				modules.GoMod: make([]byte, 1),
				modules.GoSum: make([]byte, 1),
			},
		})
		require.NoError(t, err)
		assert.Empty(t, files.GetRequiredPaths())
		assert.Empty(t, files.GetOptionalPaths())
	})
}
