package main

import (
	"fmt"
	// "github.com/scrmbld/battlesloop-go/sloopGame"
	"errors"
	"github.com/scrmbld/battlesloop-go/sloopNet"
	"net"
)

func startGame(connection *sloopNet.GameConn) error {
	// wait for client to ask to go first
	err := connection.ReadMsg()
	if err != nil {
		fmt.Println(err)
		return err
	}

	msg, err := connection.PopMsg()
	if err != nil {
		return err
	}

	// True if we are going first, false otherwise
	var goFirst bool
	err = nil
	if msg[0] == "g" && msg[1] == "first" {
		goFirst = false
		// tell the client that we are going last
		err = connection.SendMsg("_g_last:")
	} else if msg[0] == "g" && (msg[1] == "last" || msg[1] == "undecided") {
		goFirst = true
		// tell the client that we are going first
		err = connection.SendMsg("_g_first:")
	} else {
		return errors.New("Expected turn order message, received something else")
	}

	if err != nil {
		return err
	}

	fmt.Println(goFirst)

	return nil
}

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

	msg, err := connection.PopMsg()
	if err != nil {
		fmt.Println(err)
		return
	}
	if msg[0] != "g" {
		fmt.Println(msg)
		fmt.Print("Nope, not doing that right now")
		return
	} else {
		if msg[1] == "begin" {
			fmt.Println("Client has acknowledged our \"begin\"")
			err = startGame(&connection)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

}
