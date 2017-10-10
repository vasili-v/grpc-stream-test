package main

import "flag"

var (
	address     string
	limit       int
	compression bool
)

func init() {
	flag.StringVar(&address, "a", ":5555", "address:port for requests")
	flag.IntVar(&limit, "l", 500, "max number of requests to handle in parallel")
	flag.BoolVar(&compression, "c", false, "enable stream compression")

	flag.Parse()
}
