package sloopNet

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

type GameConn struct {
	sock net.Conn
	// treat this like a queue
	msgs [][]string
}

// / Initialize a GameConn with a socket
func (self *GameConn) Init(sock net.Conn) error {
	self.sock = sock
	return nil
}

func (self *GameConn) Connect(addr string, port string) error {
	uri := addr + string(':') + port
	var err error
	self.sock, err = net.Dial("tcp", uri)
	if err != nil {
		return err
	}

	return nil
}

// close the connection and deallocate msgs queue
func (self *GameConn) Quit() {
	self.sock.Close()
	self.msgs = nil
}

// read a ':' delimited message
func (self *GameConn) ReadMsg() error {

	local_buf := make([]byte, 256)
	// read until we end at a ':'
	for {
		reading := make([]byte, 64)
		idx, err := self.sock.Read(reading)
		if err != nil {
			return err
		}

		local_buf = append(local_buf, reading[:idx]...)

		if reading[idx-1] == byte(':') {
			break
		}
	}

	new_msgs, err := ParseMsgs(string(local_buf[:]))
	if err != nil {
		return err
	}

	self.msgs = append(self.msgs, new_msgs...)

	return nil
}

func (self *GameConn) PopMsg() ([]string, error) {
	// peek
	front := self.msgs[0]
	if front == nil {
		return nil, errors.New("No messages in msg queue")
	}

	// pop
	// this is actually an accepted way to do this in go (even though I don't like it)
	// https://stackoverflow.com/questions/2818852/is-there-a-queue-implementation
	self.msgs = self.msgs[1:]

	return front, nil
}

// returns the length of self.msgs
func (self *GameConn) QueueLen() int {
	return len(self.msgs)
}

func (self *GameConn) DumpQueue() string {
	var strBuilder strings.Builder

	for i := range self.msgs {
		var lineBuilder strings.Builder
		for j := range self.msgs[i] {
			lineBuilder.WriteString(self.msgs[i][j])
			lineBuilder.WriteString(", ")
		}
		strBuilder.WriteString(lineBuilder.String())
		strBuilder.WriteString("\n")
	}

	return strBuilder.String()
}

// send a message of the caller's choosing over the socket
func (self *GameConn) SendMsg(msg string) error {
	_, err := self.sock.Write([]byte(msg))
	return err
}

func ParsePos(pos string) ([2]uint8, error) {

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

	// split the string
	var posStrs []string = strings.Split(pos, "-")
	var result [2]uint8

	// row
	// This MUST be an integer
	row, err := strconv.Atoi(posStrs[1])
	if err != nil {
		return [2]uint8{}, err
	}
	// do bounds checking
	if row < 0 || row > 9 {
		return [2]uint8{}, errors.New("Row number out of range")
	}
	result[0] = uint8(row)

	// column
	col, valid := colMap[posStrs[0]]
	if !valid {
		return [2]uint8{}, errors.New("Invalid column identifier")
	}
	result[1] = col

	return result, nil

}

// Given two ints, return the string representation of the position they equate to.
func PosFromInts(y int, x int) (string, error) {
	var colMap = map[int]string{
		0: "A",
		1: "B",
		2: "C",
		3: "D",
		4: "E",
		5: "F",
		6: "G",
		7: "H",
		8: "I",
		9: "J",
	}

	col, valid := colMap[x]
	if !valid {
		return "", errors.New("Column number out of range")
	}

	if y < 0 || y > 9 {
		return "", errors.New("Row number out of range")
	}

	posString := col + "-" + strconv.Itoa(y)
	return posString, nil
}

// Parse a string that contains battlesloop messages
// Turns each message into a slice of strings, split by "_"'s
// There's no way to be sure that a socket will only return one message,
// so we have to be able to handle multiple
//
// The first index of each message should be "" because types have '_'s both before and after
func ParseMsgs(messages string) ([][]string, error) {
	var msgs []string = strings.Split(messages, ":")
	var result [][]string
	for i := range len(msgs) - 1 {
		var ss []string = strings.Split(msgs[i], "_")
		if len(ss) != 3 {
			return [][]string{}, errors.New("Invalid message: has no type or too many types")
		}
		result = append(result, ss[1:]) // first value will be "" because of the leading "_"
	}

	// if there isn't an empty slice at the end of msgs, then we were missing a delimiter
	// might not want to lose all the other messages though
	if len(msgs[len(msgs)-1]) != 0 {
		return result, errors.New("Invalid message: Last message not delimited")
	}

	return result, nil
}
