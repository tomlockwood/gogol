package gol

import (
	"math/rand"
	"sync"
	"time"
)

// goroutine safe random integer generator
var randMutex sync.Mutex
var src = rand.NewSource(time.Now().UnixNano())
var r = rand.New(src)

func randInt(i int) int {
	randMutex.Lock()
	integer := r.Intn(i)
	randMutex.Unlock()
	return integer
}
