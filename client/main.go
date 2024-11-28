package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

func Read(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Server closed its connection")
			os.Exit(0)
		}
		fmt.Print(str)
	}
}

func Write(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(conn)

	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		_, err = writer.WriteString(str)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		writer.Flush()
	}
}

var wg sync.WaitGroup

func main() {
	wg.Add(1)
	conn, err := net.Dial("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go Read(conn)
	go Write(conn)

	wg.Wait()
}
