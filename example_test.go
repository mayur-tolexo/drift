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

	//elastic v6.2 handler is printIT2
	d.AddChanelHandler("elastic", "v6.2", printIT2)

	//elastic all channels handler, except v6.2, is printIT3
	d.AddTopicHandler("elastic", printIT3)

	//port assign here is 1500
	d.Start(1500)
}
