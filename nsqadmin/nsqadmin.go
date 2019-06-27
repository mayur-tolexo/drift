package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/nsqio/nsq/nsqadmin"
)

//StringArray is a wrapper of string array
type StringArray []string

//Set will append the new value in string array
func (a *StringArray) Set(s string) error {
	*a = append(*a, s)
	return nil
}

//String will join string values
func (a *StringArray) String() string {
	return strings.Join(*a, ",")
}

func nsqAdminFlagSet(opts *nsqadmin.Options) {
	lookupdHTTPAddrs := flag.String("lookupd-http-address", "", "lookupd HTTP address (may be given multiple times)")
	httpAddrs := flag.String("http-address", "localhost:4171", "<addr>:<port> to listen on for HTTP clients")
	flag.Parse()
	opts.NSQLookupdHTTPAddresses = strings.Split(*lookupdHTTPAddrs, ",")
	opts.HTTPAddress = *httpAddrs
}

func main() {
	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	opts := nsqadmin.NewOptions()
	nsqAdminFlagSet(opts)
	daemon := nsqadmin.New(opts)
	daemon.Main()
	<-exitChan
	daemon.Exit()
}
