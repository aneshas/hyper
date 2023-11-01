package hyper

import (
	"path"
)

type AppStore interface {
	CurrentApp() (*App, error)
}

func NewApp(location string, name string, mod string) App {
	return App{Location: location, NameOnDisk: name, Mod: mod}
}

type App struct {
	// Optional
	Location   string
	NameOnDisk string
	Mod        string
}

func (a App) Dir() string {
	return path.Join(a.Location, a.NameOnDisk)
}

func (a App) CMDDir() string {
	return path.Join(a.Dir(), "cmd", a.NameOnDisk)
}

func (a App) NewUC(ns string, name string, req string, resp string) UC {
	return NewUC(ns, name, req, resp, a.Dir())
}
