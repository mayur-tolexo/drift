package drift

import (
	nsq "github.com/nsqio/go-nsq"
	"github.com/rightjoin/aqua"
)

//AddConstumer is the request format of add consumer
type AddConstumer struct {
	LookupAddr   []string    `json:"loopup_address"`
	NsqDTCPAddrs []string    `json:"nsqd_address"`
	Topic        []TopicData `json:"topic_detail"`
	MaxInFlight  int         `json:"max_in_flight"`
}

//Publish is the request format of publish request api
type Publish struct {
	NsqDHttpAddrs string      `json:"nsqd_address"`
	Topic         string      `json:"topic"`
	Data          interface{} `json:"data"`
}

//TopicData will contail the topic details
type TopicData struct {
	Topic   string `json:"topic"`
	Channel string `json:"channel"`
}

//tailHandler will implement the nsq handler
type tailHandler struct {
	topicName  string
	jobHandler JobHandler
}

//JobHandler function which will be called
type JobHandler func(value ...interface{}) error

//Drift will have the handler function
type Drift struct {
	Server        aqua.RestServer
	chanelHandler map[string]JobHandler
	jobHandler    JobHandler
	consumers     map[string][]*nsq.Consumer
	pubAddrs      string
}
