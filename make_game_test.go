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

var copyOpts Options

func TestDefaultLoad(t *testing.T) {
	copyOpts = opts
	game := MakeGame(copyOpts)

	if game.X != 50 {
		t.Fatalf("X Coordinate not set to default of 50")
	}

	if game.Y != 50 {
		t.Fatalf("Y Coordinate not set to default of 50")
	}

	if len(game.Rules.Array) > 6 || len(game.Rules.Array) < 2 {
		t.Fatalf("Random rule number not within 2-6")
	}
}

func TestGridLoad(t *testing.T) {
	copyOpts = opts
	copyOpts.Grid = MakeGrid(4, 3)
	game := MakeGame(copyOpts)

	if game.X != 4 {
		t.Fatalf("X Coordinate not set to 4")
	}

	if game.Y != 3 {
		t.Fatalf("Y Coordinate not set to 3")
	}
}

func TestGridXYSet(t *testing.T) {
	copyOpts = opts
	copyOpts.Y = 3
	copyOpts.X = 4
	game := MakeGame(copyOpts)

	if len(game.Grid.Array[0]) != 4 {
		t.Fatalf("X Coordinate not set to 4")
	}

	if len(game.Grid.Array) != 3 {
		t.Fatalf("Y Coordinate not set to 3")
	}
}

func TestRulesLoad(t *testing.T) {
	copyOpts = opts
	copyOpts.Rules = Rules{}
	copyOpts.Rules.Randomize(3)
	game := MakeGame(copyOpts)

	if len(game.Rules.Array) != 3 {
		t.Fatalf("Three rules loaded and not retained")
	}
}

func TestRulenumberSet(t *testing.T) {
	copyOpts = opts
	copyOpts.RuleNumber = 3
	game := MakeGame(copyOpts)

	if len(game.Rules.Array) != 3 {
		t.Fatalf("Rulenumber 3 set and not created")
	}
}

func TestGridvsXYMismatch(t *testing.T) {
	copyOpts = opts
	copyOpts.Grid = MakeGrid(4, 3)
	copyOpts.X = 10
	copyOpts.Y = 10

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("A mismatch between grid size and set X and Y did not cause an error")
		}
	}()

	MakeGame(copyOpts)
}

func TestRulesvsRulenumberMismatch(t *testing.T) {
	copyOpts = opts
	copyOpts.Rules = Rules{}
	copyOpts.Rules.Randomize(3)
	copyOpts.RuleNumber = 10

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("A mismatch between rules size and rulenumber did not cause an error")
		}
	}()

	MakeGame(copyOpts)
}

func TestRulesvsGridContentMismatch(t *testing.T) {
	copyOpts = opts
	copyOpts.Rules = Rules{}
	copyOpts.Rules.Randomize(3)
	copyOpts.Grid = MakeGrid(3, 3)
	copyOpts.Grid.Array[0][0] = 8

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("A mismatch between rules size and the grid content did not cause an error")
		}
	}()

	MakeGame(copyOpts)
}
