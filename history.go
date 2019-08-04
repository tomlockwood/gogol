package gol

// History stores the history of a game and allows analysis
type History struct {
	Grids    []Grid
	Seed     Grid
	Capacity int
}

// Store history
func (h *History) Store(g Grid) {
	if h.Capacity == 0 {
		return
	}

	if len(h.Grids) >= h.Capacity {
		h.Grids = append([]Grid{CopyGrid(g)}, h.Grids[:h.Capacity-1]...)
	} else {
		h.Grids = append([]Grid{CopyGrid(g)}, h.Grids...)
	}
}

// SetSeed of the game
func (h *History) SetSeed(g Grid) {
	h.Seed = CopyGrid(g)
}
