package common

import (
	"flag"
)

// Sarama configuration options
var (
	numChannels = 16
	port        = 8000
)

func init() {
	flag.IntVar(&numChannels, "numChannels", 16, "Num of Channels")
	flag.IntVar(&port, "port", 8000, "Port to run at")
	flag.Parse()
}
