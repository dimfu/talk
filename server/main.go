package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type msgpayload struct {
	clientid *int
	message  string
}

func main() {
	clientcounts := 0
	listener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	clients := make(map[net.Conn]int)
	disconnected := make(chan net.Conn)
	newconn := make(chan net.Conn)
	payload := make(chan *msgpayload)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal(err)
			}
			newconn <- conn
		}
	}()

	for {
		select {
		case c := <-newconn:
			log.Printf("accepted new connection from %s", c.RemoteAddr().String())
			clients[c] = clientcounts
			clientcounts++

			go func(conn net.Conn, clientid int) {
				reader := bufio.NewReader(c)
				for {
					inc, err := reader.ReadString('\n')
					if err != nil {
						break
					}
					payload <- &msgpayload{
						clientid: &clientid,
						message:  fmt.Sprintf("Client %d > %s", clientid, inc),
					}
				}
				disconnected <- conn
			}(c, clients[c])

		case message := <-payload:
			for conn, clientid := range clients {
				go func(conn net.Conn, clientid int, payload *msgpayload) {
					if *message.clientid != clientid {
						_, err = conn.Write([]byte(payload.message))
					}
					if err != nil {
						log.Println("Write ERROR: ", err)
						disconnected <- conn
						return
					}
				}(conn, clientid, message)
			}

		case d := <-disconnected:
			go func() {
				clientid, exists := clients[d]
				if exists {
					message := &msgpayload{
						clientid: &clientid,
						message:  fmt.Sprintf("Client %d is disconnected\n", clientid),
					}
					payload <- message
					delete(clients, d)
					d.Close()
				}
			}()
		}
	}
}
