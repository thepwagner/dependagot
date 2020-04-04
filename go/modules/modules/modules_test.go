package modules_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	dependabot_v1 "github.com/github/dependabot/go/common/dependabot/v1"
	"github.com/github/dependabot/go/modules/modules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSandbox_Dependencies(t *testing.T) {
	m := modules.NewModules(fixture(t, "testify151.mod"), "")

	deps, err := m.Dependencies()
	require.NoError(t, err)
	if assert.Len(t, deps, 1) {
		testifyDep := deps[0]
		assert.Equal(t, "github.com/stretchr/testify", testifyDep.Package)
		assert.Equal(t, "v1.5.1", testifyDep.Version)
	}
}

func TestSandbox_Upgrade(t *testing.T) {
	m := modules.NewModules(
		fixture(t, "testify150.mod"),
		fixture(t, "testify150.sum"),
	)

	files, err := m.DependencyVersion(&dependabot_v1.Dependency{
		Package: "github.com/stretchr/testify",
		Version: "v1.5.1",
	})
	require.NoError(t, err)

	if assert.Len(t, files, 2) {
		newGoMod := files[modules.GoMod]
		assert.Equal(t, fixture(t, "testify151.mod"), newGoMod)

		newGoSum := files[modules.GoSum]
		assert.Equal(t, fixture(t, "testify150-upgraded.sum"), newGoSum)
	}
}

func fixture(t *testing.T, fixture string) string {
	t.Helper()
	dir, err := os.Getwd()
	require.NoError(t, err)
	data, err := ioutil.ReadFile(filepath.Join(dir, "fixtures", fixture))
	require.NoError(t, err)
	return string(data)
}
