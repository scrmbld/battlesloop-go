package sloopGame

import (
	"testing"
)

func TestShipCovers(t *testing.T) {
	// ship starting at 3, 3 of length 4 in horizontal orientation
	s1 := newShip(3, 3, 4, false)

	// test hits
	hit := s1.Covers(3, 3)
	if !hit {
		t.Fatalf("s1.Covers(3, 3) returned %v, should be true", hit)
	}
	hit = s1.Covers(3, 6)
	if !hit {
		t.Fatalf("s1.Covers(3, 6) returned %v, should be true", hit)
	}

	// test misses
	hit = s1.Covers(4, 3)
	if hit {
		t.Fatalf("s1.Covers(4, 3) returned %v, should be false", hit)
	}
	hit = s1.Covers(2, 3)
	if hit {
		t.Fatalf("s1.Covers(2, 3) returned %v, should be false", hit)
	}
	hit = s1.Covers(4, 6)
	if hit {
		t.Fatalf("s1.Covers(4, 6) returned %v, should be false", hit)
	}
	hit = s1.Covers(2, 6)
	if hit {
		t.Fatalf("s1.Covers(2, 6) returned %v, should be false", hit)
	}
	hit = s1.Covers(3, 9)
	if hit {
		t.Fatalf("s1.Covers(3, 9) returned %v, should be false", hit)
	}
	hit = s1.Covers(3, 0)
	if hit {
		t.Fatalf("s1.Covers(3, 0) returned %v, should be false", hit)
	}
	hit = s1.Covers(3, 7)
	if hit {
		t.Fatalf("s1.Covers(3, 7) returned %v, should be false", hit)
	}
	hit = s1.Covers(3, 2)
	if hit {
		t.Fatalf("s1.Covers(3, 2) returned %v, should be false", hit)
	}

	// ship starting at 3, 3 of length 5 in vertical orientation
	s2 := newShip(3, 3, 5, true)

	// test hits
	hit = s2.Covers(3, 3)
	if !hit {
		t.Fatalf("s2.Covers(3, 3) returned %v, should be true", hit)
	}
	hit = s2.Covers(7, 3)
	if !hit {
		t.Fatalf("s2.Covers(7, 3) returned %v, should be true", hit)
	}

	// test misses
	hit = s2.Covers(3, 4)
	if hit {
		t.Fatalf("s2.Covers(3, 4) returned %v, should be false", hit)
	}
	hit = s2.Covers(3, 2)
	if hit {
		t.Fatalf("s2.Covers(3, 2) returned %v, should be false", hit)
	}
	hit = s2.Covers(6, 4)
	if hit {
		t.Fatalf("s2.Covers(6, 4) returned %v, should be false", hit)
	}
	hit = s2.Covers(6, 2)
	if hit {
		t.Fatalf("s2.Covers(6, 2) returned %v, should be false", hit)
	}
	hit = s2.Covers(8, 3)
	if hit {
		t.Fatalf("s2.Covers(8, 3) returned %v, should be false", hit)
	}
	hit = s2.Covers(2, 3)
	if hit {
		t.Fatalf("s2.Covers(2, 3) returned %v, should be false", hit)
	}
}

func TestFire(t *testing.T) {
	var sea [10][10]uint8
	/* Things the fire function should be able to handle:
	1. Correct behavior on valid input
	2. Return an error if the board position is invalid
		We don't want nil values to get everwhere in our program
	It might be good to have it do bounds checking on the coordinates
	| but that runtime error shouldn't be too hard to debug
	*/
	sea[0][0] = 1  // miss
	sea[9][5] = 2  // ship
	sea[4][7] = 3  // hit
	sea[1][1] = 42 // invalid value

	err := fire(&sea, 0, 5) // fire at 0
	if sea[0][5] != 1 || err != nil {
		t.Fatalf("Fire(&sea, 0, 5) | sea[0][5] == 0 returned %v, %v, should be 1, nil", sea[0][5], err)
	}

	err = fire(&sea, 0, 0)
	if sea[0][0] != 1 || err != nil {
		t.Fatalf("Fire(&sea, 0, 0) | sea[0][0] == 1 returned %v, %v, should be 1, nil", sea[0][0], err)

	}

	err = fire(&sea, 9, 5)
	if sea[9][5] != 3 || err != nil {
		t.Fatalf("Fire(&sea, 9, 5) | sea[9][5] == 2 returned %v, %v, should be 3, nil", sea[9][5], err)

	}

	err = fire(&sea, 4, 7)
	if sea[4][7] != 3 || err != nil {
		t.Fatalf("Fire(&sea, 4, 7) | sea[4][7] == 3 returned %v, %v, should be 3, nil", sea[4][7], err)

	}

	err = fire(&sea, 1, 1)
	if err == nil {
		t.Fatalf("Fire(&sea, 1, 1) | sea[1][1] == 42 returned %v, %v, should be *, error", sea[1][1], err)
	}
}

