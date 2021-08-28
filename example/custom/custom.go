package main

import (
	"github.com/langwan/meterun"
	"math/rand"
	"time"
)

func main() {
	meterun.Run(func(workerId int) bool {
		ri := rand.Intn(1000)
		time.Sleep(time.Duration(ri) * time.Millisecond)
		return true
	}, 2, 2, 1*time.Second, "test custom func")
}
