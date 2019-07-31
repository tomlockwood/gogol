package gol

import "testing"

// MakeGame testing

// Checking that permutations of the rules work

var opts = Options{
	0,
	0,
	Grid{},
	0,
	Rules{}}

func TestDefaultMakeGame(t *testing.T) {
	game := MakeGame(opts)

	if game.X != 50 {
		t.Fatalf("X Coordinate not set to default of 50")
	}

	if game.Y != 50 {
		t.Fatalf("Y Coordinate not set to default of 50")
	}
}
