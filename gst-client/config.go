package main

import "flag"

var (
	server  string
	total   int
	streams int
	size    int
	random  bool
)

func init() {
	flag.StringVar(&server, "s", ":5555", "address:port of server")
	flag.IntVar(&total, "n", 5, "number of requests to send")
	flag.IntVar(&streams, "c", 1, "number of parallel gRPC streams")
	flag.IntVar(&size, "size", 100, "payload size")
	flag.BoolVar(&random, "r", false, "enable random payload")

	flag.Parse()
}
