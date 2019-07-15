package gol

// Conway's Game of Life Testing

// Testing that the old classic game of life rules function correctly

import (
	"fmt"
	"testing"
)

func matchSlice(a [][]uint8, b [][]uint8) bool {
	for y := range a {
		for x := range a[y] {
			if a[y][x] != b[y][x] {
				fmt.Println(x, y)
				return false
			}
		}
	}
	return true
}

func printArray(array [][]uint8) {
	for idx := range array {
		fmt.Println(array[idx])
	}
}

func mismatchCheck(expected [][]uint8, got [][]uint8) bool {
	if !matchSlice(expected, got) {
		fmt.Println("Expected:")
		printArray(expected)
		fmt.Println("Got:")
		printArray(got)
		return true
	}
	return false
}

// Standard conway's game of life ruleset

var r0 = Rule{false, [9]uint8{0, 0, 0, 1, 0, 0, 0, 0, 0}}
var r1 = Rule{true, [9]uint8{0, 0, 1, 1, 0, 0, 0, 0, 0}}
var rs = Rules{[]Rule{r0, r1}}

// TestConwayDeath - Checking one Tick death
func TestConwayDeath(t *testing.T) {
	y0 := []uint8{0, 0, 0}
	y1 := []uint8{0, 1, 0}
	y2 := []uint8{0, 1, 0}
	array := [][]uint8{y0, y1, y2}
	grid := Grid{3, 3, array}

	game := MakeGame(GameOpts{3, 3, grid, 2, rs})
	game.Tick()

	y0Out := []uint8{0, 0, 0}
	y1Out := []uint8{0, 0, 0}
	y2Out := []uint8{0, 0, 0}
	arrayOut := [][]uint8{y0Out, y1Out, y2Out}

	if mismatchCheck(arrayOut, game.Grid.Array) {
		t.Fatalf("Slices do not match")
	}
}

// TestConwaySquare - Checking still life
func TestConwaySquare(t *testing.T) {
	y0 := []uint8{0, 0, 0, 0}
	y1 := []uint8{0, 1, 1, 0}
	y2 := []uint8{0, 1, 1, 0}
	y3 := []uint8{0, 0, 0, 0}
	array := [][]uint8{y0, y1, y2, y3}
	grid := Grid{4, 4, array}

	game := MakeGame(GameOpts{4, 4, grid, 2, rs})
	game.Tick()

	y0Out := []uint8{0, 0, 0, 0}
	y1Out := []uint8{0, 1, 1, 0}
	y2Out := []uint8{0, 1, 1, 0}
	y3Out := []uint8{0, 0, 0, 0}
	arrayOut := [][]uint8{y0Out, y1Out, y2Out, y3Out}

	if mismatchCheck(arrayOut, game.Grid.Array) {
		t.Fatalf("Slices do not match")
	}
}

// TestConwayGlider - Testing gliders work
func TestConwayGlider(t *testing.T) {
	// Test for 1 Tick death
	y0 := []uint8{0, 1, 0, 0, 0, 0, 0, 0, 0, 0}
	y1 := []uint8{0, 0, 1, 0, 0, 0, 0, 0, 0, 0}
	y2 := []uint8{1, 1, 1, 0, 0, 0, 0, 0, 0, 0}
	y3 := []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	y4 := []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	y5 := []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	y6 := []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	y7 := []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	y8 := []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	y9 := []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	array := [][]uint8{y0, y1, y2, y3, y4, y5, y6, y7, y8, y9}
	grid := Grid{10, 10, array}

	g := MakeGame(GameOpts{10, 10, grid, 2, rs})
	for i := 0; i <= 23; i++ {
		g.Tick()
	}

	y0Out := []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	y1Out := []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	y2Out := []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	y3Out := []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	y4Out := []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	y5Out := []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	y6Out := []uint8{0, 0, 0, 0, 0, 0, 0, 1, 0, 0}
	y7Out := []uint8{0, 0, 0, 0, 0, 0, 0, 0, 1, 0}
	y8Out := []uint8{0, 0, 0, 0, 0, 0, 1, 1, 1, 0}
	y9Out := []uint8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	arrayOut := [][]uint8{y0Out, y1Out, y2Out, y3Out, y4Out, y5Out, y6Out, y7Out, y8Out, y9Out}

	if mismatchCheck(arrayOut, g.Grid.Array) {
		t.Fatalf("Slices do not match")
	}
}
