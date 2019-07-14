package gol

import (
	"fmt"
	"math/rand"
	"time"
)

// random generator
var src = rand.NewSource(time.Now().UnixNano())
var r = rand.New(src)

// Grid to play on
type Grid struct {
	x, y  int
	array [][]uint8
}

func (gr *Grid) validate() {
	if len(gr.array) != gr.y {
		panic(fmt.Sprintf("Grid array length %d does not equal grid y %d", len(gr.array), gr.y))
	}

	for idx := range gr.array {
		if len(gr.array[idx]) != gr.x {
			panic(fmt.Sprintf("Grid array length at line %d, %d does not equal grid x: %d", idx, len(gr.array[idx]), gr.x))
		}
	}
}

func makeGrid(x int, y int) Grid {
	array := make([][]uint8, y)
	for idx := range array {
		array[idx] = make([]uint8, x)
	}
	return Grid{x, y, array}
}

func (gr *Grid) randomize(ruleAmount int) {
	for idxy := range gr.array {
		for idxx := range gr.array[idxy] {
			gr.array[idxy][idxx] = uint8(r.Intn(ruleAmount))
		}
	}
}

func (gr *Grid) print() {
	for idx := range gr.array {
		fmt.Println(gr.array[idx])
	}
}

type rule struct {
	alive       bool
	transitions [9]uint8
}

func (ru *rule) randomize(ruleAmount int) {
	ru.alive = r.Intn(2) == 0
	for idx := range ru.transitions {
		ru.transitions[idx] = uint8(r.Intn(ruleAmount))
	}
}

// Rules for game of life
type Rules struct {
	array []rule
}

func (rs *Rules) randomize(ruleAmount int) {
	rs.array = make([]rule, ruleAmount)
	for idx := range rs.array {
		rs.array[idx].randomize(ruleAmount)
	}
}

type alives struct {
	x, y  int
	array [][]bool
}

func makeAlives(x int, y int) alives {
	array := make([][]bool, y)
	for idx := range array {
		array[idx] = make([]bool, x)
	}
	return alives{x, y, array}
}

// Game for fun
type Game struct {
	x, y       int
	grid       Grid
	rules      Rules
	alives     alives
	aliveCount Grid
}

func (g *Game) validate() {
	// Check grid exists
	if len(g.grid.array) == 0 {
		panic("Grid not loaded")
	}

	var ruleNumber uint8

	ruleNumber = uint8(len(g.rules.array))

	// Check rules exist
	if ruleNumber == 0 {
		panic("Rules not loaded")
	}

	// Check grid has no cells outside rule number
	for y := range g.grid.array {
		for x := range g.grid.array[y] {
			if g.grid.array[y][x] > ruleNumber {
				panic(fmt.Sprintf("X: %d Y: %d not consistent with rule count", x, y))
			}
		}
	}
}

func (g *Game) updateAliveState(x int, y int, aliveState bool) {
	var absoluteY, absoluteX int
	for relY := -1; relY <= 1; relY++ {
		for relX := -1; relX <= 1; relX++ {
			absoluteY = relY + y
			absoluteX = relX + x
			if (relY == 0 && relX == 0) || absoluteY < 0 || absoluteX < 0 || absoluteY >= g.y || absoluteX >= g.x {
				continue
			}
			if aliveState {
				g.aliveCount.array[absoluteY][absoluteX]++
			} else {
				g.aliveCount.array[absoluteY][absoluteX]--
			}
		}
	}
	g.alives.array[y][x] = aliveState
}

func (g *Game) init() {
	var cellAlive bool
	for y := 0; y < g.y; y++ {
		for x := 0; x < g.x; x++ {
			cellAlive = g.rules.array[g.grid.array[y][x]].alive
			if cellAlive {
				g.updateAliveState(x, y, cellAlive)
			}
		}
	}
}

func (g *Game) tick() {
	var oldCellRule, newCellRule rule
	var nextRuleIdx uint8
	var cellAlive bool
	oldAliveCount := makeGrid(g.x, g.y)
	for y := range g.aliveCount.array {
		for x := range g.aliveCount.array[y] {
			oldAliveCount.array[y][x] = g.aliveCount.array[y][x]
		}
	}
	newGrid := makeGrid(g.x, g.y)
	for y := 0; y < g.y; y++ {
		for x := 0; x < g.x; x++ {
			oldCellRule = g.rules.array[g.grid.array[y][x]]
			nextRuleIdx = oldCellRule.transitions[oldAliveCount.array[y][x]]
			newGrid.array[y][x] = nextRuleIdx
			newCellRule = g.rules.array[nextRuleIdx]
			cellAlive = newCellRule.alive
			if cellAlive != g.alives.array[y][x] {
				g.updateAliveState(x, y, cellAlive)
			}
		}
	}
	copy(g.grid.array, newGrid.array)
}

// Run a game of life
func (g *Game) Run(count int, interactive bool) {
	var response int
	for i := 0; i <= count; i++ {

		if interactive {
			g.grid.print()
			fmt.Scanf("%c", &response)
			fmt.Println()
		}
		g.tick()
	}
}

// GameOpts of life Options
type GameOpts struct {
	X, Y       int
	Grid       Grid
	RuleNumber int
	Rules      Rules
}

// MakeGame is for fun
func MakeGame(options GameOpts) Game {

	// Get/set rules amount if needed
	if options.Rules.array == nil {
		if options.RuleNumber == 0 {
			options.RuleNumber = r.Intn(4) + 2
		}
		options.Rules.randomize(options.RuleNumber)
	} else if options.RuleNumber == 0 {
		options.RuleNumber = len(options.Rules.array)
	} else if options.RuleNumber != len(options.Rules.array) {
		panic(fmt.Sprintf("Rule number in options %d does not equal rules in array %d", options.RuleNumber, len(options.Rules.array)))
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
		if options.Grid.array == nil {
			setX = 50
		} else if options.Grid.x == 0 {
			setX = len(options.Grid.array[0])
		} else {
			setX = options.Grid.x
		}
	} else {
		setX = options.X
	}

	if options.Y == 0 {
		if options.Grid.array == nil {
			setY = 50
		} else if options.Grid.y == 0 {
			setY = len(options.Grid.array)
		} else {
			setY = options.Grid.y
		}
	} else {
		setY = options.Y
	}

	options.Y = setY
	options.X = setX

	// Create the grid if it doesn't exist
	// or validate the grid
	if options.Grid.array == nil {
		gameGrid = makeGrid(options.X, options.Y)
		gameGrid.randomize(options.RuleNumber)
	} else {
		options.Grid.x = options.X
		options.Grid.y = options.Y
		options.Grid.validate()
		gameGrid = options.Grid
	}

	options.Grid = gameGrid

	// Make alive bool and counts arrays
	alives := makeAlives(options.X, options.Y)
	aliveCounts := makeGrid(options.X, options.Y)

	currentGame := Game{
		options.X,
		options.Y,
		options.Grid,
		options.Rules,
		alives,
		aliveCounts}

	currentGame.validate()
	currentGame.init()

	return currentGame
}
