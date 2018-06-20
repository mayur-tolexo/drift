package drift

import (
	"fmt"
	"github.com/mayur-tolexo/drift"
)

// new pub created to publish message to nsqd
func ExampleNewPub() {
	d := drift.NewPub("127.0.0.1:4151")
	msg := "This is a test"
	if resp, err := d.Publish("elastic", msg); err == nil {
		fmt.Println(resp)
	} else {
		fmt.Println(err.Error())
	}
}
