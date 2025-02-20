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
	OurSea   [10][10]uint8
	EnemySea [10][10]uint8
	// uints represent "health" of ships (aka their length)
	OurFleet []uint8
	// do not track enemy fleet -- opponent will tell us when we win

	// True: our turn, False: enemy turn
	WhoseTurn bool
}

func RenderSea(sea [10][10]uint8, cellMap map[uint8]rune) (string, error) {

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

func (b Board) PrintBoard() {
	// map a cell value (how they are stored in the ourSea and enemySea uint8 arrays)
	// to its visual representation
	var cellMap = map[uint8]rune{0: '~', 1: 'M', 2: 'S', 3: 'H'}

	seaStr, err := RenderSea(b.EnemySea, cellMap)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("ENEMY SEA\n%v", seaStr)

	seaStr, err = RenderSea(b.OurSea, cellMap)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("OUR SEA\n%v", seaStr)
}

// TODO: Board setup
func (b *Board) PlaceShip(origin_y int, origin_x int, size int, orientation bool) error {
	if size < 1 {
		// errors will be printed to console and can be caused by invalid user input
		// if this one happens, it will be because of a bug and not invalid user input
		return errors.New("Invalid ship size (this is a bug)")
	}
	end_y := origin_y
	end_x := origin_x
	if orientation {
		end_y = origin_y + size
	} else {
		end_x = origin_x + size
	}

	// bounds checking
	if origin_y < 0 || origin_y >= 10 || origin_x < 0 || origin_x > 10 ||
		end_y < 0 || end_y > 10 || end_x < 0 || end_x > 10 {
		return errors.New("Ship placement goes out of bounds")
	}

	// collision checking (ships cannot intersect)
	for i := 0; i < size; i++ {
		y := origin_y
		x := origin_x
		if orientation {
			y = y + i
		} else {
			x = x + i
		}

		if b.OurSea[y][x] == 2 {
			return errors.New("Ship placement intersects with existing ship")
		}
	}

	// actually adding the ship to the board
	for i := 0; i < size; i++ {
		y := origin_y
		x := origin_x
		if orientation {
			y = y + i
		} else {
			x = x + i
		}

		b.OurSea[y][x] = 2
	}

	return nil
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
