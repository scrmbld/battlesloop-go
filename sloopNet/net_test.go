package main

import (
	"testing"
)

// TODO: make sure the parser can handle invalid message types

func TestParsePositional(t *testing.T) {

	/* Ensure parsePositional does the following things:
	1. gives correct output on valid inputs, as defined in the Battlesloop Protocol spec
	2. handles invalid row or column values
	*/
	var message string = "_h_A-7" // hit on A7
	result, err := parsePositional(message)
	expected := PositionalMsg{"h", 7, 0}
	if result != expected || err != nil {
		t.Fatalf("parsePositional(%s) returned %v, %v | expected %v, %v", message, result, err, expected, nil)
	}

	message = "_m_B-5" // miss on B5
	result, err = parsePositional(message)
	expected = PositionalMsg{"m", 5, 1}
	if result != expected || err != nil {
		t.Fatalf("parsePositional(%s) returned %v, %v | expected %v %v", message, result, err, expected, nil)
	}

	message = "_f_J-9" // fire at I9
	result, err = parsePositional(message)
	expected = PositionalMsg{"f", 9, 9}
	if result != expected || err != nil {
		t.Fatalf("parsePositional(%s) returned %v, %v | expected %v %v", message, result, err, expected, nil)
	}

	message = "_f_f-9" // fire at f9 -- 'f' is an invalid column, so this should return an error
	result, err = parsePositional(message)
	expected = PositionalMsg{"", 0, 0}
	if result != expected || err == nil {
		t.Fatalf("parsePositional(%s) returned %v, %v | expected %v Invalid column error", message, result, err, expected)
	}

	message = "_f_F-370" // fire at F370 -- row is out of range, so this should return an error
	result, err = parsePositional(message)
	expected = PositionalMsg{"", 0, 0}
	if result != expected || err == nil {
		t.Fatalf("parsePositional(%s) returned %v, %v | expected %v Row out of range error", message, result, err, expected)
	}

	message = "_f_F-A" // fire at FA -- A is not a valid row because it is not a number
	result, err = parsePositional(message)
	expected = PositionalMsg{"", 0, 0}
	if result != expected || err == nil {
		t.Fatalf("parsePositional(%s) returned %v, %v | expected %v Atoi() invalid syntax error", message, result, err, expected)
	}
}

func TestParseShip(t *testing.T) {

	message := "_a_7" // sank ship 7
	result, err := parseShip(message)
	expected := ShipMsg{"a", 7}
	if result != expected || err != nil {
		t.Fatalf("parseShip(%s) returned %v, %v | expected %v, %v", message, result, err, expected, nil)
	}

	message = "_a_g" // sank ship g -- illegal ship id, should return an error
	result, err = parseShip(message)
	expected = ShipMsg{}
	if result != expected || err == nil {
		t.Fatalf("parseShip(%s) returned %v, %v | expected %v, Atoi() invalid syntax error", message, result, err, expected)
	}
}

func TestParseGame(t *testing.T) {
	// TODO: enforce better input sanitization
	message := "_g_win"
	result, err := parseGame(message)
	expected := GameMsg{"win"}
	if result != expected || err != nil {
		t.Fatalf("parseGame(%s) returned %v, %v | expected %v, %v", message, result, err, expected, nil)
	}
}

func TestParseConn(t *testing.T) {
	// TODO: enforce better input sanitization
	message := "_c_begin"
	result, err := parseConn(message)
	expected := ConnMsg{"begin"}
	if result != expected || err != nil {
		t.Fatalf("parseGame(%s) returned %v, %v | expected %v, %v", message, result, err, expected, nil)
	}
}
