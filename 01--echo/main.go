package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Server struct {
	node *maelstrom.Node
}

func main() {
	s := &Server{
		node: maelstrom.NewNode(),
	}

	s.node.Handle("echo", s.echoHandler)

	if err := s.node.Run(); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) echoHandler(msg maelstrom.Message) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	body["type"] = "echo_ok"
	return s.node.Reply(msg, body)
}
