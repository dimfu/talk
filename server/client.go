package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

type client struct {
	conn       net.Conn
	id         int
	outbound   chan<- command
	register   chan<- *client
	deregister chan<- *client
}

func newClient(conn net.Conn, o chan<- command, r chan<- *client, d chan<- *client) *client {
	c := &client{
		conn:       conn,
		outbound:   o,
		register:   r,
		deregister: d,
	}
	r <- c
	return c
}

func (c *client) read() error {
	for {
		inc, err := bufio.NewReader(c.conn).ReadString('\n')
		if err == io.EOF {
			c.deregister <- c
			return nil
		}

		if err != nil {
			return err
		}

		c.handle(inc)
	}
}

func (c *client) handle(inc string) {
	msg := append([]byte(fmt.Sprintf("[Client %d] ", c.id)), inc...)
	c.outbound <- command{
		id:     MSG,
		sender: fmt.Sprint(c.id),
		body:   msg,
		client: c,
	}
}
