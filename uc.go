package hyper

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"path"
)

func GenUC(location string, app string, ns string, mod string, uc string, req string, resp string) error {
	appDir := path.Join(location, app)
	ucDir := path.Join(appDir, "pkg", ns)

	f := jen.NewFile(ns)

	// TODO - If req or resp exist write info log and skip

	if req != "" {
		f.Comment(fmt.Sprintf("%s represents ...", req))
		f.Type().Id(req).Struct(
			jen.Id("Foo").String(),
		)
	}

	if resp != "" {
		f.Line()
		f.Comment(fmt.Sprintf("%s represents ...", resp))
		f.Type().Id(resp).Struct(
			jen.Id("Bar").String(),
		)
	}

	f.Line()

	params := []jen.Code{
		jen.Id("ctx").Qual("context", "Context"),
	}

	if req != "" {
		params = append(params, jen.Id("cmd").Id(req))
	}

	var ret []jen.Code

	if resp != "" {
		ret = append(ret, jen.Id(fmt.Sprintf("*%s", resp)))
	}

	ret = append(ret, jen.Error())

	ucf := fmt.Sprintf("%sFunc", uc)

	f.Comment(fmt.Sprintf("%s accomplishes ...", ucf))

	f.Type().
		Id(ucf).
		Func().
		Params(params...).
		Params(ret...)

	f.Line()

	ucn := fmt.Sprintf("New%s", uc)

	f.ImportName("github.com/aneshas/tx", "")

	rets := []jen.Code{
		jen.Nil(),
	}

	if resp != "" {
		rets = append(rets, jen.Nil())
	}

	f.Comment(fmt.Sprintf("%s instantiates %s use case", ucn, uc))
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

	err := mkDir(ucDir)
	if err != nil {
		return err
	}

	err = f.Save(path.Join(ucDir, fmt.Sprintf("%s.go", strcase.ToSnake(uc))))
	if err != nil {
		return err
	}

	err = genUCHandler(location, app, ns, mod, uc, req, resp)
	if err != nil {
		return err
	}

	return modTidy(appDir)
}

func genUCHandler(location string, app string, ns string, mod string, uc string, req string, resp string) error {
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

	err := mkDir(hDir)
	if err != nil {
		return err
	}

	return f.Save(path.Join(hDir, fmt.Sprintf("%s.go", strcase.ToSnake(uc))))
}
