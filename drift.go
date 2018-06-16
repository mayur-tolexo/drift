package drift

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/rightjoin/aqua"
	"github.com/tolexo/aero/conf"
	"github.com/tolexo/tachyon/lib"
)

//AddTopicHandler will add a new handler with the given topic
func (d *drift) AddTopicHandler(topic string, jobHandler JobHandler) {
	d.chanelHandler[hash(topic, allKey)] = jobHandler
}

//AddChanelHandler will add a new handler with the channel of given topic
func (d *drift) AddChanelHandler(topic, channel string, jobHandler JobHandler) {
	d.chanelHandler[hash(topic, channel)] = jobHandler
}

//Start will start the drift server
func (d *drift) Start(port int) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	d.Server = aqua.NewRestServer()
	d.Server.AddModule("access", aqua.ModAccessLog(conf.String("drift.access", "")))
	d.Server.Modules = "access"
	d.Server.Port = lib.GetPriorityValue(port, conf.Int("drift.port", 0)).(int)
	d.Server.AddService(&ds{drift: d})
	go d.SysInterrupt()
	d.Server.Run()
}

//SysInterrupt will handle system interrupt
func (d *drift) SysInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTSTP)
	fmt.Println("System Exit: ", <-c)
	for _, topic := range d.consumers {
		for _, consumer := range topic {
			consumer.Stop()
		}
	}
	for _, topic := range d.consumers {
		for _, consumer := range topic {
			<-consumer.StopChan
		}
	}
	os.Exit(1)
}
