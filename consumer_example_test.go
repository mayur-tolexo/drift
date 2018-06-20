package drift

import (
	"fmt"

	"github.com/mayur-tolexo/drift"
)

func printIT(value ...interface{}) error {
	fmt.Println("In 1st Print", value)
	return nil
}

func printIT2(value ...interface{}) error {
	fmt.Println("In 2nd Print", value)
	return nil
}

func printIT3(value ...interface{}) error {
	fmt.Println("In 3rd Print", value)
	return nil
}

// New consumer created with handel to call by the consumer.
// This will start new server to receive request over HTTP
func ExampleNewConsumer() {
	//Default handler is printIT
	d := drift.NewConsumer(printIT)

	// This will map a new handeler with specified topic's channel
	d.AddChanelHandler("elastic", "v6.2", printIT2)

	// This will map a new handeler with all channels of the specified topic.
	// If a channelHandler is already mapped with any channel of the specified topic then that handler will be called
	// and in rest of the channel this handler will be called.
	d.AddTopicHandler("elastic", printIT3)

	//port assign here is 1500
	d.Start(1500)
}
