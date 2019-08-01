package gol

import (
	"testing"
)

func tick(g Game, gameNumber int) {
	for i := 0; i <= 1500; i++ {
		g.Tick()
	}
}

var cr0 = Rule{false, [9]uint8{0, 0, 0, 1, 0, 0, 0, 0, 0}, Colour{}}
var cr1 = Rule{true, [9]uint8{0, 0, 1, 1, 0, 0, 0, 0, 0}, Colour{}}
var crs = Rules{[]Rule{r0, r1}}

var conwayOpts = Options{
	0,
	0,
	Grid{},
	2,
	crs}

func BenchRunMany(opts Options, gameAmount int, TickFunction TickFunction, b *testing.B) {
	RunMany(opts, gameAmount, TickFunction)
}

func Benchmark1000Conways(b *testing.B) {
	RunMany(conwayOpts, 1000, tick)
}

func Benchmark100Conways(b *testing.B) {
	RunMany(conwayOpts, 100, tick)
}
