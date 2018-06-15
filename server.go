package drift

import (
	"github.com/mayur-tolexo/drift/lib"
	"github.com/rightjoin/aqua"
)

//DS is the drift service
type ds struct {
	aqua.RestService `prefix:"drift" root:"/" version:"1"`
	addConsumer      aqua.POST `url:"add/consumer/"`
	drift            *drift
}

//AddConsumer will add new consumer to the given topic
func (d *ds) AddConsumer(req aqua.Aide) (int, interface{}) {
	var (
		data    interface{}
		payload AddConstumer
		err     error
	)
	if payload, err = vAddConsumer(req); err == nil {
		data, err = d.drift.pAddConsumer(payload)
	}
	return lib.BuildResponse(data, err)
}
