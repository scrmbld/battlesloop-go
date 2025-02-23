package main

import (
	"bufio"
	"fmt"
	"github.com/scrmbld/battlesloop-go/sloopNet"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	// connect to the server
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the ip address of the server you would like to connect to:")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	server_addr := strings.Trim(text, "\n")

	connection := sloopNet.GameConn{}
	err = connection.Connect(server_addr, "8080")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Quit()

	// decide what we want to do
	fmt.Println("Would you like to play against a computer (1) or against another player (2)?")
	text, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	choice, err := strconv.Atoi(strings.Trim(text, "\n"))
	if err != nil {
		fmt.Println(err)
		return
	}

	// tell the server our choice and begin the game (when the time comes)
	if choice == 1 {
		err = connection.SendMsg("_s_computer:")
		if err != nil {
			fmt.Println(err)
			return
		}
		playGame(&connection, true)
	} else if choice == 2 {
		err = connection.SendMsg("_s_player:")
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

		if msg[0] == "s" {
			parts := strings.Split(msg[1], "-")

			if parts[0] == "l" {
				// wait for the opponent to connect to us
				l, err := net.Listen("tcp", ":8081")
				if err != nil {
					fmt.Println(err)
					return
				}
				// tell server that we are listening
				connection.SendMsg("_s_ready:")

				sock, err := l.Accept()
				if err != nil {
					fmt.Println(err)
					return
				}

				var opponent sloopNet.GameConn
				opponent.Init(sock)
				err = playGame(&opponent, true)
				if err != nil {
					fmt.Println(err)
					return
				}

			} else { // connect to the opponent
				// wait for server to tell us opponent is ready
				err = connection.ReadMsg()
				if err != nil {
					fmt.Println(err)
					return
				}
				msg, err = connection.PopMsg()
				if msg[0] != "s" && msg[1] != "ready" {
					fmt.Println("Error: opponent did not ready up")
				}

				// connect to opponent
				var opponent sloopNet.GameConn
				err = opponent.Connect(parts[1], ":8081")
				if err != nil {
					fmt.Println(err)
					return
				}

				err = playGame(&opponent, false)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}

		fmt.Println("Game over")
	} else {
		fmt.Println("Error: value out of range")
		return
	}

}
