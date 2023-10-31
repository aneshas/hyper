package main

import (
	"github.com/aneshas/hyper"
	"log"
	"os"
	"path"
)

func main() {
	loc, _ := os.Getwd()
	appPath := path.Join(loc, "tmp")
	mod := "github.com/hyper/webapp"

	// err := os.RemoveAll(appPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// err = hyper.InitApp(appPath, "webapp", mod)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// // TODO - This will be run from inside project so path, name, mod would be parsed
	// // err = hyper.GenUC(appPath, "task", "ScheduleTask", "", "")
	// err = hyper.GenUC(appPath, "webapp", "task", mod, "ScheduleTask", "UnscheduledTask", "ScheduledTask")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// err = hyper.GenUC(appPath, "webapp", "task", mod, "CancelTask", "ToCancel", "")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// err = hyper.GenUC(appPath, "webapp", "employee", mod, "ListEmployees", "ListFilter", "Employees")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err := hyper.GenStore(appPath, "webapp", "task", mod)
	if err != nil {
		log.Fatal(err)
	}
}
