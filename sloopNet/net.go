package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type PositionalMsg struct {
	msgContent string
	row        uint8
	col        uint8
}

func (m PositionalMsg) GetType() string {
	return m.msgContent
}

type ShipMsg struct {
	msgContent string
	shipId     uint8
}

func (m ShipMsg) GetType() string {
	return m.msgContent
}

type GameMsg struct {
	msgContent string
}

func (m GameMsg) GetType() string {
	return m.msgContent
}

type ConnMsg struct {
	msgContent string
}

func (m ConnMsg) GetType() string {
	return m.msgContent
}

type SloopMsg interface {
	GetType() string
}

// see the documentation for positional message type in Battlesloop Protocol
func parsePositional(message string) (PositionalMsg, error) {
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

	// the return value which we will build up
	var result PositionalMsg

	// separate type from content
	var parts []string = strings.Split(message, "_")
	result.msgContent = parts[1]

	// parse the content
	var posStrs []string = strings.Split(parts[2], "-")

	// row
	// This MUST be an integer
	row, err := strconv.Atoi(posStrs[1])
	if err != nil {
		return PositionalMsg{}, err
	}
	// do bounds checking
	if row < 0 || row > 9 {
		return PositionalMsg{}, errors.New("Row number out of range")
	}
	result.row = uint8(row)

	// column
	col, valid := colMap[posStrs[0]]
	if !valid {
		return PositionalMsg{}, errors.New("Invalid column identifier")
	}
	result.col = col

	return result, nil
}

func parseShip(message string) (ShipMsg, error) {

	var result ShipMsg

	// separate type from content
	var parts []string = strings.Split(message, "_")
	result.msgContent = parts[1]

	// parse the ship ID
	val, err := strconv.Atoi(parts[2])
	if err != nil {
		return ShipMsg{}, err
	}
	result.shipId = uint8(val)

	return result, nil
}

func parseGame(message string) (GameMsg, error) {
	var result GameMsg

	// separate type from content
	var parts []string = strings.Split(message, "_")
	result.msgContent = parts[2]

	return result, nil
}

func parseConn(message string) (ConnMsg, error) {
	var result ConnMsg

	// separate type from content
	var parts []string = strings.Split(message, "_")
	result.msgContent = parts[2]

	return result, nil
}

// func ParseMessage(msg string) SloopMsg {
// 	var typeMap = map[rune]interface{} {
// 		'h': parsePositional,
// 		'm': parsePositional,
// 		'f': parsePositional,
// 		'a': parseShip,
// 		'g': parseGame,
// 		'c': parseConn,
// 	}
// 	// split msg based on ':' delimiters
// 	var messages []string = strings.Split(msg, ":")
//
// 	// figure out what each part is
// 	for _,m := range(messages) {
// 		// separate the type from the contents
// 		var parts = strings.Split(m, "_")
//
// 	}
//
// }

func main() {
	message := "_h_A-7"
	result, err := parsePositional(message)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}
