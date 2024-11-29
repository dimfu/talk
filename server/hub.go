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
		if client.id != from.id {
			msg := append([]byte(fmt.Sprintf("Client %d: ", from.id)), m...)
			_, err = conn.Write(msg)
		}

		if err != nil {
			log.Println(err)
		}
	}
}

func (h *hub) register(c *client) {
	clientCounts := len(h.clients)
	if _, exists := h.clients[c.conn]; !exists {
		log.Printf("accepted new connection from %s", c.conn.RemoteAddr().String())
		c.id = clientCounts
		h.clients[c.conn] = c
	}
}

func (h *hub) deregister(c *client) {
	if _, exists := h.clients[c.conn]; exists {
		log.Printf("%s is disconnected", c.conn.RemoteAddr().String())
		delete(h.clients, c.conn)
	}
}

func (h *hub) message(c *client, m []byte) {
	if sender, ok := h.clients[c.conn]; ok {
		h.broadcast(sender, m)
	}
}
