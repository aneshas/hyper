package main

import (
	"context"
	"github.com/aneshas/hyper"
	"log"
)

func main() {
	mod := "github.com/wisag/lucie"

	var err error

	// err = os.RemoveAll("lucie")
	// if err != nil {
	// 	log.Println(err)
	// }

	ctx := context.Background()

	fs := new(hyper.DiskFS)
	tools := new(hyper.GoCLI)
	parser := hyper.NewGoParser()
	storeGen := hyper.NewBoilStoreGen(tools, fs)
	gen := hyper.NewFuncGen(fs, parser, tools)

	initApp := hyper.NewInitApp(fs, gen, tools)
	genUC := hyper.NewGenUC(fs, gen, storeGen)

	_ = initApp
	_ = mod

	// err = initApp(ctx, hyper.AppDetails{
	// 	Name: "lucie",
	// 	Mod:  mod,
	// })
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// TODO - Concentrate on generating the skeleton and making it extendable so
	// you can then only write modules for various adapters etc ...
	// - maybe in functional style / data oriented style eg. every module gets the whole app definition
	//   and modifies it as it sees fit
	// eg. save rendering for the end - only build a representation
	// don't focus on different types of use cases, make it opinionated instead focus on extendability with different adapters
	// databases etc ...

	// _, err = genUC(ctx, hyper.UCDetails{
	// 	NS:     "contract",
	// 	UCName: "ContactInfo",
	// 	Req:    "ContactInfoReq",
	// 	Resp:   "ContactInfoResp",
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	_, err = genUC(ctx, hyper.UCDetails{
		NS:     "contract",
		UCName: "RemoveContact",
		Req:    "ContactInfoReqB",
		Resp:   "ContactInfoRespC",
	})
	if err != nil {
		log.Fatal(err)
	}
}
