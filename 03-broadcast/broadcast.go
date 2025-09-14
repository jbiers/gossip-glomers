package main

import maelstrom "github.com/jepsen-io/maelstrom/demo/go"

const (
	BroadcastType   = "broadcast"
	BroadcastOKType = "broadcast_ok"
	ReadType        = "read"
	ReadOKType      = "read_ok"
	TopologyType    = "topology"
	TopologyOKType  = "topology_ok"
)

type BroadcastBody struct {
	maelstrom.MessageBody
	Message *float64 `json:"message,omitempty"`
}

type ReadBody struct {
	maelstrom.MessageBody
	Messages []float64 `json:"messages,omitempty"`
}

type TopologyBody struct {
	maelstrom.MessageBody
	Topology map[string][]string `json:"topology,omitempty"`
}
