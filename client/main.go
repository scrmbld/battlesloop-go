package main

import (
	"fmt"
	// "github.com/scrmbld/battlesloop-go/sloopGame"
	"github.com/scrmbld/battlesloop-go/sloopNet"
)

func main() {

	// 1. connect to server
	connection := sloopNet.GameConn{}
	err := connection.Connect("localhost", "8080")

	if err != nil {
		fmt.Println(err)
		return
	}

	// 2. wait for server to send a "start game" message
	err = connection.ReadMsg()
	if err != nil {
		fmt.Println(err)
		connection.Quit()
		return
	}

	// 3. read the message
	msg, err := connection.PopMsg()
	if err != nil {
		fmt.Println(err)
		connection.Quit()
		return
	}

	fmt.Printf("%s\n", msg)

	connection.Quit()
}
