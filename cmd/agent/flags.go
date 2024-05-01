package main

import (
	"flag"
	"os"
	"time"
)

var flagRunAddr string

var options struct {
	pollInterval   time.Duration
	reportInterval time.Duration
}

func parseFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flag.StringVar(&flagRunAddr, "a", "localhost:8080", "server address and port")
	p := flag.Int64("p", 2, "frequency of polling metrics from the runtime package")
	r := flag.Int64("r", 10, "frequency of sending metrics to the server")
	flag.Parse()
	options.pollInterval = time.Duration(*p) * time.Second
	options.reportInterval = time.Duration(*r) * time.Second
}
