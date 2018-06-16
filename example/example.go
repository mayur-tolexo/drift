package main

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

func main() {
	//Default handler is printIT
	d := drift.Newdrift(printIT)
	//elastic v6.2 handler is printIT2
	d.AddChanelHandler("elastic", "v6.2", printIT2)
	//elastic all channels, except v6.2, is printIT3
	d.AddTopicHandler("elastic", printIT3)
	d.Start(1500)
}
