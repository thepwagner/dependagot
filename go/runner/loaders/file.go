package loaders

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/github/dependabot/go/runner/runner"
)

type File struct {
	root string
}

var _ runner.Loader = (*File)(nil)

func NewFile(root string) *File {
	return &File{root: root}
}

func (f *File) Load(_ context.Context, path string) ([]byte, bool, error) {
	p := filepath.Join(f.root, path)
	file, err := ioutil.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return file, true, err
}
