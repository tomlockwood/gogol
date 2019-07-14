package gol

import (
	"fmt"
	"math/rand"
	"time"
)

// random generator
var src = rand.NewSource(time.Now().UnixNano())
var r = rand.New(src)

type grid struct {
	x, y  int
	array [][]uint8
}

func makeGrid(x int, y int) grid {
	array := make([][]uint8, y)
	for idx := range array {
		array[idx] = make([]uint8, x)
	}
	return grid{x, y, array}
}

func (gr *grid) randomize(ruleAmount int) {
	for idxy := range gr.array {
		for idxx := range gr.array[idxy] {
			gr.array[idxy][idxx] = uint8(r.Intn(ruleAmount))
		}
	}
}

func (gr *grid) print() {
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

type rules struct {
	array []rule
}

func (rs *rules) randomize(ruleAmount int) {
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

type game struct {
	x, y       int
	grid       grid
	rules      rules
	alives     alives
	aliveCount grid
}

func (g *game) updateAliveState(x int, y int, aliveState bool) {
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

func (g *game) init() {
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

func (g *game) tick() {
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

func (g *game) run(count int, interactive bool) {
	var response int
	for i := 0; i <= count; i++ {
		if interactive {
			fmt.Scanf("%c", &response)
		}
		g.grid.print()
		fmt.Println()
		g.tick()
	}
}

func main() {
	x := 10
	y := 10
	gr := makeGrid(x, y)
	gr.randomize(2)
	r0 := rule{false, [9]uint8{0, 0, 0, 1, 0, 0, 0, 0, 0}}
	r1 := rule{true, [9]uint8{0, 0, 1, 1, 0, 0, 0, 0, 0}}
	rs := rules{[]rule{r0, r1}}
	a := makeAlives(x, y)
	ac := makeGrid(x, y)
	g := game{x, y, gr, rs, a, ac}
	g.init()
	g.run(10000, false)
}
