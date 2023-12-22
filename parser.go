package hyper

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path"
)

type Parser interface {
	// TODO - return Store object with name etc ...

	HasStore(app App, uc UC) (bool, error)
}

func NewGoParser() *GoParser {
	return &GoParser{}
}

type GoParser struct{}

// TODO
// if interface present - add it to newly generated use cases
// if present AND embeds hyper.Store then also generate the store

func (g *GoParser) HasStore(app App, uc UC) (bool, error) {
	data, err := os.ReadFile(path.Join(app.Location, app.NameOnDisk, "pkg", uc.NS, uc.NS+".go"))
	if err != nil {
		_, ok := err.(*fs.PathError)
		if ok {
			return false, nil
		}

		return false, err
	}

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "entity", data, parser.ParseComments)
	if err != nil {
		return false, err
	}

	hasStore := false

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

		hasStore = true

		_ = i

		return true
	})

	return hasStore, nil
}
