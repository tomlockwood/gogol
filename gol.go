package gol

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"sync"
	"time"
)

// random generator
var randMutex sync.Mutex
var src = rand.NewSource(time.Now().UnixNano())
var r = rand.New(src)

func randInt(i int) int {
	randMutex.Lock()
	integer := r.Intn(i)
	randMutex.Unlock()
	return integer
}

// Grid is a 2D grid of uint8s
// Represents the game board and alive counts
type Grid struct {
	X, Y  int
	Array [][]uint8
}

// Validate a Grid
func (gr *Grid) Validate() {
	if len(gr.Array) != gr.Y {
		panic(fmt.Sprintf("Grid array length %d does not equal grid y %d", len(gr.Array), gr.Y))
	}

	for idx := range gr.Array {
		if len(gr.Array[idx]) != gr.X {
			panic(fmt.Sprintf("Grid array length at line %d, %d does not equal grid x: %d", idx, len(gr.Array[idx]), gr.X))
		}
	}
}

// MakeGrid creates a Grid with given dimensions
func MakeGrid(x int, y int) Grid {
	array := make([][]uint8, y)
	for idx := range array {
		array[idx] = make([]uint8, x)
	}
	return Grid{x, y, array}
}

// Randomize a Grid based on the amount of Rules
// it represents
func (gr *Grid) Randomize(RuleAmount int) {
	for idxy := range gr.Array {
		for idxx := range gr.Array[idxy] {
			gr.Array[idxy][idxx] = uint8(randInt(RuleAmount))
		}
	}
}

func (gr *Grid) print() {
	for idx := range gr.Array {
		fmt.Println(gr.Array[idx])
	}
}

// Rule with alive status and transitions which
// represent what the rule changes to based on
// amount of adjacent alive cells (0-8)
type Rule struct {
	Alive       bool
	Transitions [9]uint8
}

// Randomize a single Rule
func (ru *Rule) Randomize(RuleAmount int) {
	ru.Alive = randInt(2) == 0
	for idx := range ru.Transitions {
		ru.Transitions[idx] = uint8(randInt(RuleAmount))
	}
}

// Rules is an ordered array of Rule structs
type Rules struct {
	Array []Rule
}

// Randomize an array of Rules
func (rs *Rules) Randomize(RuleAmount int) {
	rs.Array = make([]Rule, RuleAmount)
	for idx := range rs.Array {
		rs.Array[idx].Randomize(RuleAmount)
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

// Game contains all game state required to progress a game of life
type Game struct {
	x, y       int
	grid       Grid
	rules      Rules
	alives     alives
	aliveCount Grid
}

// Validate that a game's contents are consistent
// If this does not pass the game cannot Tick properly
func (g *Game) Validate() {
	// Check grid exists
	if len(g.grid.Array) == 0 {
		panic("Grid not loaded")
	}

	var ruleNumber uint8

	ruleNumber = uint8(len(g.rules.Array))

	// Check rules exist
	if ruleNumber == 0 {
		panic("Rules not loaded")
	}

	// Check grid has no cells outside rule number
	for y := range g.grid.Array {
		for x := range g.grid.Array[y] {
			if g.grid.Array[y][x] > ruleNumber {
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
				g.aliveCount.Array[absoluteY][absoluteX]++
			} else {
				g.aliveCount.Array[absoluteY][absoluteX]--
			}
		}
	}
	g.alives.array[y][x] = aliveState
}

func (g *Game) init() {
	var cellAlive bool
	for y := 0; y < g.y; y++ {
		for x := 0; x < g.x; x++ {
			cellAlive = g.rules.Array[g.grid.Array[y][x]].Alive
			if cellAlive {
				g.updateAliveState(x, y, cellAlive)
			}
		}
	}
}

// Tick progresses the game one step forward
func (g *Game) Tick() {
	var oldCellRule, newCellRule Rule
	var nextRuleIdx uint8
	var cellAlive bool
	oldAliveCount := MakeGrid(g.x, g.y)
	for y := range g.aliveCount.Array {
		for x := range g.aliveCount.Array[y] {
			oldAliveCount.Array[y][x] = g.aliveCount.Array[y][x]
		}
	}
	newGrid := MakeGrid(g.x, g.y)
	for y := 0; y < g.y; y++ {
		for x := 0; x < g.x; x++ {
			oldCellRule = g.rules.Array[g.grid.Array[y][x]]
			nextRuleIdx = oldCellRule.Transitions[oldAliveCount.Array[y][x]]
			newGrid.Array[y][x] = nextRuleIdx
			newCellRule = g.rules.Array[nextRuleIdx]
			cellAlive = newCellRule.Alive
			if cellAlive != g.alives.array[y][x] {
				g.updateAliveState(x, y, cellAlive)
			}
		}
	}
	copy(g.grid.Array, newGrid.Array)
}

// TickFunction is called to Tick a Game
type TickFunction func(g Game, gameNumber int)

// Run a Game
func Run(Options GameOpts, TickFunction TickFunction, gameNumber int) {
	g := MakeGame(Options)
	TickFunction(g, gameNumber)
}

// GameOpts represents all the options necessary to make
// a valid game
type GameOpts struct {
	X, Y       int
	Grid       Grid
	RuleNumber int
	Rules      Rules
}

// MakeGame constructs a game from a given set of options,
// Which may be missing some options
func MakeGame(options GameOpts) Game {

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

	// Create the game object
	currentGame := Game{
		options.X,
		options.Y,
		options.Grid,
		options.Rules,
		alives,
		aliveCounts}

	// Ensure nothing mismatches
	currentGame.Validate()

	// Initialize the game
	currentGame.init()

	return currentGame
}

// RunMany games of life concurrently
// TODO - Get rid of tickAmount and add TickFunction
func RunMany(Options GameOpts, gameAmount int, TickFunction TickFunction) {
	var wg sync.WaitGroup
	wg.Add(gameAmount)
	for i := 0; i < gameAmount; i++ {
		go func(i int) {
			defer wg.Done()
			Run(Options, TickFunction, i)
		}(i)
	}
	wg.Wait()
}

type gameSave struct {
	Rules []Rule    `json:"rules"`
	Grid  [][]uint8 `json:"grid"`
}

// Save game of life to a file
func (g *Game) Save(filename string) {
	saveInfo := gameSave{g.rules.Array, g.grid.Array}
	json, err := json.Marshal(saveInfo)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(filename, json, 0644)
}
