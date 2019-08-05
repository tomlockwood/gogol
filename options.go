package gol

import "fmt"

// Options represents all the options necessary to make
// a valid game
type Options struct {
	X, Y       int
	Grid       Grid
	RuleNumber int
	Rules      Rules
}

// MakeGame constructs a game from a given set of options,
// Which may be missing some options
func MakeGame(options Options) Game {

	// Get/set rules amount if needed
	if options.Rules.Array == nil {
		if options.RuleNumber == 0 {
			options.RuleNumber = randInt(4) + 2
		}
		options.Rules.Randomize(options.RuleNumber)
	} else if options.RuleNumber == 0 {
		options.RuleNumber = len(options.Rules.Array)
	} else if options.RuleNumber != len(options.Rules.Array) {
		panic(fmt.Sprintf("Rule number in options %d does not equal rules in array %d", options.RuleNumber, len(options.Rules.Array)))
	}

	// Grid check
	if options.X < 0 || options.Y < 0 {
		panic("X/Y values cannot be negative")
	}

	var setX, setY int
	var gameGrid Grid

	// Make the proposed game x, y 50 by default
	// or based on the length of the grid slices
	// or by the x, y set by the grid
	// or by the game options struct
	if options.X == 0 {
		if options.Grid.Array == nil {
			setX = 50
		} else if options.Grid.X == 0 {
			setX = len(options.Grid.Array[0])
		} else {
			setX = options.Grid.X
		}
	} else {
		setX = options.X
	}

	if options.Y == 0 {
		if options.Grid.Array == nil {
			setY = 50
		} else if options.Grid.Y == 0 {
			setY = len(options.Grid.Array)
		} else {
			setY = options.Grid.Y
		}
	} else {
		setY = options.Y
	}

	options.Y = setY
	options.X = setX

	// Create the grid if it doesn't exist
	// or validate the grid
	if options.Grid.Array == nil {
		gameGrid = MakeGrid(options.X, options.Y)
		gameGrid.Randomize(options.RuleNumber)
	} else {
		options.Grid.X = options.X
		options.Grid.Y = options.Y
		options.Grid.Validate()
		gameGrid = options.Grid
	}

	options.Grid = gameGrid

	// Make alive bool and counts arrays
	alives := makeAlives(options.X, options.Y)
	aliveCounts := MakeGrid(options.X, options.Y)
	field := MakeGridBuffers(options.X, options.Y, false)

	field.Front = options.Grid.Array
	// Create the game object
	currentGame := Game{
		options.X,
		options.Y,
		field,
		options.Rules,
		alives,
		aliveCounts,
		0}

	// Ensure nothing mismatches
	currentGame.Validate()

	// Initialize the game
	currentGame.init()

	return currentGame
}
