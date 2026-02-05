//go:build dev
// +build dev

package frontend

import (
	"io/fs"
)

// DistDirFS is a no-op filesystem for development mode
var DistDirFS fs.FS = noopFS{}

type noopFS struct{}

func (noopFS) Open(name string) (fs.File, error) {
	return nil, fs.ErrNotExist
}
