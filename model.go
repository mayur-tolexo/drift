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
	StartAdmin     bool        `json:"start_admin"`
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

//AddAdmin is the add admin request
type AddAdmin struct {
	AdminUser      []string `json:"user"`
	HTTPAddrs      string   `json:"http_address"`
	LookupHTTPAddr []string `json:"lookup_http_address"`
	NsqDTCPAddrs   []string `json:"nsqd_tcp_address"`
}

//tailHandler will implement the nsq handler
type tailHandler struct {
	topicName  string
	jobHandler JobHandler
}

//JobHandler function which will be called
type JobHandler func(value ...interface{}) error

//Drift have the consumer/publisher model
type Drift struct {
	Server        aqua.RestServer
	chanelHandler map[string]JobHandler
	jobHandler    JobHandler
	consumers     map[string][]*nsq.Consumer
	pubAddrs      string
	admin         DAdmin
}

//DAdmin have the drift admin model
type DAdmin struct {
	adminUser      []string `json:"user"`
	httpAddrs      string   `json:"http_address"`
	lookupHTTPAddr []string `json:"lookup_http_address"`
	nsqDTCPAddrs   []string `json:"nsqd_tcp_address"`
	adminRunning   bool
	exitAdmin      chan int
}
