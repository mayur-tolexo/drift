package drift

import nsq "github.com/nsqio/go-nsq"

//AddConstumer is the request format of add consumer
type AddConstumer struct {
	LookupAddr   []string    `json:"loopup_address"`
	NsqDTCPAddrs []string    `json:"nsqd_tcp_address"`
	Topic        []TopicData `json:"topic_detail"`
	MaxInFlight  int         `json:"max_in_flight"`
}

//TopicData will contail the topic details
type TopicData struct {
	Topic   string `json:"topic"`
	Channel string `json:"channel"`
}

//TailHandler will implement the nsq handler
type TailHandler struct {
	topicName  string
	jobHandler JobHandler
}

//JobHandler function which will be called
type JobHandler func(value ...interface{}) error

//Drift will have the handler function
type Drift struct {
	jobHandler JobHandler
	consumers  []*nsq.Consumer
}