func TestPlaceShip(t *testing.T) {
	var board Board
	var empty_sea [10][10]uint8

	// test bounds checking
	// in bounds
	correct_sea := empty_sea
	correct_sea[0][0] = 2
	err := board.PlaceShip(0, 0, 1, false)
	if err != nil {
		t.Fatalf("board.PlaceShip(0, 0, 1, false) | returned %v, should be nil", err)
	} else if board.OurSea != correct_sea {
		t.Fatalf("board.PlaceShip(0, 0, 1, false) | gave board.OurSea of %v, should be %v", board.OurSea, correct_sea)
	}
	board.OurSea = empty_sea

	// in bounds
	correct_sea = empty_sea
	correct_sea[9][9] = 2
	err = board.PlaceShip(9, 9, 1, false)
	if err != nil {
		t.Fatalf("board.PlaceShip(9, 9, 1, false) | returned %v, should be nil", err)
	} else if board.OurSea != correct_sea {
		t.Fatalf("board.PlaceShip(9, 9, 1, false) | gave board.OurSea of %v, should be %v", board.OurSea, correct_sea)
	}
	board.OurSea = empty_sea

	// in bounds
	correct_sea = empty_sea
	correct_sea[9][8] = 2
	correct_sea[9][9] = 2
	err = board.PlaceShip(9, 8, 2, false)
	if err != nil {
		t.Fatalf("board.PlaceShip(9, 8, 2, false) | returned %v, should be nil", err)
	} else if board.OurSea != correct_sea {
		t.Fatalf("board.PlaceShip(9, 8, 2, false) | gave board.OurSea of %v, should be %v", board.OurSea, correct_sea)
	}
	board.OurSea = empty_sea

	// in bounds
	correct_sea = empty_sea
	correct_sea[9][9] = 2
	correct_sea[8][9] = 2
	err = board.PlaceShip(8, 9, 2, true)
	if err != nil {
		t.Fatalf("board.PlaceShip(8, 9, 2, true) | returned %v, should be nil", err)
	} else if board.OurSea != correct_sea {
		t.Fatalf("board.PlaceShip(8, 9, 2, true) | gave board.OurSea of %v, should be %v", board.OurSea, correct_sea)
	}
	board.OurSea = empty_sea

	// out of bounds -- y too small
	err = board.PlaceShip(-1, 1, 1, false)
	if err == nil {
		t.Fatalf("board.PlaceShip(-1, 1, 1, false) | returned %v, should be error", err)
	} else if board.OurSea != empty_sea {
		t.Fatalf("board.PlaceShip(-1, 1, 1, false) | gave board.OurSea of %v, should be %v", board.OurSea, empty_sea)
	}

	// out of bounds -- x too small
	err = board.PlaceShip(1, -1, 1, false)
	if err == nil {
		t.Fatalf("board.PlaceShip(1, -1, 1, false) | returned %v, should be error", err)
	} else if board.OurSea != empty_sea {
		t.Fatalf("board.PlaceShip(1, -1, 1, false) | gave board.OurSea of %v, should be %v", board.OurSea, empty_sea)
	}

	// out of bounds -- y too big
	err = board.PlaceShip(10, 1, 1, false)
	if err == nil {
		t.Fatalf("board.PlaceShip(10, 1, 1, false) | returned %v, should be error", err)
	} else if board.OurSea != empty_sea {
		t.Fatalf("board.PlaceShip(10, 1, 1, false) | gave board.OurSea of %v, should be %v", board.OurSea, empty_sea)
	}

	// out of bounds -- x too big
	err = board.PlaceShip(1, 10, 1, false)
	if err == nil {
		t.Fatalf("board.PlaceShip(1, 10, 1, false) | returned %v, should be error", err)
	} else if board.OurSea != empty_sea {
		t.Fatalf("board.PlaceShip(1, 10, 1, false) | gave board.OurSea of %v, should be %v", board.OurSea, empty_sea)
	}

	// out of bounds -- off the bottom
	err = board.PlaceShip(7, 0, 4, true)
	if err == nil {
		t.Fatalf("board.PlaceShip(7, 0, 4, false) | returned %v, should be error", err)
	} else if board.OurSea != empty_sea {
		t.Fatalf("board.PlaceShip(7, 0, 4, false) | gave board.OurSea of %v, should be %v", board.OurSea, empty_sea)
	}

	// out of bounds -- off the right
	err = board.PlaceShip(0, 7, 4, false)
	if err == nil {
		t.Fatalf("board.PlaceShip(0, 7, 4, false) | returned %v, should be error", err)
	} else if board.OurSea != empty_sea {
		t.Fatalf("board.PlaceShip(0, 7, 4, false) | gave board.OurSea of %v, should be %v", board.OurSea, empty_sea)
	}

	// test intersections
	// illegal, as there is an intersection
	err = board.PlaceShip(2, 2, 3, true) // starting at 2, 2 going down
	if err != nil {
		t.Fatalf("board.PlaceShip(2, 2, 3, true) | returned %v, should be nil", err)
	}
	err = board.PlaceShip(2, 1, 3, false) // starting at 2, 1 going right (should intersect)
	if err == nil {
		t.Fatalf("board.PlaceShip(2, 1, 3, false) | returned %v, should be error (intersection)", err)
	}
	board.OurSea = empty_sea

	// legal, as there is no intersection
	err = board.PlaceShip(2, 2, 3, true) // starting at 2, 2 going down
	if err != nil {
		t.Fatalf("board.PlaceShip(2, 2, 3, true) | returned %v, should be nil", err)
	}
	err = board.PlaceShip(2, 1, 3, true) // starting at 2, 1 going down
	if err != nil {
		t.Fatalf("board.PlaceShip(2, 2, 3, true) | returned %v, should be nil", err)
	}
}

