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
	array [][]int8
}

func (g* grid) init() {
	if (len(g.array) > 0) {
		panic("Array contains data")
	}
	g.array = make([][]int8, g.y)
	for idx := range g.array {
		g.array[idx] = make([]int8, g.x)
	}
}

func (g* grid) randomize(rule_amount int) {
	for idxy := range g.array {
		for idxx := range g.array[idxy] {
			g.array[idxy][idxx] = int8(r.Intn(rule_amount))
		}
	}
}

type rule struct {
	alive bool
	transitions [9]int8
}

func (ru* rule) randomize(rule_amount int) {
	ru.alive = r.Intn(2) == 0
	for idx := range ru.transitions {
		ru.transitions[idx] = int8(r.Intn(rule_amount))
	}
}

type rules struct {
	array []rule
}

func main() {
	g := grid{10,5,nil}
	g.init()
	g.randomize(10)
	r := rule{}
	r.randomize(10)
	fmt.Println(g)
	fmt.Println(r)
}