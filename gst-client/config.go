package main

import (
	"flag"
	"time"
)

var (
	server      string
	total       int
	size        int
	timeout     time.Duration
	compression bool
	random      bool
)

func init() {
	flag.StringVar(&server, "s", ":5555", "address:port of server")
	flag.IntVar(&total, "n", 5, "number of requests to send")
	flag.IntVar(&size, "size", 60, "payload size")
	flag.DurationVar(&timeout, "t", 2*time.Minute, "time to wait for responses")
	flag.BoolVar(&compression, "c", false, "enable stream compression")
	flag.BoolVar(&random, "r", false, "enable random payload")

	flag.Parse()
}
