package hyper

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"path"
)

func NewUC(NS string, name string, req string, resp string, appDir string) UC {
	return UC{NS: NS, Name: name, Req: req, Resp: resp, AppDir: appDir}
}

type UC struct {
	NS     string
	Name   string
	Req    string
	Resp   string
	AppDir string
}

func (uc UC) Dir() string {
	return path.Join(uc.AppDir, "pkg", uc.NS)
}

func (uc UC) Path() string {
	return path.Join(uc.AppDir, fmt.Sprintf("%s.go", strcase.ToSnake(uc.Name)))
}
