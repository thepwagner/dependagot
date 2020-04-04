package modules

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dependabot/gomodules-extracted/cmd/go/_internal_/modfile"
	"github.com/github/dependabot/go/common/dependabot/v1"
	"github.com/sirupsen/logrus"
)

const (
	GoMod = "go.mod"
	GoSum = "go.sum"
)

type Modules struct {
	GoMod string
	GoSum string
	Paths map[string]string
}

func NewModules(gomod, gosum string, paths map[string]string) *Modules {
	if paths == nil {
		paths = map[string]string{}
	}
	return &Modules{
		GoMod: gomod,
		GoSum: gosum,
		Paths: paths,
	}
}

func (s *Modules) Dependencies() ([]*dependabot_v1.Dependency, error) {
	mod, err := s.parsedModFile()
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

func (s *Modules) AdditionalPaths() ([]string, error) {
	mod, err := s.parsedModFile()
	if err != nil {
		return nil, err
	}

	// Search for modules mapped to local paths, which may be in a multi-module/monorepo
	// The client can decide to provide them or not:
	var res []string
	for _, replace := range mod.Replace {
		if newPath := replace.New.Path; newPath != "" {
			newGoMod := fmt.Sprintf("%s/%s", newPath, GoMod)
			if _, ok := s.Paths[newGoMod]; !ok {
				res = append(res, newGoMod)
			}
		}
	}
	return res, nil
}

func (s *Modules) DependencyVersion(dep *dependabot_v1.Dependency) (map[string]string, error) {
	return nil, nil
}

func (s *Modules) UpdateDependencies(deps []*dependabot_v1.Dependency) (map[string]string, error) {
	mod, err := s.parsedModFile()
	if err != nil {
		return nil, err
	}
	for _, dep := range deps {
		if err := mod.AddRequire(dep.Package, dep.Version); err != nil {
			return nil, err
		}
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
		GoSum: goSum,
	}, nil

}

func (s *Modules) parsedModFile() (*modfile.File, error) {
	return modfile.Parse("go.mod", []byte(s.GoMod), nil)
}

var (
	emptyMain = []byte("package main\nfunc main() {}\n")
)

func (s *Modules) generateGoSum(mod *modfile.File) (string, error) {
	gomod, err := mod.Format()
	if err != nil {
		return "", err
	}

	dir, err := ioutil.TempDir("", "gomod-*")
	if err != nil {
		return "", fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(dir)

	// Nest the sandbox deep enough that relative paths are contained:
	var maxDepth int
	for path := range s.Paths {
		depth := len(strings.Split(path, "../"))
		if depth > maxDepth {
			maxDepth = depth
		}
	}
	logrus.WithField("depth", maxDepth).Debug("calculated max depth")
	if maxDepth > 1 {
		dir = fmt.Sprintf("%s/%s", dir, strings.Repeat("f/", maxDepth-1))
		if err := os.MkdirAll(dir, 0700); err != nil {
			return "", err
		}
	}

	// Write files to temporary directory:
	gomodFile := filepath.Join(dir, GoMod)
	if err := ioutil.WriteFile(gomodFile, gomod, 0400); err != nil {
		return "", err
	}
	gosumFile := filepath.Join(dir, GoSum)
	if len(s.GoSum) > 0 {
		if err := ioutil.WriteFile(gosumFile, []byte(s.GoSum), 0600); err != nil {
			return "", err
		}
	}
	mainGoFile := filepath.Join(dir, "main.go")
	if err := ioutil.WriteFile(mainGoFile, emptyMain, 0400); err != nil {
		return "", err
	}
	logrus.WithFields(logrus.Fields{
		"go.mod": gomodFile,
		"go.sum": gosumFile,
	}).Debug("Wrote files to sandbox")
	for path, data := range s.Paths {
		fullPath := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0700); err != nil {
			return "", err
		}
		logrus.WithField("path", fullPath).Debug("Wrote additional sandbox file")
		if err := ioutil.WriteFile(fullPath, []byte(data), 0400); err != nil {
			return "", err
		}
	}

	var buf bytes.Buffer
	cmd := exec.Command("go", "get", "-d", "-v")
	cmd.Env = []string{
		"GO111MODULE=on",
		fmt.Sprintf("GOCACHE=%s", filepath.Join(dir, "cache")),
		fmt.Sprintf("GOPATH=%s", filepath.Join(dir, "go")),
	}
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "PATH=") {
			cmd.Env = append(cmd.Env, env)
		}
	}

	cmd.Dir = dir
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		fmt.Println(buf.String())
		return "", err
	}
	fmt.Println(buf.String())

	gosum, err := ioutil.ReadFile(gosumFile)
	if err != nil {
		return "", err
	}
	return string(gosum), nil
}
