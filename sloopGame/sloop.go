package main

import (
	"errors"
	"fmt"
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

func printSea(sea [10][10]uint8, cellMap map[uint8]rune) {

	fmt.Printf("--- A  B  C  D  E  F  G  H  I  J -\n")
	for i := range sea {
		fmt.Printf("%v -", i)
		for j := range sea {
			var cellRune = cellMap[sea[i][j]]
			fmt.Printf(" %c ", cellRune)
		}
		fmt.Printf("-\n")
	}
	fmt.Printf("--------------------------------\n")
}

func fire(sea *[10][10]uint8, y uint8, x uint8) error {
	var effectMap = map[uint8]uint8{0: 1, 1: 1, 2: 3, 3: 3}
	var newVal, valid = effectMap[sea[y][x]]
	if !valid {
		return errors.New("Invalid cell value for sea")
	}

	sea[y][x] = newVal

	return nil
}

func (b Board) printBoard() {
	// map a cell value (how they are stored in the ourSea and enemySea uint8 arrays)
	// to its visual representation
	var cellMap = map[uint8]rune{0: '~', 1: 'M', 2: 'S', 3: 'H'}

	fmt.Printf("ENEMY SEA\n")
	printSea(b.enemySea, cellMap)

	fmt.Printf("OUR SEA\n")
	printSea(b.ourSea, cellMap)
}

func main() {
	var b Board

	b.ourSea[1][2] = 1 // miss at C1
	b.ourSea[4][1] = 2 // hit at B4
	b.ourSea[3][8] = 3 // ship at I3
	b.printBoard()

	_ = fire(&b.ourSea, 1, 2)
	_ = fire(&b.ourSea, 4, 1)
	_ = fire(&b.ourSea, 3, 8)
	_ = fire(&b.ourSea, 0, 0)
	b.printBoard()
}
