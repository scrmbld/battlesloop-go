package sloopGame

import (
	"testing"
)

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

	err := Fire(&sea, 0, 5) // fire at 0
	if sea[0][5] != 1 || err != nil {
		t.Fatalf("Fire(&sea, 0, 5) | sea[0][5] == 0 returned %v, %v, should be 1, nil", sea[0][5], err)
	}

	err = Fire(&sea, 0, 0)
	if sea[0][0] != 1 || err != nil {
		t.Fatalf("Fire(&sea, 0, 0) | sea[0][0] == 1 returned %v, %v, should be 1, nil", sea[0][0], err)

	}

	err = Fire(&sea, 9, 5)
	if sea[9][5] != 3 || err != nil {
		t.Fatalf("Fire(&sea, 9, 5) | sea[9][5] == 2 returned %v, %v, should be 3, nil", sea[9][5], err)

	}

	err = Fire(&sea, 4, 7)
	if sea[4][7] != 3 || err != nil {
		t.Fatalf("Fire(&sea, 4, 7) | sea[4][7] == 3 returned %v, %v, should be 3, nil", sea[4][7], err)

	}

	err = Fire(&sea, 1, 1)
	if err == nil {
		t.Fatalf("Fire(&sea, 1, 1) | sea[1][1] == 42 returned %v, %v, should be *, error", sea[1][1], err)
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
