package hyper

import (
	"context"
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"path"
)

func NewFuncGen(fs FS, tools GoTools) *FuncGen {
	return &FuncGen{fs: fs, tools: tools}
}

type FuncGen struct {
	fs    FS
	tools GoTools
}

func (fg *FuncGen) GenMain(app App) error {
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

	p := path.Join(app.CMDDir(), "main.go")

	return f.Save(p)
}

func (fg *FuncGen) GenUC(ctx context.Context, app App, uc UC) error {
	ucDir := uc.Dir()

	f := jen.NewFile(uc.NS)

	// TODO - If uc.Req or uc.Resp exist write info log and skip

	if uc.Req != "" {
		f.Comment(fmt.Sprintf("%s represents ...", uc.Req))
		f.Type().Id(uc.Req).Struct(
			jen.Id("Foo").String(),
		)
	}

	if uc.Resp != "" {
		f.Line()
		f.Comment(fmt.Sprintf("%s represents ...", uc.Resp))
		f.Type().Id(uc.Resp).Struct(
			jen.Id("Bar").String(),
		)
	}

	f.Line()

	params := []jen.Code{
		jen.Id("ctx").Qual("context", "Context"),
	}

	if uc.Req != "" {
		params = append(params, jen.Id("uc").Id(uc.Req))
	}

	var ret []jen.Code

	if uc.Resp != "" {
		ret = append(ret, jen.Id(fmt.Sprintf("*%s", uc.Resp)))
	}

	ret = append(ret, jen.Error())

	ucf := fmt.Sprintf("%sFunc", uc.Name)

	f.Comment(fmt.Sprintf("%s accomplishes ...", ucf))

	f.Type().
		Id(ucf).
		Func().
		Params(params...).
		Params(ret...)

	f.Line()

	ucn := fmt.Sprintf("New%s", uc.Name)

	f.ImportName("github.com/aneshas/tx", "")

	rets := []jen.Code{
		jen.Nil(),
	}

	if uc.Resp != "" {
		rets = append(rets, jen.Nil())
	}

	f.Comment(fmt.Sprintf("%s instantiates %s use case", ucn, uc.Name))
	f.Func().
		Id(ucn).
		Params(
			jen.Id("tx").Qual("github.com/aneshas/tx", "Transactor"),
		).
		Params(jen.Id(ucf)).
		Block(
			jen.Return(
				jen.Func().Params(params...).Params(ret...).
					Block(
						jen.Return(rets...),
					),
			),
		)

	err := fg.fs.MkAppDir(ucDir)
	if err != nil {
		return err
	}

	err = f.Save(path.Join(ucDir, fmt.Sprintf("%s.go", strcase.ToSnake(uc.Name))))
	if err != nil {
		return err
	}

	err = fg.genUCHandler(app.Location, app.NameOnDisk, uc.NS, app.Mod, uc.Name, uc.Req, uc.Resp)
	if err != nil {
		return err
	}

	return fg.tools.ModTidy(app)
}

func (fg *FuncGen) genUCHandler(location string, app string, ns string, mod string, uc string, req string, resp string) error {
	appDir := path.Join(location, app)
	hDir := path.Join(appDir, "internal", "http")

	f := jen.NewFile("http")

	// TODO - Generating models optionally
	// if req != "" {
	// 	f.Comment(fmt.Sprintf("%s represents ...", req))
	// 	f.Type().Id(req).Struct(
	// 		jen.Id("Foo").String(),
	// 	)
	// }
	//
	// if resp != "" {
	// 	f.Line()
	// 	f.Comment(fmt.Sprintf("%s represents ...", resp))
	// 	f.Type().Id(resp).Struct(
	// 		jen.Id("Bar").String(),
	// 	)
	// }

	echo := "github.com/labstack/echo/v4"
	nsPkg := fmt.Sprintf("%s/pkg/%s", mod, ns)

	f.ImportName(echo, "echo")

	hn := fmt.Sprintf("Reg%s", uc)
	hnn := fmt.Sprintf("New%s", uc)

	f.Func().Id(hn).
		Params(
			jen.Id("e").Op("*").Qual(echo, "Echo"),
			jen.Id(strcase.ToLowerCamel(uc)).Qual(nsPkg, fmt.Sprintf("%sFunc", uc)),
			jen.Id("mw").Op("...").Qual(echo, "MiddlewareFunc"),
		).Block(
		jen.Id("e").Dot("GET").Call(
			jen.Lit(fmt.Sprintf("%s/%s", ns, strcase.ToLowerCamel(uc))),
			jen.Id(hnn).Call(
				jen.Id(strcase.ToLowerCamel(uc)),
			),
			jen.Id("mw").Op("..."),
		),
	)

	f.Line()

	// f.Comment(fmt.Sprintf("%s accomplishes ...", ucf))

	f.
		Func().
		Id(hnn).
		Params(
			jen.Id(strcase.ToLowerCamel(uc)).Qual(nsPkg, fmt.Sprintf("%sFunc", uc)),
		).
		Params(
			jen.Qual(echo, "HandlerFunc"),
		).Block(
		jen.Return(
			jen.Func().Params(
				jen.Id("c").Qual(echo, "Context"),
			).Params(jen.Error()).Block(
				jen.Return(
					jen.Id("c").Dot("String").Call(
						jen.Qual("net/http", "StatusOK"),
						jen.Lit("Ok"),
					),
				),
			),
		),
	)

	err := fg.fs.MkAppDir(hDir)
	if err != nil {
		return err
	}

	return f.Save(path.Join(hDir, fmt.Sprintf("%s.go", strcase.ToSnake(uc))))
}
