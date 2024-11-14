package sloopGame

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Board struct {
	// uints represent state of each cell
	// 0: empty, 1: miss, 2: ship (not hit), 3: hit
	ourSea   [10][10]uint8
	enemySea [10][10]uint8
	// uints represent "health" of ships (aka their length)
	ourFleet []uint8
	// do not track enemy fleet -- opponent will tell us when we win

	// True: our turn, False: enemy turn
	whoseTurn bool
}

func renderSea(sea [10][10]uint8, cellMap map[uint8]rune) (string, error) {

	var seaText strings.Builder
	seaText.WriteString("=== A  B  C  D  E  F  G  H  I  J =\n")
	for i := range sea {
		var row string = strconv.Itoa(i)
		seaText.WriteString(row)
		seaText.WriteString(" -")
		for j := range sea {
			var cellRune, valid = cellMap[sea[i][j]]
			if !valid {
				return "", errors.New("Invalid cell value for sea")
			}

			seaText.WriteString(" ")
			seaText.WriteRune(cellRune)
			seaText.WriteString(" ")
		}
		seaText.WriteString("-\n")
	}
	seaText.WriteString("==================================\n")

	return seaText.String(), nil
}

func (b Board) printBoard() {
	// map a cell value (how they are stored in the ourSea and enemySea uint8 arrays)
	// to its visual representation
	var cellMap = map[uint8]rune{0: '~', 1: 'M', 2: 'S', 3: 'H'}

	seaStr, err := renderSea(b.enemySea, cellMap)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("ENEMY SEA\n%v", seaStr)

	seaStr, err = renderSea(b.ourSea, cellMap)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("OUR SEA\n%v", seaStr)
}

func Fire(sea *[10][10]uint8, y uint8, x uint8) error {
	var effectMap = map[uint8]uint8{0: 1, 1: 1, 2: 3, 3: 3}
	var newVal, valid = effectMap[sea[y][x]]
	if !valid {
		return errors.New("Invalid cell value for sea")
	}

	sea[y][x] = newVal

	return nil
}
