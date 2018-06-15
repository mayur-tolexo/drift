package drift

import (
	"runtime"

	"github.com/rightjoin/aqua"
	"github.com/tolexo/aero/conf"
)

//Newdrift will create the drift model
func Newdrift(jh JobHandler) *drift {
	return &drift{
		jobHandler: jh,
	}
}

//Start will start the drift server
func (d *drift) Start() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	server := aqua.NewRestServer()
	server.AddModule("access", aqua.ModAccessLog(conf.String("drift.access", "")))
	server.Modules = "access"
	server.Port = conf.Int("drift.port", 0)
	server.AddService(&DS{drift: d})
	go d.SystemInterrupt()
	server.Run()
}
