package loaders

import (
	"context"

	"github.com/github/dependabot/go/runner/runner"
)

type Memory struct {
	data map[string][]byte
}

var _ runner.Loader = (*Memory)(nil)

func NewMemory(data map[string][]byte) *Memory {
	return &Memory{data: data}
}

func (m *Memory) Load(_ context.Context, path string) ([]byte, bool, error) {
	if data, ok := m.data[path]; ok {
		return data, true, nil
	}
	return nil, false, nil
}
