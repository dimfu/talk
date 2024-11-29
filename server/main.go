package main

import (
	"log"
	"net"
)

type msgpayload struct {
	clientid *int
	message  string
}

func main() {
	listener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	hub := newHub()
	go hub.run()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("%v", err)
		}

		c := newClient(
			conn,
			hub.commands,
			hub.registrations,
			hub.deregistrations,
		)

		go c.read()
	}
}
