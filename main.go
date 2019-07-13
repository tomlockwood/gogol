package main

import (
	"fmt"
	"math/rand"
	"time"
)

// random generator
var src = rand.NewSource(time.Now().UnixNano())
var r = rand.New(src)

type grid struct {
	x,y int
	array [][]uint8
}

func (gr* grid) init() {
	if (len(gr.array) > 0) {
		panic("Array contains data")
	}
	gr.array = make([][]uint8, gr.y)
	for idx := range gr.array {
		gr.array[idx] = make([]uint8, gr.x)
	}
}

func (gr* grid) randomize(rule_amount int) {
	for idxy := range gr.array {
		for idxx := range gr.array[idxy] {
			gr.array[idxy][idxx] = uint8(r.Intn(rule_amount))
		}
	}
}

func (gr* grid) print() {
	for idx := range gr.array {
		fmt.Println(gr.array[idx])
	}
}

type rule struct {
	alive bool
	transitions [9]uint8
}

func (ru* rule) randomize(rule_amount int) {
	ru.alive = r.Intn(2) == 0
	for idx := range ru.transitions {
		ru.transitions[idx] = uint8(r.Intn(rule_amount))
	}
}

type rules struct {
	array []rule
}

func (rs* rules) randomize(rule_amount int) {
	rs.array = make([]rule, rule_amount)
	for idx := range rs.array {
		rs.array[idx].randomize(rule_amount)
	}
}

type alives struct {
	x,y int
	array [][]bool
}

func (a* alives) init() {
	if (len(a.array) > 0) {
		panic("Array contains data")
	}
	a.array = make([][]bool, a.y)
	for idx := range a.array {
		a.array[idx] = make([]bool, a.x)
	}
}

type game struct {
	x,y int
	grid grid
	rules rules
	alives alives
	alive_count grid
}

func (g* game) update_alive_state(x int, y int, alive_state bool) {
	var absolute_y, absolute_x int
	for rel_y := -1; rel_y <= 1; rel_y++ {
		for rel_x := -1; rel_x <= 1; rel_x++ {
			absolute_y = rel_y + y
			absolute_x = rel_x + x
			if (rel_y == 0 && rel_x == 0) || absolute_y < 0 || absolute_x < 0 || absolute_y >= g.y || absolute_x >= g.x {
				continue
			}
			if (alive_state) {
				g.alive_count.array[absolute_y][absolute_x]++
			} else {
				g.alive_count.array[absolute_y][absolute_x]--
			}
		}
	}
	g.alives.array[y][x] = alive_state
}

func (g* game) init() {
	var cell_alive bool
	for y := 0; y < g.y; y++ {
		for x := 0; x < g.x; x++ {
			cell_alive = g.rules.array[g.grid.array[y][x]].alive
			if (cell_alive) {
				g.update_alive_state(x,y,cell_alive)
			}
		}
	}
}

func (g* game) tick() {
	var cell_rule rule
	var cell_alive bool
	old_alive_count := grid{g.x,g.y,nil}
	old_alive_count.init()
	for y := range g.alive_count.array {
		for x := range g.alive_count.array[y] {
			old_alive_count.array[y][x] = g.alive_count.array[y][x]
		}
	}
	new_grid := grid{g.x,g.y,nil}
	new_grid.init()
	for y := 0; y < g.y; y++ {
		for x := 0; x < g.x; x++ {
			cell_rule = g.rules.array[g.grid.array[y][x]]
			new_grid.array[y][x] = cell_rule.transitions[old_alive_count.array[y][x]]
			cell_alive = cell_rule.alive
			if (cell_alive != g.alives.array[y][x]) {
				g.update_alive_state(x,y,cell_alive)
			}
		}
	}
	copy(g.grid.array,new_grid.array)
}

func (g* game) run(count int, interactive bool) {
	var response int
	for i := 0; i <= count; i++ {
		if (interactive) {
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
	gr := grid{x,y,nil}
	gr.init()
	gr.randomize(2)
	r0 := rule{false,[9]uint8{0,0,0,1,0,0,0,0,0}}
	r1 := rule{true,[9]uint8{0,0,1,1,0,0,0,0,0}}
	rs := rules{[]rule{r0,r1}}
	a := alives{x,y,nil}
	a.init()
	ac := grid{x,y,nil}
	ac.init()
	g := game{x,y,gr,rs,a,ac}
	g.init()
	g.run(10000,false)
}