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
	Count   int    `json:"count"`
}

//AddAdmin is the add admin request
type AddAdmin struct {
	AdminUser                []string `json:"user"`
	HTTPAddrs                string   `json:"http_address"`
	LookupHTTPAddr           []string `json:"lookup_http_address"`
	NsqDTCPAddrs             []string `json:"nsqd_tcp_address"`
	ACLHTTPHeader            string   `json:"acl_http_header"`
	NotificationHTTPEndpoint string   `json:"notification_http_endpoint"`
}

//Admin is the request format of admin api to permorm action.
//allowed actions are - create/empty/delete/pause/unpause
type Admin struct {
	Topic   string `json:"topic"`
	Channel string `json:"channel"`
	Action  string `json:"action"`
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
	admin         dAdmin
}

//dAdmin have the drift admin model
type dAdmin struct {
	adminUser                []string
	httpAddrs                string
	lookupHTTPAddr           []string
	nsqDTCPAddrs             []string
	aclHTTPHeader            string
	notificationHTTPEndpoint string
	adminRunning             bool
	exitAdmin                chan int
}
