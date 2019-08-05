package gol

import (
	"fmt"
	"sync"
)

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

// CopyGrid creates a copy of a Grid
func CopyGrid(g Grid) Grid {
	array := make([][]uint8, g.Y)
	for idx := range array {
		array[idx] = make([]uint8, g.X)
		copy(array[idx], g.Array[idx])
	}
	return Grid{g.X, g.Y, array}
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

// Print grid contents
func (gr *Grid) Print() {
	for idx := range gr.Array {
		fmt.Println(gr.Array[idx])
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

// GridBuffers are a set of Grids for flippin'
type GridBuffers struct {
	X, Y        int
	Front, back [][]uint8
	mutexes     [][]sync.Mutex
}

// MakeGridBuffers from x,y ranges
func MakeGridBuffers(x int, y int, lockable bool) GridBuffers {
	front := make([][]uint8, y)
	back := make([][]uint8, y)
	var mutexes [][]sync.Mutex
	if lockable {
		mutexes = make([][]sync.Mutex, y)
	} else {
		mutexes = [][]sync.Mutex{}
	}
	for idx := range front {
		front[idx] = make([]uint8, x)
		back[idx] = make([]uint8, x)
		if lockable {
			mutexes[idx] = make([]sync.Mutex, x)
		}
	}
	return GridBuffers{x, y, front, back, mutexes}
}

func (grb *GridBuffers) flip() {
	grb.back, grb.Front = grb.Front, grb.back
}

// Randomize a Grid based on the amount of Rules
// it represents
func (grb *GridBuffers) Randomize(RuleAmount int) {
	for idxy := range grb.Front {
		for idxx := range grb.Front[idxy] {
			grb.Front[idxy][idxx] = uint8(randInt(RuleAmount))
		}
	}
}
