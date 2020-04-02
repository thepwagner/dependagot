package sandbox

import (
	"github.com/dependabot/gomodules-extracted/cmd/go/_internal_/modfile"
	"github.com/dependabot/gomodules-extracted/cmd/go/_internal_/modload"
	dependabot_v1 "github.com/github/dependabot/go/common/dependabot/v1"
)

type Sandbox struct {
	gomod []byte
	gosum []byte
}

func (s *Sandbox) RequiredPaths() []string {
	if len(s.gomod) == 0 {
		return []string{"go.mod"}
	}
	return nil
}

func (s *Sandbox) OptionalPaths() []string {
	if len(s.gosum) == 0 {
		return []string{"go.sum"}
	}
	return nil
}

func (s *Sandbox) File(path string, data []byte) {
	switch path {
	case "go.mod":
		s.gomod = data
	case "go.sum":
		s.gosum = data
	}
}

func (s *Sandbox) Dependencies() ([]*dependabot_v1.Dependency, error) {
	mod, err := modfile.Parse("go.mod", s.gomod, nil)
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

	modload.InitMod()

	return deps, nil
}
