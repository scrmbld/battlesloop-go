package sloopNet

import (
	"testing"
)

// TODO: make sure the parser can handle invalid message types

func msgEq(s1 [][]string, s2 [][]string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if len(s1) != len(s2) {
			return false
		}
		for j := range s1[i] {
			if s1[i][j] != s2[i][j] {
				return false
			}
		}
	}

	return true
}

func TestParsePos(t *testing.T) {
	var message string = "A-7"
	result, err := ParsePos(message)
	expected := [2]uint8{7, 0}
	if result != expected || err != nil {
		t.Fatalf("ParsePos(%s) returned %v, %v | expected %v, %v", message, result, err, expected, nil)
	}

	message = "B-5"
	result, err = ParsePos(message)
	expected = [2]uint8{5, 1}
	if result != expected || err != nil {
		t.Fatalf("ParsePos(%s) returned %v, %v | expected %v, %v", message, result, err, expected, nil)
	}

	message = "J-9"
	result, err = ParsePos(message)
	expected = [2]uint8{9, 9}
	if result != expected || err != nil {
		t.Fatalf("ParsePos(%s) returned %v, %v | expected %v, %v", message, result, err, expected, nil)
	}

	message = "f-9" // 'f' is an invalid column, so this should return an error
	result, err = ParsePos(message)
	expected = [2]uint8{0, 0}
	if result != expected || err == nil {
		t.Fatalf("ParsePos(%s) returned %v, %v | expected %v Invalid column error", message, result, err, expected)
	}

	message = "F-370" // fire at F370 -- row is out of range, so this should return an error
	result, err = ParsePos(message)
	expected = [2]uint8{0, 0}
	if result != expected || err == nil {
		t.Fatalf("ParsePos(%s) returned %v, %v | expected %v Row out of range error", message, result, err, expected)
	}

	message = "F-A" // A is not a valid row because it is not a number
	result, err = ParsePos(message)
	expected = [2]uint8{0, 0}
	if result != expected || err == nil {
		t.Fatalf("ParsePos(%s) returned %v, %v | expected %v Atoi() invalid syntax error", message, result, err, expected)
	}
}

func TestParseMsgs(t *testing.T) {

	var message string = "_g_begin:_f_A-7"
	result, err := ParseMsgs(message)
	expected := [][]string{{"g", "begin"}, {"f", "A-7"}}
	if !msgEq(expected, result) || err != nil {
		t.Fatalf("ParseMsgs(%s) returned %v, %v | expected %v, %v", message, result, err, expected, nil)
	}

	message = "_g_begin"
	result, err = ParseMsgs(message)
	expected = [][]string{{"g", "begin"}}
	if !msgEq(expected, result) || err != nil {
		t.Fatalf("ParseMsgs(%s) returned %v, %v | expected %v, %v", message, result, err, expected, nil)
	}

	// some invalid messages
	message = "begin"
	result, err = ParseMsgs(message)
	expected = [][]string{}
	if !msgEq(expected, result) || err == nil {
		t.Fatalf("ParseMsgs(%s) returned %v, %v | expected %v, Invalid message: has no type or too many types", message, result, err, expected)
	}

	message = "_g__f_begin"
	result, err = ParseMsgs(message)
	expected = [][]string{}
	if !msgEq(expected, result) || err == nil {
		t.Fatalf("ParseMsgs(%s) returned %v, %v | expected %v, Invalid message: has no type or too many types", message, result, err, expected)
	}

	message = "_g_begin:"
	result, err = ParseMsgs(message)
	expected = [][]string{}
	if !msgEq(expected, result) || err == nil {
		t.Fatalf("ParseMsgs(%s) returned %v, %v | expected %v, Invalid message: has no type or too many types", message, result, err, expected)
	}
}
