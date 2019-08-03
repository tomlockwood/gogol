package gol

import "testing"

func TestStorage(t *testing.T) {

	grid := MakeGrid(3, 4)

	history := History{[]Grid{}, Grid{}, 3}

	history.Store(grid)
	history.Store(grid)
	history.Store(grid)
	history.Store(grid)
	history.Store(grid)
	grid.Array[0][0] = 69
	history.Store(grid)

	if len(history.Grids) > 3 {
		t.Fatalf("Stored more than history capacity of %d, stored %d", history.Capacity, len(history.Grids))
	}

	if history.Grids[0].Array[0][0] != 69 {
		t.Fatalf("Didn't store state properly in latest storage")
	}

	if history.Grids[1].Array[0][0] != 0 {
		t.Fatalf("Didn't store state properly in older storage")
	}
}
