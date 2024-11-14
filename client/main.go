package main

import (
	"fmt"
	"net"
	// "github.com/scrmbld/battlesloop-go/sloopGame"
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

	fmt.Printf("%s\n", buf)

	// 3. play a game with the server

	// 4. clean up
	conn.Close()
}
