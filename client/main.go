package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/scrmbld/battlesloop-go/sloopGame"
	"github.com/scrmbld/battlesloop-go/sloopNet"
	"os"
	"strconv"
	"strings"
)

func askShipLoc(board *sloopGame.Board, length uint8) {
	reader := bufio.NewReader(os.Stdin)
	for { // try over and over until nothing goes wrong (aka user enters valid input)

		// ask user for orientation, x, and y
		fmt.Printf("Please determine placement for ship of length %v\n", length)

		// orienation
		fmt.Print("Pick an orientation, horizontal (0) or vertical (1): ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print("Invalid input: ")
			fmt.Println(err)
			continue
		}
		text = strings.Trim(text, "\n")

		num, err := strconv.Atoi(text)
		if err != nil {
			fmt.Print("Invalid input: ")
			fmt.Println(err)
			continue
		}
		if num > 1 || num < 0 {
			fmt.Println("Invalid input: number out of range")
		}

		orientation := num != 0
		fmt.Print("\n")

		// x
		fmt.Println("Enter a horizontal coordinate (A-J): ")
		text, err = reader.ReadString('\n')
		if err != nil {
			fmt.Print("Invalid input: ")
			fmt.Println(err)
			continue
		}
		text = strings.Trim(text, "\n")

		var colMap = map[string]uint8{
			"A": 0,
			"B": 1,
			"C": 2,
			"D": 3,
			"E": 4,
			"F": 5,
			"G": 6,
			"H": 7,
			"I": 8,
			"J": 9,
		}
		x, valid := colMap[text]
		if !valid {
			fmt.Println("Invalid input: Column not in range")
			continue
		}
		fmt.Print("\n")

		// y
		fmt.Println("Enter a vertical coordinate (0-9): ")
		text, err = reader.ReadString('\n')
		if err != nil {
			fmt.Print("Invalid input: ")
			fmt.Println(err)
			continue
		}
		text = strings.Trim(text, "\n")

		y, err := strconv.Atoi(text)
		if err != nil {
			fmt.Print("Invalid input: ")
			fmt.Println(err)
			continue
		}
		if y < 0 || y > 9 {
			fmt.Println("Invalid input: row not in range")
			continue
		}
		fmt.Print("\n")

		//check the validity of the entered position
		err = board.PlaceShip(int(y), int(x), int(length), orientation)
		if err != nil {
			fmt.Print("Invalid input: ")
			fmt.Println(err)
			continue
		}

		break
	}

}

func playGame(connection *sloopNet.GameConn) error {
	var board sloopGame.Board

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
	if err != nil {
		fmt.Println(err)
		return err
	}
	if msg[0] != "g" && msg[1] != "first" {
		fmt.Println("Unexpected response")
		return errors.New("Unexpected response")
	}

	fmt.Println("Connection Successful! Beginning game...")
	board.PrintBoard()

	// begin placing ships
	ships := []uint8{5, 4, 3, 3, 2}
	for _, v := range ships {
		askShipLoc(&board, v)

		// draw the new state of the board
		board.PrintBoard()
	}

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
			err = playGame(&connection)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("Game over")
		}
	}
}
