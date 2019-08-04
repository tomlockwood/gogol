package gol

import "sync"

// TickFunction is called to Tick a Game
type TickFunction func(g Game, gameNumber int)

// Run a Game
func Run(g Game, TickFunction TickFunction, gameNumber int) {
	TickFunction(g, gameNumber)
}

// RunMany games of life concurrently
// TickFunction is run on every tick of the game, so it
// can be used to halt execution early or change the state
func RunMany(Options Options, gameAmount int, TickFunction TickFunction) {
	var wg sync.WaitGroup
	wg.Add(gameAmount)
	for i := 0; i < gameAmount; i++ {
		go func(i int) {
			defer wg.Done()
			g := MakeGame(Options)
			Run(g, TickFunction, i)
		}(i)
	}
	wg.Wait()
}
