package main

import "flag"

var (
	address    string
	maxStreams int
)

func init() {
	flag.StringVar(&address, "a", ":5555", "address:port for requests")
	flag.IntVar(&maxStreams, "s", 0, "limit on the number of concurrent streams (0 - use gRPC default)")

	flag.Parse()
}
