package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	options "github.com/mreiferson/go-options"
	"github.com/nsqio/nsq/nsqlookupd"
)

func nsqlookupdFlagSet(opts *nsqlookupd.Options) *flag.FlagSet {
	flagSet := flag.NewFlagSet("nsqlookupd", flag.ExitOnError)

	flagSet.String("log-level", "info", "set log verbosity: debug, info, warn, error, or fatal")
	flagSet.String("log-prefix", "[nsqlookupd] ", "log message prefix")
	flagSet.Bool("verbose", false, "deprecated in favor of log-level")

	flagSet.String("tcp-address", opts.TCPAddress, "<addr>:<port> to listen on for TCP clients")
	flagSet.String("http-address", opts.HTTPAddress, "<addr>:<port> to listen on for HTTP clients")
	flagSet.String("broadcast-address", opts.BroadcastAddress, "address of this lookupd node, (default to the OS hostname)")

	flagSet.Duration("inactive-producer-timeout", opts.InactiveProducerTimeout, "duration of time a producer will remain in the active list since its last ping")
	flagSet.Duration("tombstone-lifetime", opts.TombstoneLifetime, "duration of time a producer will remain tombstoned if registration remains")

	return flagSet
}

func main() {
	exitChan := make(chan int)
	signalChan := make(chan os.Signal, 1)
	go func() {
		<-signalChan
		exitChan <- 1
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	opts := nsqlookupd.NewOptions()
	flagSet := nsqlookupdFlagSet(opts)
	flagSet.Parse(os.Args[1:])

	options.Resolve(opts, flagSet, nil)
	daemon := nsqlookupd.New(opts)

	err := daemon.Main()
	if err != nil {
		log.Fatalf("ERROR: failed to start nsqlookupd: %v", err)
	}

	<-exitChan
	daemon.Exit()
}
