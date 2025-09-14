package main

import (
	"encoding/json"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Server struct {
	node        maelstrom.Node
	storage     []float64
	storageLock sync.RWMutex
	neighbors   []string
}

func main() {

	s := &Server{
		node:      *maelstrom.NewNode(),
		storage:   []float64{},
		neighbors: []string{},
	}

	// goroutine that every 1 second broadcasts all of the cluster storage

	s.node.Handle("broadcast", s.broadcastHandler)
	s.node.Handle("read", s.readHandler)
	s.node.Handle("topology", s.topologyHandler)

	if err := s.node.Run(); err != nil {
		log.Fatal(err)
	}

}

func (s *Server) broadcastHandler(msg maelstrom.Message) error {
	defer s.storageLock.Unlock()

	var body BroadcastBody

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.storageLock.Lock()
	s.storage = append(s.storage, *body.Message)

	body.Type = BroadcastOKType
	body.Message = nil

	return s.node.Reply(msg, body)
}

func (s *Server) readHandler(msg maelstrom.Message) error {
	defer s.storageLock.RUnlock()

	var body ReadBody

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.storageLock.RLock()
	body.Type = ReadOKType
	body.Messages = s.storage

	return s.node.Reply(msg, body)
}

func (s *Server) topologyHandler(msg maelstrom.Message) error {
	var body TopologyBody

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.neighbors = body.Topology[s.node.ID()]

	body.Type = TopologyOKType
	body.Topology = nil

	return s.node.Reply(msg, body)
}
