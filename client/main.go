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

// asks the user for a board position
// returns y, x
func getUserPos() (uint8, uint8) {
	reader := bufio.NewReader(os.Stdin)
	// try over and over again until the user gives us valid input
	for {
		fmt.Println("Enter coordinates (e.g., A0, B6, J9): ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Print("Invalid input: ")
			fmt.Println(err)
			continue
		}
		text = strings.Trim(text, "\n")

		if len(text) != 2 {
			fmt.Println("Invalid input: please enter inputs in the format {Row Letter}{Column Letter} e.g. A0")
			continue
		}

		posString := text[:1] + "-" + text[1:]
		pos, err := sloopNet.ParsePos(posString)
		if err != nil {
			fmt.Println(err)
			continue
		}

		return pos[0], pos[1]
	}

}

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

		// x and y
		y, x := getUserPos()

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

// Play our turn. Returns 0 if game is not over, 1 on loss (impossible), 2 on win, 3 on other ending.
func ourTurn(connection *sloopNet.GameConn, board *sloopGame.Board) (int, error) {

	// ask the user where they want to shoot
	fire_y, fire_x := getUserPos()
	posString, err := sloopNet.PosFromInts(int(fire_y), int(fire_x))
	if err != nil {
		return 3, err
	}

	// send the position to the opponent
	err = connection.SendMsg("_f_" + posString + ":")
	if err != nil {
		return 3, err
	}

	// now, wait for their response
	err = connection.ReadMsg()
	if err != nil {
		return 3, err
	}
	msgs, err := connection.PopMsg()
	if err != nil {
		return 3, err
	}

	if msgs[0] == "a" { // we hit & sunk
		displayPos := (posString[:1] + posString[2:])
		fmt.Printf("We hit on %v, sinking their ship!\n", displayPos)

		pos, _ := sloopNet.ParsePos(posString)
		board.EnemySea[pos[0]][pos[1]] = 3
		return 0, nil
	} else if msgs[0] == "h" { // we hit
		displayPos := (posString[:1] + posString[2:])
		fmt.Printf("We hit on %v\n", displayPos)

		pos, _ := sloopNet.ParsePos(posString)
		board.EnemySea[pos[0]][pos[1]] = 3
		return 0, nil

	} else if msgs[0] == "m" { // we missed
		displayPos := (posString[:1] + posString[2:])
		fmt.Printf("We missed on %v\n", displayPos)

		pos, _ := sloopNet.ParsePos(posString)
		board.EnemySea[pos[0]][pos[1]] = 1
		return 0, nil
	} else if msgs[0] == "g" && msgs[1] == "lose" { // we won
		// assume this means that we hit
		pos, _ := sloopNet.ParsePos(posString)
		board.EnemySea[pos[0]][pos[1]] = 3
		connection.SendMsg("_g_win:")
		return 2, nil
	} else if msgs[0] == "g" && msgs[1] == "end" { // no contest
		return 3, nil
	} else {
		return 3, errors.New("Opponent sent unexpected message: " + "_" + msgs[0] + "_" + msgs[1] + ":")
	}

}

// Opponent's turn. Returns 0 if game is not over, 1 on loss, 2 on win (impossible), 3 on other ending.
func oppTurn(connection *sloopNet.GameConn, board *sloopGame.Board) (int, error) {

	// wait for the opponent's decision
	err := connection.ReadMsg()
	if err != nil {
		return 3, err
	}
	msg, err := connection.PopMsg()

	if msg[0] == "g" && msg[1] == "end" {
		return 3, nil
	} else if msg[0] == "f" {
		pos, err := sloopNet.ParsePos(msg[1])
		if err != nil {
			return 3, err
		}

		hit, sunk, err := board.FireFriendly(pos[0], pos[1])
		if err != nil {
			return 3, err
		}

		// check win/loss
		if board.CheckLoss() {
			connection.SendMsg("_g_lose:")
			return 1, nil
		}

		// tell the opponent if they hit, missed, etc.
		if hit && sunk {
			fmt.Printf("Opponent hit %v%v, sinking our ship!\n", msg[1][:1], msg[1][2:])
			// tell the opponent that they sunk our ship
			connection.SendMsg("_a_" + msg[1] + ":")
		} else if hit {
			fmt.Printf("Opponent hit %v%v\n", msg[1][:1], msg[1][2:])
			// tell the opponent that they hit our ship
			connection.SendMsg("_h_" + msg[1] + ":")
		} else {
			fmt.Printf("Opponent missed on %v%v\n", msg[1][:1], msg[1][2:])
			// tell the opponent that they missed
			connection.SendMsg("_m_" + msg[1] + ":")
		}

	} else {
		return 3, errors.New("Opponent sent unexpected message: " + "_" + msg[0] + "_" + msg[1] + ":")
	}

	return 0, nil
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
	board.WhoseTurn = false

	fmt.Println("Connection Successful! Beginning game...")
	board.PrintBoard()

	// begin placing ships
	ships := []uint8{5, 4, 3, 3, 2}
	for _, v := range ships {
		askShipLoc(&board, v)

		// draw the new state of the board
		board.PrintBoard()
	}

	// now we are done with the setup phase
	// tell opponent that we are ready to continue
	connection.SendMsg("_g_ready:")

	// wait for opponent to tell us that they are ready
	err = connection.ReadMsg()
	if err != nil {
		return err
	}
	msg, err = connection.PopMsg()
	if err != nil {
		return err
	}
	if msg[0] != "g" && msg[1] != "ready" {
		return errors.New("Opponent did not ready up")
	}

	// once opponent is ready, begin the main loop
	var result int
	for {

		if board.WhoseTurn { // our turn
			board.PrintBoard()
			result, err = ourTurn(connection, &board)
			if err != nil {
				return err
			}
			if result != 0 {
				break
			}
		} else { // their turn
			board.PrintBoard()
			result, err = oppTurn(connection, &board)
			if err != nil {
				return err
			}

			if result != 0 {
				break
			}
		}

		board.WhoseTurn = !board.WhoseTurn
	}

	fmt.Printf("Result: %v\n", result)

	if result == 1 {
		fmt.Println("You Lost...")
	} else if result == 2 {
		fmt.Println("You Won!")
	} else if result == 3 {
		fmt.Println("Game ended: opponent decided to end the game")
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
