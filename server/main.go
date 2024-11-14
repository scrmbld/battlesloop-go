package main

import (
	"fmt"
	"net"
)

func main() {

	// listen for new connections from clients
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		// send the new client off to join a room
	}
}

func handleJoin(conn net.Conn) {
	// close the connection when we're done
	// this might not be appropriate to have here depending on architecture
	defer conn.Close()

	// play a game against an AI, which runs on the server
}
