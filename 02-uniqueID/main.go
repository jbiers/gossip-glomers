package main

import (
	"fmt"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// TODO: implement logging
// TODO: improve variable naming, clean up code

type Server struct {
	node   *maelstrom.Node
	id     uint64
	idLock sync.Mutex // this controls thread sync when running with multiple threads
}

func main() {
	s := &Server{
		node: maelstrom.NewNode(),
		id:   0,
	}

	s.node.Handle("generate", s.generateHandler)

	if err := s.node.Run(); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) generateHandler(msg maelstrom.Message) error {
	s.idLock.Lock()
	res := map[string]any{
		"type": "generate_ok",
		"id":   fmt.Sprintf("%v%v", s.id, encodeDecimal(s.node.ID())),
	}
	s.id++
	s.idLock.Unlock()

	return s.node.Reply(msg, res)
}
