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

	_ = mod
	_ = initApp

	// err = initApp(ctx, hyper.AppDetails{
	// 	Name: "lucie",
	// 	Mod:  mod,
	// })
	//
	// if err != nil {
	// 	log.Fatal(err)
	// }

	_, err = genUC(ctx, hyper.UCDetails{
		NS:     "contract",
		UCName: "ContactInfo",
		Req:    "ContactInfoReq",
		Resp:   "ContactInfoResp",
	})
	if err != nil {
		log.Fatal(err)
	}
}
