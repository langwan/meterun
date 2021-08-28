package main

import (
	"github.com/langwan/meterun"
	"time"
)

func main() {
	workers := 2
	works := 2
	meterun.RequestStart(workers, 60 * time.Second)
	meterun.Run(func(workerId int) bool {
		code, err := meterun.Post(workerId,"https://github.com/langwan/meterun", nil,  "")
		if err != nil {
			return false
		}
		if code != 200 {
			return false
		}
		return true
	}, workers, works, 1*time.Second, "test get home")
}
