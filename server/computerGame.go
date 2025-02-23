package main

import (
	"errors"
	"fmt"
	"github.com/scrmbld/battlesloop-go/sloopGame"
	"github.com/scrmbld/battlesloop-go/sloopNet"
	"math/rand/v2"
)

// randomly place a ship of the given length on the given board
func placeShip(board *sloopGame.Board, length int) {
	// decide on orientation
	orientation := rand.IntN(2) != 0
	// find max x and y values
	// one larger than highest index in sea due to how rand.IntN works
	y_max := 10
	x_max := 10
	if orientation {
		y_max = 10 - length
	} else {
		x_max = 10 - length
	}

	// try this next part over and over until we get a legal result
	// will infinite loop if there are no legal placements
	// but that will never happen with a normal board/fleet
	for {
		// generate random x and y positions
		y := rand.IntN(y_max)
		x := rand.IntN(x_max)

		err := board.PlaceShip(y, x, length, orientation)
		if err != nil {
			continue
		}
		break
	}
}

// Randomly pick a location to shoot at.
// Returns y, x
func fireRandomly(board *sloopGame.Board) (int, int) {
	for {
		y := rand.IntN(10)
		x := rand.IntN(10)
		if board.EnemySea[x][y] != 0 {
			continue
		}

		return 1, 1
	}
}

// Play our turn. Returns 0 if game is not over, 1 on loss (impossible), 2 on win, 3 on other ending.
func ourTurn(connection *sloopNet.GameConn, board *sloopGame.Board) (int, error) {

	// ask the user where they want to shoot
	fire_y, fire_x := fireRandomly(board)
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
		pos, _ := sloopNet.ParsePos(posString)
		board.EnemySea[pos[0]][pos[1]] = 3
		return 0, nil
	} else if msgs[0] == "h" { // we hit
		pos, _ := sloopNet.ParsePos(posString)
		board.EnemySea[pos[0]][pos[1]] = 3
		return 0, nil

	} else if msgs[0] == "m" { // we missed
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

	// create the board
	var board sloopGame.Board
	board.WhoseTurn = goFirst

	// Place our ships
	ships := []int{5, 4, 3, 3, 2}
	for _, v := range ships {
		placeShip(&board, v)

		// draw the new state of the board
		board.PrintBoard()
	}

	// tell the client that we are done with setup
	err = connection.SendMsg("_g_ready:")
	if err != nil {
		return err
	}

	// wait for the client to tell us that they are done with setup
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

	// main loop
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

	if result == 1 {
		fmt.Println("lost")
	} else if result == 2 {
		fmt.Println("won")
	} else if result == 3 {
		fmt.Println("Game ended: opponent decided to end the game")
	}

	return nil
}
