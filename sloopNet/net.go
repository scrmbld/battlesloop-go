package main

import (
	"errors"
	"strconv"
	"strings"
)

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

func ParseMsgs(message string) ([][]string, error) {
	var msgs []string = strings.Split(message, ":")
	var result [][]string
	for i := range msgs {
		var ss []string = strings.Split(msgs[i], "_")
		if len(ss) != 3 {
			return [][]string{}, errors.New("Invalid message: has no type or too many types")
		}
		result = append(result, ss[1:]) // first value will be "" because of the leading "_"
	}

	return result, nil
}
