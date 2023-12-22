package hyper

import (
	"os"
)

type FS interface {
	MkAppDir(p string) error
}

type DiskFS struct {
}

func (d *DiskFS) CurrentApp() (*App, error) {
	// TODO - This will be run from inside project so path, name, mod would be parsed
	// loc, _ := os.Getwd()
	// appPath := path.Join(loc, "tmp")
	mod := "github.com/wisag/lucie"

	app := NewApp("", "lucie", mod)

	return &app, nil
}

func (d *DiskFS) MkAppDir(p string) error {
	return os.MkdirAll(p, 0755)
}
