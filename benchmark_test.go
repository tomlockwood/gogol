package gol

import (
	"testing"
)

func tick(g Game, gameNumber int) {
	for i := 0; i <= 1500; i++ {
		g.Tick()
	}
}

func BenchmarkRunManyConways(b *testing.B) {

	var r0 = Rule{false, [9]uint8{0, 0, 0, 1, 0, 0, 0, 0, 0}, Colour{}}
	var r1 = Rule{true, [9]uint8{0, 0, 1, 1, 0, 0, 0, 0, 0}, Colour{}}
	var rs = Rules{[]Rule{r0, r1}}

	opts := Options{
		0,
		0,
		Grid{},
		2,
		rs}

	RunMany(opts, 1000, tick)
}
