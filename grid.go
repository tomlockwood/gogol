package gol

import (
	"fmt"
	"sync"
)

// MakeGrid creates a Grid with given dimensions
func MakeGrid(x int, y int) [][]uint8 {
	array := make([][]uint8, y)
	for idx := range array {
		array[idx] = make([]uint8, x)
	}
	return array
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

// CheckGrid is up to spec
func CheckGrid(grid [][]uint8, x int, y int) {
	if len(grid) != y {
		panic(fmt.Sprintf("Grid array length %d does not equal grid y %d", len(grid), y))
	}

	for idx := range grid {
		if len(grid[idx]) != x {
			panic(fmt.Sprintf("Grid array length at line %d, %d does not equal grid x: %d", idx, len(grid), x))
		}
	}
}

// CheckBoolGrid is up to spec
func CheckBoolGrid(grid [][]sync.Mutex, x int, y int) {
	if len(grid) != y {
		panic(fmt.Sprintf("Grid array length %d does not equal grid y %d", len(grid), y))
	}

	for idx := range grid {
		if len(grid[idx]) != x {
			panic(fmt.Sprintf("Grid array length at line %d, %d does not equal grid x: %d", idx, len(grid), x))
		}
	}
}

// Validate all grids in GridBuffers
func (grb *GridBuffers) Validate() {
	CheckGrid(grb.Front, grb.X, grb.Y)
	CheckGrid(grb.back, grb.X, grb.Y)
	if len(grb.mutexes) != 0 {
		CheckBoolGrid(grb.mutexes, grb.X, grb.Y)
	}
}

// CopyFrontToBack copies the front GridBuffer to the back one
func (grb *GridBuffers) CopyFrontToBack() {
	for idx := range grb.back {
		copy(grb.back[idx], grb.Front[idx])
	}
}

// Print the GridBuffers Arrays
func (grb *GridBuffers) Print() {
	fmt.Println("Front Field")
	printArray(grb.Front)
	fmt.Println("Back Field")
	printArray(grb.back)
}
