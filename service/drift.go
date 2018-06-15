package service

import (
	"github.com/mayur-tolexo/drift/lib"
	"github.com/mayur-tolexo/drift/model"
	"github.com/mayur-tolexo/drift/util"
	"github.com/rightjoin/aqua"
)

//Drift service
type Drift struct {
	aqua.RestService `prefix:"drift" root:"/" version:"1"`
	addConsumer      aqua.POST `url:"add/consumer/"`
}

//AddConsumer will add new consumer to the given topic
func (*Drift) AddConsumer(req aqua.Aide) (int, interface{}) {
	var (
		data    interface{}
		payload model.AddConstumer
		err     error
	)
	if payload, err = util.VAddConsumer(req); err == nil {
		data, err = util.PAddConsumer(payload)
	}
	return lib.BuildResponse(data, err)
}
