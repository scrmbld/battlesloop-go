package main

import (
	"fmt"
	"github.com/scrmbld/battlesloop-go/sloopGame"
	"github.com/scrmbld/battlesloop-go/sloopNet"
)

func startGame(connection *sloopNet.GameConn) error {
	var board sloopGame.Board
	board.PrintBoard()

	// say that we want to go first
	err := connection.SendMsg("_g_undecided:")
	if err != nil {
		fmt.Println(err)
		return err
	}

	// wait for the opponent to agree or whatever
	err = connection.ReadMsg()
	if err != nil {
		fmt.Println(err)
		if connection.QueueLen() == 0 {
			return err
		}
	}
	msg, err := connection.PopMsg()
	fmt.Println(msg)

	return nil
}

func main() {

	// 1. connect to server
	connection := sloopNet.GameConn{}
	err := connection.Connect("localhost", "8080")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Quit()

	// 2. wait for server to send a "start game" message
	err = connection.ReadMsg()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 3. read the message
	msg, err := connection.PopMsg()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 4. do stuff based on what that message is
	if msg[0] != "g" {
		fmt.Println("Nope, we not doing that")
		return
	} else {
		if msg[1] == "begin" {
			//acknowledge the _g_begin
			err = connection.SendMsg("_g_begin:")
			if err != nil {
				fmt.Println(err)
				connection.SendMsg("_c_end:")
				return
			}

			// start the game!
			fmt.Println("starting the game")
			err = startGame(&connection)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("Game over")
		}
	}
}