// not sure if I want to unit test this component, since it defines UI appearance
// So im just gonna put it off for later
// func TestRenderSea(t* testing.T) {
// 	ENEMY SEA
// === A  B  C  D  E  F  G  H  I  J =
// 0 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 1 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 2 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 3 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 4 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 5 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 6 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 7 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 8 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 9 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// ==================================
// OUR SEA
// === A  B  C  D  E  F  G  H  I  J =
// 0 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 1 - ~  ~  M  ~  ~  ~  ~  ~  ~  ~ -
// 2 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 3 - ~  ~  ~  ~  ~  ~  ~  ~  H  ~ -
// 4 - ~  S  ~  ~  ~  ~  ~  ~  ~  ~ -
// 5 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 6 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 7 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 8 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 9 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// ==================================
// ENEMY SEA
// === A  B  C  D  E  F  G  H  I  J =
// 0 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 1 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 2 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 3 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 4 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 5 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 6 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 7 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 8 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 9 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// ==================================
// OUR SEA
// === A  B  C  D  E  F  G  H  I  J =
// 0 - M  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 1 - ~  ~  M  ~  ~  ~  ~  ~  ~  ~ -
// 2 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 3 - ~  ~  ~  ~  ~  ~  ~  ~  H  ~ -
// 4 - ~  H  ~  ~  ~  ~  ~  ~  ~  ~ -
// 5 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 6 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 7 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 8 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// 9 - ~  ~  ~  ~  ~  ~  ~  ~  ~  ~ -
// ==================================
//
// }
