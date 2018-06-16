package drift

import (
	"runtime"

	"github.com/rightjoin/aqua"
	"github.com/tolexo/aero/conf"
	"github.com/tolexo/tachyon/lib"
)

//AddTopicHandler will add a new handler with the given topic
func (d *Drift) AddTopicHandler(topic string, jobHandler JobHandler) {
	d.chanelHandler[hash(topic, allKey)] = jobHandler
}

//AddChanelHandler will add a new handler with the channel of given topic
func (d *Drift) AddChanelHandler(topic, channel string, jobHandler JobHandler) {
	d.chanelHandler[hash(topic, channel)] = jobHandler
}

//Start will start the drift server
func (d *Drift) Start(port int) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	d.Server = aqua.NewRestServer()
	d.Server.AddModule("access", aqua.ModAccessLog(conf.String("drift.access", "")))
	d.Server.Modules = "access"
	d.Server.Port = lib.GetPriorityValue(port, conf.Int("drift.port", 0)).(int)
	d.Server.AddService(&ds{drift: d})
	go d.sysInterrupt()
	d.Server.Run()
}

//Publish will broadcast the data to the nsqd
func (d *Drift) Publish(topic string, data interface{}) (resp interface{}, err error) {
	payload := Publish{
		NsqDHTTPAddrs: d.pubAddrs,
		Topic:         topic,
		Data:          data,
	}
	resp, err = pPublishReq(payload)
	return
}
