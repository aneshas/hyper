package hyper

import (
	"context"
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"path"
)

const (
	echoPkg = "github.com/labstack/echo/v4"
	txPkg   = "github.com/aneshas/tx"
)

func NewFuncGen(fs FS, parser Parser, tools GoTools) *FuncGen {
	return &FuncGen{fs: fs, tools: tools, parser: parser}
}

type FuncGen struct {
	fs     FS
	tools  GoTools
	parser Parser
}

func (fg *FuncGen) GenMain(app App) error {
	f := jen.NewFile("main")

	f.ImportName(echoPkg, "echo")
	f.ImportName(txPkg, "")

	f.Func().
		Id("main").
		Params().
		Block(
			jen.Id("e").Op(":=").Qual(echoPkg, "New").Call(),
			f.Line(),
			jen.Var().Id("conn").Op("*").Qual("database/sql", "DB"),
			f.Line(),
			jen.Id("transactor").Op(":=").Qual(txPkg, "New").Call(jen.Id("conn")),
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

	f.ImportName(txPkg, "")

	rets := []jen.Code{
		jen.Nil(),
	}

	if uc.Resp != "" {
		rets = append(rets, jen.Nil())
	}

	ucParams := []jen.Code{
		jen.Id("tx").Qual(txPkg, "Transactor"),
	}

	has, err := fg.parser.HasStore(app, uc)
	if err != nil {
		return err
	}

	if has {
		// TODO - Use name from the store details

		ucParams = append(
			ucParams,
			jen.Id("store").Id("Store"),
		)
	}

	f.Comment(fmt.Sprintf("%s instantiates %s use case", ucn, uc.Name))
	f.Func().
		Id(ucn).
		Params(ucParams...).
		Params(jen.Id(ucf)).
		Block(
			jen.Return(
				jen.Func().Params(params...).Params(ret...).
					Block(
						jen.Return(rets...),
					),
			),
		)

	err = fg.fs.MkAppDir(ucDir)
	if err != nil {
		return err
	}

	err = f.Save(path.Join(ucDir, fmt.Sprintf("%s.go", strcase.ToSnake(uc.Name))))
	if err != nil {
		return err
	}

	err = fg.genUCHandler(app, uc)
	if err != nil {
		return err
	}

	return fg.tools.ModTidy(app)
}

func (fg *FuncGen) genUCHandler(app App, uc UC) error {
	hDir := path.Join(app.Dir(), "internal", "echo")

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

	nsPkg := fmt.Sprintf("%s/pkg/%s", app.Mod, uc.NS)

	f.ImportName(echoPkg, "echo")

	hn := fmt.Sprintf("Reg%s", uc.Name)
	hnn := fmt.Sprintf("New%s", uc.Name)

	f.Func().Id(hn).
		Params(
			jen.Id("e").Op("*").Qual(echoPkg, "Echo"),
			jen.Id(strcase.ToLowerCamel(uc.Name)).Qual(nsPkg, fmt.Sprintf("%sFunc", uc.Name)),
			jen.Id("mw").Op("...").Qual(echoPkg, "MiddlewareFunc"),
		).Block(
		jen.Id("e").Dot("GET").Call(
			jen.Lit(fmt.Sprintf("%s/%s", uc.NS, strcase.ToLowerCamel(uc.Name))),
			jen.Id(hnn).Call(
				jen.Id(strcase.ToLowerCamel(uc.Name)),
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
			jen.Id(strcase.ToLowerCamel(uc.Name)).Qual(nsPkg, fmt.Sprintf("%sFunc", uc.Name)),
		).
		Params(
			jen.Qual(echoPkg, "HandlerFunc"),
		).Block(
		jen.Return(
			jen.Func().Params(
				jen.Id("c").Qual(echoPkg, "Context"),
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

	return f.Save(path.Join(hDir, fmt.Sprintf("%s.go", strcase.ToSnake(uc.Name))))
}
