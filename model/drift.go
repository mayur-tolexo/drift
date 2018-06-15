package model

//AddConstumer is the request format of add consumer
type AddConstumer struct {
	LookupAddr string      `json:"loopup_address"`
	Topic      []TopicData `json:"topic_detail"`
}

//TopicData will contail the topic details
type TopicData struct {
	Topic   string `json:"topic"`
	Channel string `json:"channel"`
}
