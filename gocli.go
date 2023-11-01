package hyper

import (
	"os"
	"os/exec"
)

type GoTools interface {
	ModInit(app App) error
	ModTidy(app App) error
}

type GoCLI struct {
}

func (gc *GoCLI) ModInit(app App) error {
	cmd := exec.Command("go", "mod", "init", app.Mod)

	cmd.Dir = app.Dir()
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func (gc *GoCLI) ModTidy(app App) error {
	cmd := exec.Command("go", "mod", "tidy")

	cmd.Dir = app.Dir()
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
