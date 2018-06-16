package drift

import (
	nsq "github.com/nsqio/go-nsq"
	"github.com/rightjoin/aqua"
)

//AddConstumer is the request format of add consumer api
type AddConstumer struct {
	LookupHTTPAddr []string    `json:"lookup_http_address"`
	NsqDTCPAddrs   []string    `json:"nsqd_tcp_address"`
	Topic          []TopicData `json:"topic_detail"`
	MaxInFlight    int         `json:"max_in_flight"`
}

//KillConsumer is the request format of kill consumer api
type KillConsumer struct {
	Topic   string `json:"topic"`
	Channel string `json:"channel"`
	Count   int    `json:"count"`
}

//Publish is the request format of publish request api
type Publish struct {
	NsqDHTTPAddrs string      `json:"nsqd_http_address"`
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
