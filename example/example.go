package main

import (
	"fmt"

	"github.com/mayur-tolexo/drift"
)

func printIT(value ...interface{}) error {
	fmt.Println(value)
	return nil
}

func main() {
	d := drift.Newdrift(printIT)
	d.Start()
}
