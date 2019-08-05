package gol

import "fmt"

// Game contains all game state required to progress a game of life
type Game struct {
	X, Y       int
	Field      GridBuffers
	Rules      Rules
	alives     alives
	aliveCount GridBuffers
	ticks      int
}

// Validate that a game's contents are consistent
// If this does not pass the game cannot Tick properly
func (g *Game) Validate() {
	// Check grid exists
	if len(g.Field.Front) == 0 {
		panic("Grid not loaded")
	}

	var ruleNumber uint8

	ruleNumber = uint8(len(g.Rules.Array))

	// Check rules exist
	if ruleNumber == 0 {
		panic("Rules not loaded")
	}

	// Check grid has no cells outside rule number
	for y := range g.Field.Front {
		for x := range g.Field.Front[y] {
			if g.Field.Front[y][x] > ruleNumber {
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
				g.aliveCount.back[absoluteY][absoluteX]++
			} else {
				g.aliveCount.back[absoluteY][absoluteX]--
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
			cellAlive = g.Rules.Array[g.Field.Front[y][x]].Alive
			if cellAlive {
				g.updateAliveState(x, y, cellAlive)
			}
		}
	}
	g.aliveCount.flip()
}

// Reset the game to a random initial state
// But with the same rules
func (g *Game) Reset() {
	g.ticks = 0
	g.Field.Randomize(len(g.Rules.Array))
	g.alives = makeAlives(g.X, g.Y)
	g.aliveCount = MakeGridBuffers(g.X, g.Y, true)
	g.init()
}

// Tick progresses the game one step forward
func (g *Game) Tick() {
	var oldCellRule, newCellRule Rule
	var nextRuleIdx uint8
	var cellAlive bool
	g.aliveCount.CopyFrontToBack()
	for y := 0; y < g.Y; y++ {
		for x := 0; x < g.X; x++ {
			oldCellRule = g.Rules.Array[g.Field.Front[y][x]]
			nextRuleIdx = oldCellRule.Transitions[g.aliveCount.Front[y][x]]
			g.Field.back[y][x] = nextRuleIdx
			newCellRule = g.Rules.Array[nextRuleIdx]
			cellAlive = newCellRule.Alive
			if cellAlive != g.alives.array[y][x] {
				g.updateAliveState(x, y, cellAlive)
			}
		}
	}
	g.ticks++
	g.Field.flip()
	g.aliveCount.flip()
}
