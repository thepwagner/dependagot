package modules

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dependabot/gomodules-extracted/cmd/go/_internal_/modfile"
	"github.com/github/dependabot/go/common/dependabot/v1"
)

const (
	GoMod = "go.mod"
	GoSum = "go.sum"
)

type Modules struct {
	GoMod []byte
	GoSum []byte
}

func NewModules(gomod, gosum string) *Modules {
	return &Modules{
		GoMod: []byte(gomod),
		GoSum: []byte(gosum),
	}
}

func (s *Modules) Dependencies() ([]*dependabot_v1.Dependency, error) {
	mod, err := modfile.Parse("go.mod", s.GoMod, nil)
	if err != nil {
		return nil, err
	}

	// TODO: what about replacements?
	deps := make([]*dependabot_v1.Dependency, 0, len(mod.Require))
	for _, r := range mod.Require {
		deps = append(deps, &dependabot_v1.Dependency{
			Package: r.Mod.Path,
			Version: r.Mod.Version,
		})
	}
	return deps, nil
}

func (s *Modules) DependencyVersion(dep *dependabot_v1.Dependency) (map[string]string, error) {
	mod, err := modfile.Parse("go.mod", s.GoMod, nil)
	if err != nil {
		return nil, err
	}
	if err := mod.AddRequire(dep.Package, dep.Version); err != nil {
		return nil, err
	}

	goSum, err := s.generateGoSum(mod)
	if err != nil {
		return nil, err
	}
	goMod, err := mod.Format()
	if err != nil {
		return nil, err
	}

	return map[string]string{
		GoMod: string(goMod),
		GoSum: string(goSum),
	}, nil
}

var (
	emptyMain = []byte("package main\nfunc main() {}\n")
)

func (s *Modules) generateGoSum(mod *modfile.File) ([]byte, error) {
	gomod, err := mod.Format()
	if err != nil {
		return nil, err
	}

	dir, err := ioutil.TempDir("", "gomod-*")
	if err != nil {
		return nil, fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(dir)

	// Write files to temporary directory:
	gomodFile := filepath.Join(dir, GoMod)
	if err := ioutil.WriteFile(gomodFile, gomod, 0400); err != nil {
		return nil, err
	}
	gosumFile := filepath.Join(dir, GoSum)
	if len(s.GoSum) > 0 {
		if err := ioutil.WriteFile(gosumFile, s.GoSum, 0600); err != nil {
			return nil, err
		}
	}
	mainGoFile := filepath.Join(dir, "main.go")
	if err := ioutil.WriteFile(mainGoFile, emptyMain, 0400); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	cmd := exec.Command("go", "get", "-d", "-v")
	cmd.Env = []string{
		"GO111MODULE=on",
		fmt.Sprintf("GOCACHE=%s", filepath.Join(dir, "cache")),
		fmt.Sprintf("GOPATH=%s", filepath.Join(dir, "go")),
	}
	cmd.Dir = dir
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		fmt.Println(buf.String())
		return nil, err
	}
	fmt.Println(buf.String())

	gosum, err := ioutil.ReadFile(gosumFile)
	if err != nil {
		return nil, err
	}
	return gosum, nil
}
