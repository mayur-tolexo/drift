package service

import (
	"github.com/mayur-tolexo/drift/lib"
	"github.com/rightjoin/aqua"
)

//Drift service
type Drift struct {
	aqua.RestService `prefix:"tune" root:"/" version:"1"`
	sellerMySQLToPG  aqua.POST `url:"add/consumer/"`
}

//AddConsumer will add new consumer to the given topic
func (*Drift) AddConsumer(req aqua.Aide) (int, interface{}) {
	var (
		data interface{}
		// payload model.AddConstumer
		err error
	)
	return lib.BuildResponse(data, err)
}
