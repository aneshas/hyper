package main

import (
	"context"
	"github.com/aneshas/hyper"
	"log"
	"os"
)

func main() {
	mod := "github.com/hyper/webapp"

	err := os.RemoveAll("webapp")
	if err != nil {
		log.Println(err)
	}

	ctx := context.Background()

	fs := new(hyper.DiskFS)
	tools := new(hyper.GoCLI)
	gen := hyper.NewFuncGen(fs, tools)

	initApp := hyper.NewInitApp(fs, gen, tools)
	genUC := hyper.NewGenUC(fs, gen)

	err = initApp(ctx, hyper.AppDetails{
		Name: "webapp",
		Mod:  mod,
	})

	if err != nil {
		log.Fatal(err)
	}

	_, err = genUC(ctx, hyper.UCDetails{
		NS:     "task",
		UCName: "ScheduleTask",
		Req:    "UnscheduledTask",
		Resp:   "ScheduledTask",
	})
	if err != nil {
		log.Fatal(err)
	}
}
