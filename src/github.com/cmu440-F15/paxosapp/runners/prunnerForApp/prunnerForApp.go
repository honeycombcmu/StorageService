package main

import (
	"flag"
	"github.com/cmu440-F15/paxosapp/paxos"
	"strings"
	"time"
)

var (
	ports      = flag.String("ports", "32767,32768,32769", "ports for all paxos nodes,split by comma")
	numRetries = flag.Int("retries", 5, "number of times a node should retry dialing another node")
)

func main() {
	flag.Parse()

	portStrings := strings.Split(*ports, ",")

	hostMap := make(map[int]string)
	for i, port := range portStrings {
		hostMap[i] = "localhost:" + port
	}
	numNodes := len(portStrings)
	for i := 0; i < numNodes; i++ {
		go paxos.NewPaxosNode(hostMap[i], hostMap, numNodes, i, *numRetries, false)
		time.Sleep(time.Second)
	}
	// Run the paxos node forever.
	select {}
}
