package main

import (
	"runtime"

	"github.com/mayur-tolexo/drift/service"
	"github.com/mayur-tolexo/drift/util"
	"github.com/rightjoin/aqua"
	"github.com/tolexo/aero/conf"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	server := aqua.NewRestServer()

	server.AddModule("access", aqua.ModAccessLog(conf.String("drift.access", "")))
	server.Modules = "access"

	server.Port = conf.Int("drift.port", 0)
	server.AddService(&service.Drift{})
	go util.SystemInterrupt()
	server.Run()
}
