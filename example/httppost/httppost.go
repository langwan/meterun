package main

import (
	"github.com/langwan/meterun"
	"time"
)

func main() {
	meterun.Run(func() bool {
		code, err := meterun.Post("https://github.com/langwan/meterun", nil, "", "application/json", 60*time.Second)
		if err != nil {
			return false
		}
		if code != 200 {
			return false
		}
		return true
	}, 2, 2, 1*time.Second, "test login http post")
}
