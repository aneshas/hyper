package hyper

import (
	"github.com/dave/jennifer/jen"
	"os"
	"os/exec"
	"path"
)

func InitApp(location, name, modPath string) error {
	appDir := path.Join(location, name)

	err := mkDir(appDir)
	if err != nil {
		return err
	}

	err = initMod(modPath, appDir)
	if err != nil {
		return err
	}

	cmdDir := path.Join(appDir, "cmd", name)

	err = mkDir(cmdDir)
	if err != nil {
		return err
	}

	return writeMain(path.Join(cmdDir, "main.go"))
}

func mkDir(p string) error {
	return os.MkdirAll(p, 0755)
}

func initMod(modPath string, appDir string) error {
	cmd := exec.Command("go", "mod", "init", modPath)

	cmd.Dir = appDir
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func modTidy(appDir string) error {
	cmd := exec.Command("go", "mod", "tidy")

	cmd.Dir = appDir
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func writeMain(p string) error {
	echo := "github.com/labstack/echo/v4"
	tx := "github.com/aneshas/tx"

	f := jen.NewFile("main")

	f.ImportName(echo, "echo")
	f.ImportName(tx, "")

	f.Func().
		Id("main").
		Params().
		Block(
			jen.Id("e").Op(":=").Qual(echo, "New").Call(),
			f.Line(),
			jen.Var().Id("conn").Op("*").Qual("database/sql", "DB"),
			f.Line(),
			jen.Id("transactor").Op(":=").Qual(tx, "New").Call(jen.Id("conn")),
			f.Line(),
			jen.Id("_").Op("=").Id("transactor"),
			f.Line(),
			jen.Qual("log", "Fatal").Call(jen.Id("e").Dot("Start").Call(jen.Lit(":8080"))),
		)

	return f.Save(p)
}
