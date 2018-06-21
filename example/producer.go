package main

import (
	"fmt"

	"flag"

	"github.com/mayur-tolexo/drift"
)

func main() {
	msg := flag.String("msg", "Hi this is a test", "Message to broadcast")
	flag.Parse()
	nsqdTCPAddrs := []string{"127.0.0.1:4150"}
	d := drift.NewPub(nsqdTCPAddrs)
	topic := "elastic"
	if resp, err := d.Publish(topic, *msg); err == nil {
		fmt.Println(resp)
	} else {
		fmt.Println(err.Error())
	}
}
