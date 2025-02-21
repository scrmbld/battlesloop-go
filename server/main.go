package main

import (
	"errors"
	"fmt"
	"github.com/scrmbld/battlesloop-go/sloopGame"
	"github.com/scrmbld/battlesloop-go/sloopNet"
	"math/rand/v2"
	"net"
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

		return y, x
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

	if msgs[0] == "h" { // we hit
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

		hit, err := board.FireFriendly(pos[0], pos[1])
		if err != nil {
			return 3, err
		}

		if hit {
			fmt.Printf("Opponent hit %v\n", pos)
			// tell the opponent that they hit
			connection.SendMsg("_h_" + msg[1] + ":")
		} else {
			fmt.Printf("Opponent missed %v\n", pos)
			// tell the opponent that they missed
			connection.SendMsg("_m_" + msg[1] + ":")
		}

		// check win/loss
		if board.CheckLoss() {
			connection.SendMsg("_g_loss:")
			return 1, nil
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
			result, err := ourTurn(connection, &board)
			if err != nil {
				return err
			}
			if result != 0 {
				break
			}
		} else { // their turn
			board.PrintBoard()
			result, err := oppTurn(connection, &board)
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
			} else {
				fmt.Println("game over")
			}
		}
	}

}
