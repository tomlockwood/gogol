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

func (g* game) init_alives() {
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

func main() {
	x := 5
	y := 5
	gr := grid{x,y,nil}
	gr.init()
	gr.randomize(2)
	rs := rules{}
	rs.randomize(2)
	a := alives{x,y,nil}
	a.init()
	ac := grid{x,y,nil}
	ac.init()
	g := game{x,y,gr,rs,a,ac}
	g.init_alives()
	fmt.Println(g)
	fmt.Println("GRID")
	g.grid.print()
	fmt.Println("ALIVE COUNT")
	g.alive_count.print()
}