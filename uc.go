package hyper

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
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

func GenStore(location string, app string, ns string, mod string) error {
	data, err := os.ReadFile(path.Join(location, app, "pkg", ns, ns+".go"))
	if err != nil {
		// if errors.Is(err, fs.PathError) {
		// 	return nil
		// }

		return err
	}

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "entity", data, parser.ParseComments)
	if err != nil {
		return err
	}

	ast.Inspect(file, func(x ast.Node) bool {
		s, ok := x.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if s.Type == nil {
			return true
		}

		i, ok := s.Type.(*ast.InterfaceType)
		if !ok {
			return true
		}

		if s.Name.Name != "Store" {
			return true
		}

		_ = i

		return true
	})

	return nil
}
