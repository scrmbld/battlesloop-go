package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/scrmbld/battlesloop-go/sloopNet"
)

var queueMtx = sync.RWMutex{}
var queue = PlayerQueue{}

// The matchmaking queue.
// Stores players identified by their socket connection.
type PlayerQueue struct {
	players []*sloopNet.GameConn
	length  int
}

// Push a player to the matchmaking queue
func (self *PlayerQueue) Push(p *sloopNet.GameConn) {
	queueMtx.Lock()
	self.players = append(self.players, p)
	self.length++
	queueMtx.Unlock()
}

// Pop a pair of players from the matchmaking queue
func (self *PlayerQueue) PopPair() (*sloopNet.GameConn, *sloopNet.GameConn, error) {
	queueMtx.Lock()
	defer queueMtx.Unlock()

	if self.length < 2 {
		return &sloopNet.GameConn{}, &sloopNet.GameConn{}, errors.New("Not enough players in queue")
	}
	ps := self.players[:2]
	self.players = self.players[2:]
	self.length -= 2
	return ps[0], ps[1], nil
}

func matchMaker() {
	for {
		p1, p2, err := queue.PopPair()
		if err != nil {
			time.Sleep(time.Millisecond * 10) // check every hundredth of a second
			continue
		}
		// send each player the address of their opponent
		// additionally, tell one to dial and the other to listen
		p1_addr := strings.Split(p1.RemoteAddr(), ":")[0]
		p2_addr := strings.Split(p2.RemoteAddr(), ":")[0]
		fmt.Println(p1_addr)
		fmt.Println(p2_addr)
		p1.SendMsg("_s_l-" + p2_addr + ":")
		p2.SendMsg("_s_d-" + p1_addr + ":")

		// wait for p1 to confirm that they are listening
		err = p1.ReadMsg()
		if err != nil {
			p2.SendMsg("_c_end:") // report the failure to p2
			continue
		}
		msg, err := p1.PopMsg()
		if err != nil {
			p2.SendMsg("_c_end:")
			continue
		}

		if msg[0] == "s" && msg[1] == "ready" { // p1 readied up
			fmt.Println("p1 is ready")
			err = p2.SendMsg("_s_ready:")
			if err != nil {
				p1.SendMsg("_c_end:") // report the failure to p1
				continue
			}
		} else { // p1 did not ready up
			p2.SendMsg("_c_end:")
			continue
		}

		p1.Quit()
		p2.Quit()
	}
}

func main() {

	// listen for new connections from clients
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	go matchMaker()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleJoin(conn)
	}
}

// Called whenever someone connects to the server
func handleJoin(sock net.Conn) {
	var connection sloopNet.GameConn
	connection.Init(sock)

	// wait for client to tell us what they want to do
	err := connection.ReadMsg()
	if err != nil {
		fmt.Println(err)
		connection.Quit()
		return
	}

	msg, err := connection.PopMsg()
	if err != nil {
		fmt.Println(err)
		connection.Quit()
		return
	}

	if msg[0] == "s" && msg[1] == "computer" {
		// begin match vs. computer
		err := connection.SendMsg("_g_begin:")
		if err != nil {
			fmt.Println(err)
			connection.Quit()
			return
		}
		err = connection.ReadMsg()
		if err != nil {
			fmt.Println(err)
			connection.Quit()
			return
		}
		msg, err = connection.PopMsg()
		if msg[0] != "g" && msg[1] != "begin" {
			fmt.Println("Opponent did not confirm game start")
			connection.Quit()
			return
		}

		err = startGame(&connection)
		if err != nil {
			fmt.Println(err)
			connection.Quit()
		} else {
			fmt.Println("game over")
		}

		connection.Quit()

	} else if msg[0] == "s" && msg[1] == "player" {
		// add the current player to the queue
		queue.Push(&connection)
	} else {
		fmt.Println("Invalid choice from client")
		connection.Quit()
	}

}
