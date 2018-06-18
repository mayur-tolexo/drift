package main

import (
	"fmt"

	"flag"

	"github.com/mayur-tolexo/drift"
)

func main() {
	msg := flag.String("msg", "Hi this is a test", "Message to broadcast")
	flag.Parse()
	d := drift.NewPub("127.0.0.1:4151")
	if resp, err := d.Publish("elastic", *msg); err == nil {
		fmt.Println(resp)
	} else {
		fmt.Println(err.Error())
	}
}
