package main

import (
	"encoding/json"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"github.com/sirupsen/logrus"
)

type NodeID string
type NodeStorageVersion float64

type Server struct {
	node                  maelstrom.Node
	storage               []float64
	storageIndex          map[float64]bool
	storageLock           sync.Mutex
	neighbors             map[NodeID]NodeStorageVersion
	logger                logrus.Logger
	currentStorageVersion float64
}

func main() {
	s := &Server{
		node:                  *maelstrom.NewNode(),
		storage:               []float64{},
		storageIndex:          map[float64]bool{},
		neighbors:             map[NodeID]NodeStorageVersion{},
		logger:                *logrus.New(),
		currentStorageVersion: 0,
	}

	// how are these handlers implemented? are they async?
	s.node.Handle("broadcast", s.broadcastHandler)
	s.node.Handle("read", s.readHandler)
	s.node.Handle("topology", s.topologyHandler)

	if err := s.node.Run(); err != nil {
		s.logger.Fatal(err)
	}

}

func (s *Server) sendBroadcast(dstNodeID NodeID, message float64) {
	body := &BroadcastBody{
		Type: BroadcastType,
	}

	body.Message = &message
	err := s.node.RPC(string(dstNodeID), body, s.broadcastOKHandler)

	if err != nil {
		s.logger.Error(err)
	}
}

func (s *Server) broadcastOKHandler(msg maelstrom.Message) error {
	return nil
}

func (s *Server) broadcastHandler(msg maelstrom.Message) error {
	defer s.storageLock.Unlock()

	var body BroadcastBody

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	s.storageLock.Lock()

	if !s.storageIndex[*body.Message] {
		s.storage = append(s.storage, *body.Message)
		s.storageIndex[*body.Message] = true

		for neighborID := range s.neighbors {
			if neighborID != NodeID(msg.Src) {
				s.sendBroadcast(neighborID, *body.Message)
			}
		}
	}

	s.currentStorageVersion = *body.Message

	body.Type = BroadcastOKType
	body.Message = nil

	return s.node.Reply(msg, body)
}

func (s *Server) readHandler(msg maelstrom.Message) error {

	var body ReadBody

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	body.Type = ReadOKType
	body.Messages = &s.storage

	return s.node.Reply(msg, body)
}

func (s *Server) topologyHandler(msg maelstrom.Message) error {
	var body TopologyBody

	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	neighbors := body.Topology[s.node.ID()]

	for _, neighbor := range neighbors {
		s.neighbors[NodeID(neighbor)] = 0
	}

	body.Type = TopologyOKType
	body.Topology = nil

	return s.node.Reply(msg, body)
}
