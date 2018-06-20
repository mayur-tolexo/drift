package drift

import (
	"flag"
	"fmt"

	"github.com/mayur-tolexo/drift"
)

// new pub created to publish message to nsqd
func ExampleNewPub() {
	msg := flag.String("msg", "Hi this is a test", "Message to broadcast")
	flag.Parse()
	d := drift.NewPub("127.0.0.1:4151")
	if resp, err := d.Publish("elastic", *msg); err == nil {
		fmt.Println(resp)
	} else {
		fmt.Println(err.Error())
	}
}

// This will map a new handeler with specified topic's channel
func ExampleAddChanelHandler() {
	d := drift.NewConsumer(nil)
	topic := "elastic"
	channel := "v6.2"
	d.AddChanelHandler(topic, channel, printIT)
}

// This will map a new handeler with all channels of the specified topic.
// If a channelHandler is already mapped with any channel of the specified topic then that handler will be called
// and in rest of the channel this handler will be called.
func ExampleAddTopicHandler() {
	d := drift.NewConsumer(nil)
	topic := "elastic"
	d.AddChanelHandler(topic, printIT)
}

// This will start the drift server to receive request over HTTP
func ExampleStart() {
	d := drift.NewConsumer(printIT)
	port := 1500
	d.Start(port)
}
