package main

import (
	"fmt"
	"net"
	// "github.com/scrmbld/battlesloop-go/sloopGame"
	"github.com/scrmbld/battlesloop-go/sloopNet"
)

func main() {

	// 1. connect to server
	conn, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		fmt.Println(err)
		return
	}

	// 2. wait for server to send a "start game" message
	buf := make([]byte, 512)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	var msg = string(buf[:])
	t, err := sloopNet.GetType(msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	if t == "c" {
		conn_msg, err := sloopNet.ParseConn(msg)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	// 3. play a game with the server

	// 4. clean up
	conn.Close()
}
