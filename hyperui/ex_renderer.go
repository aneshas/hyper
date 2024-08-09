package hyperui

import (
	"fmt"
	"github.com/aneshas/hyper/hyperui/flash"
	"io"

	"github.com/Masterminds/sprig/v3"
	"github.com/dannyvankooten/extemplate"
	"github.com/labstack/echo/v4"
)

func NewRenderer(tplPath string) (*Renderer, error) {
	xt := extemplate.New().
		Funcs(sprig.FuncMap()).
		Funcs(flash.FuncMap())

	err := xt.ParseDir(tplPath, []string{".tpl"})
	if err != nil {
		return nil, err
	}

	return &Renderer{
		tpl: xt,
	}, nil
}

type Renderer struct {
	tpl *extemplate.Extemplate
}

func (r *Renderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	bag, ok := data.(Bag)
	if ok {
		bag.Context = c
		data = bag.WithKV("Page", name)
	}

	return r.tpl.ExecuteTemplate(w, fmt.Sprintf("%s.go.tpl", name), data)
}
