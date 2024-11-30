package main

import (
	"fmt"
	"log"
	"net"
)

type hub struct {
	clients         map[net.Conn]*client
	commands        chan command
	registrations   chan *client
	deregistrations chan *client
}

func newHub() *hub {
	return &hub{
		clients:         make(map[net.Conn]*client),
		commands:        make(chan command),
		registrations:   make(chan *client),
		deregistrations: make(chan *client),
	}
}

func (h *hub) run() {
	for {
		select {
		case client := <-h.registrations:
			h.register(client)
		case client := <-h.deregistrations:
			h.deregister(client)
		case cmd := <-h.commands:
			switch cmd.id {
			case MSG:
				h.message(cmd.client, cmd.body)
			default:
			}
		}
	}
}

func (h *hub) broadcast(from *client, m []byte) {
	for conn, client := range h.clients {
		var err error
		if from == nil || client.id != from.id {
			_, err = conn.Write(m)
		}

		if err != nil {
			log.Println(err)
		}
	}
}

func (h *hub) register(c *client) {
	clientCounts := len(h.clients)
	if _, exists := h.clients[c.conn]; !exists {
		c.id = clientCounts
		h.clients[c.conn] = c
		h.broadcast(c, []byte(fmt.Sprintf("> Client %d logged on\n", c.id)))
	}
}

func (h *hub) deregister(c *client) {
	if _, exists := h.clients[c.conn]; exists {
		delete(h.clients, c.conn)
		h.broadcast(c, []byte(fmt.Sprintf("> Client %d logged off\n", c.id)))
	}
}

func (h *hub) message(c *client, m []byte) {
	if sender, ok := h.clients[c.conn]; ok {
		h.broadcast(sender, m)
	}
}
