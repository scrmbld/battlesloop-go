package main

import (
	"fmt"
	"github.com/scrmbld/battlesloop-go/sloopNet"
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

		handleJoin(conn)
	}
}

func handleJoin(sock net.Conn) {
	var connection sloopNet.GameConn
	connection.Init(sock)
	defer connection.Quit()

	// wait for client to tell us what they want to do
	err := connection.ReadMsg()
	if err != nil {
		fmt.Println(err)
		return
	}

	msg, err := connection.PopMsg()
	if err != nil {
		fmt.Println(err)
		return
	}

	if msg[0] == "s" && msg[1] == "computer" {
		// begin match vs. computer
		err := connection.SendMsg("_g_begin:")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = connection.ReadMsg()
		if err != nil {
			fmt.Println(err)
			return
		}
		msg, err = connection.PopMsg()
		if msg[0] != "g" && msg[1] != "begin" {
			fmt.Println("Opponent did not confirm game start")
			return
		}

		err = startGame(&connection)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("game over")
		}

	} else {
		fmt.Println("Invalid choice from client")
	}

}
