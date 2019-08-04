package gol

import "fmt"

// Game contains all game state required to progress a game of life
type Game struct {
	X, Y                int
	Grids               [2]Grid
	Rules               Rules
	alives              alives
	aliveCount          Grid
	ticks               int
	FrontGrid, backGrid *Grid
}

// Validate that a game's contents are consistent
// If this does not pass the game cannot Tick properly
func (g *Game) Validate() {
	// Check grid exists
	if len(g.FrontGrid.Array) == 0 {
		panic("Grid not loaded")
	}

	var ruleNumber uint8

	ruleNumber = uint8(len(g.Rules.Array))

	// Check rules exist
	if ruleNumber == 0 {
		panic("Rules not loaded")
	}

	// Check grid has no cells outside rule number
	for y := range g.FrontGrid.Array {
		for x := range g.FrontGrid.Array[y] {
			if g.FrontGrid.Array[y][x] > ruleNumber {
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
			if (relY == 0 && relX == 0) || absoluteY < 0 || absoluteX < 0 || absoluteY >= g.Y || absoluteX >= g.X {
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
	g.ticks = 0
	var cellAlive bool
	for y := 0; y < g.Y; y++ {
		for x := 0; x < g.X; x++ {
			cellAlive = g.Rules.Array[g.FrontGrid.Array[y][x]].Alive
			if cellAlive {
				g.updateAliveState(x, y, cellAlive)
			}
		}
	}
}

// Reset the game to a random initial state
// But with the same rules
func (g *Game) Reset() {
	g.ticks = 0
	g.FrontGrid.Randomize(len(g.Rules.Array))
	g.alives = makeAlives(g.X, g.Y)
	g.aliveCount = MakeGrid(g.X, g.Y)
	g.init()
}

func (g *Game) flipGrid() {
	if (g.ticks % 2) == 0 {
		g.FrontGrid = &g.Grids[0]
		g.backGrid = &g.Grids[1]
	} else {
		g.FrontGrid = &g.Grids[1]
		g.backGrid = &g.Grids[0]
	}
}

// Tick progresses the game one step forward
func (g *Game) Tick() {
	var oldCellRule, newCellRule Rule
	var nextRuleIdx uint8
	var cellAlive bool
	oldAliveCount := CopyGrid(g.aliveCount)
	for y := 0; y < g.Y; y++ {
		for x := 0; x < g.X; x++ {
			oldCellRule = g.Rules.Array[g.FrontGrid.Array[y][x]]
			nextRuleIdx = oldCellRule.Transitions[oldAliveCount.Array[y][x]]
			g.backGrid.Array[y][x] = nextRuleIdx
			newCellRule = g.Rules.Array[nextRuleIdx]
			cellAlive = newCellRule.Alive
			if cellAlive != g.alives.array[y][x] {
				g.updateAliveState(x, y, cellAlive)
			}
		}
	}
	g.ticks++
	g.flipGrid()
}
