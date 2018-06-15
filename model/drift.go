package model

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
