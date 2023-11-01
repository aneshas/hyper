package hyper

import (
	"context"
)

// AppDetails represents ...
type AppDetails struct {
	Location string
	Name     string
	Mod      string
}

// InitAppFunc accomplishes ...
type InitAppFunc func(ctx context.Context, cmd AppDetails) error

type AppGenerator interface {
	GenMain(app App) error
}

// NewInitApp instantiates InitApp use case
func NewInitApp(fs FS, appGen AppGenerator, tools GoTools) InitAppFunc {
	return func(ctx context.Context, cmd AppDetails) error {
		app := NewApp(cmd.Location, cmd.Name, cmd.Mod)

		err := fs.MkAppDir(app.Dir())
		if err != nil {
			return err
		}

		err = tools.ModInit(app)
		if err != nil {
			return err
		}

		err = fs.MkAppDir(app.CMDDir())
		if err != nil {
			return err
		}

		return appGen.GenMain(app)
	}
}
