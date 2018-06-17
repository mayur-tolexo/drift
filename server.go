package drift

import (
	"strings"

	"github.com/mayur-tolexo/drift/lib"
	"github.com/rightjoin/aqua"
)

//DS is the drift service
type ds struct {
	aqua.RestService `prefix:"drift" root:"/" version:"1"`
	consumerCount    aqua.GET  `url:"consumer/"`
	stopAdmin        aqua.GET  `url:"stop/admin/"`
	startAdmin       aqua.POST `url:"start/admin/"`
	addConsumer      aqua.POST `url:"add/consumer/"`
	publishReq       aqua.POST `url:"pub/request/"`
	killConsumer     aqua.POST `url:"kill/consumer/"`
	admin            aqua.POST `url:"admin/"`
	drift            *Drift
}

//AddConsumer will add new consumer to the given topic
func (d *ds) AddConsumer(req aqua.Aide) (int, interface{}) {
	var (
		data    interface{}
		payload AddConstumer
		err     error
	)
	if payload, err = vAddConsumer(req); err == nil {
		data, err = d.drift.addConsumer(payload)
	}
	return lib.BuildResponse(data, err)
}

//PublishReq will publish request
func (d *ds) PublishReq(req aqua.Aide) (int, interface{}) {
	var (
		data    interface{}
		payload Publish
		err     error
	)
	if payload, err = vPublishReq(req); err == nil {
		data, err = pPublishReq(payload)
	}
	return lib.BuildResponse(data, err)
}

//ConsumerCount will return the consumer count of the channel of given topic
//pass topic and channel in query params.
//Is need to get count of total consumer of a topic then only pass topic
func (d *ds) ConsumerCount(qParam aqua.Aide) (int, interface{}) {
	count := 0
	qParam.LoadVars()
	topic := qParam.QueryVar["topic"]
	channel := qParam.QueryVar["channel"]
	if channel == "" {
		for key, val := range d.drift.consumers {
			if strings.HasPrefix(key, topic) {
				count += len(val)
			}
		}
	} else {
		count = len(d.drift.consumers[hash(topic, channel)])
	}
	return lib.BuildResponse(count, nil)
}

//KillConsumer will kill consumer of given topic
func (d *ds) KillConsumer(req aqua.Aide) (int, interface{}) {
	var (
		data    interface{}
		payload KillConsumer
		err     error
	)
	if payload, err = vKillConsumer(req); err == nil {
		data, err = d.drift.killConsumer(payload)
	}
	return lib.BuildResponse(data, err)
}

//StartAdmin will start admin
func (d *ds) StartAdmin(req aqua.Aide) (int, interface{}) {
	var (
		data interface{}
		err  error
	)
	if err = d.drift.admin.vStartAdmin(req); err == nil {
		if d.drift.admin.adminRunning {
			err = lib.VError("Already Running at", d.drift.admin.httpAddrs)
		} else {
			go d.drift.admin.startAdmin()
			data = "Admin started at " + d.drift.admin.httpAddrs
		}
	}
	return lib.BuildResponse(data, err)
}

//StopAdmin will stop admin
func (d *ds) StopAdmin(req aqua.Aide) (int, interface{}) {
	var data interface{}
	if d.drift.admin.adminRunning {
		d.drift.admin.exitAdmin <- 1
		<-d.drift.admin.exitAdmin
		data = "DONE"
	} else {
		data = "Not Running"
	}
	return lib.BuildResponse(data, nil)
}

//Admin will do the admin actions
func (d *ds) Admin(req aqua.Aide) (int, interface{}) {
	var (
		data    interface{}
		payload Admin
		err     error
	)
	if payload, err = vAdmin(req); err == nil {
		data, err = d.drift.admin.doAction(payload)
	}
	return lib.BuildResponse(data, err)
}
