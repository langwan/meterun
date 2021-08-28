# meterun

golang benchmarking tool http post or custom function, support auto save the log file

## Installation

```sh
$ go get -u github.com/langwan/meterun
```

2. Import it in your code:

```go
import "github.com/langwan/meterun"
```

## Examples

1. test http post

```go
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
```

output

```
- START - 2021-08-29 04:00 - test custom func
|       SEC |                 QPS |            MAX TIME |            MIN TIME |                 P90 |                 BAD |
|         1 |                   2 |           892.044ms |         84.835542ms |           892.044ms |                   0 |

|       SEC |                 QPS |            MAX TIME |            MIN TIME |                 P90 |                 BAD |
|         2 |                   2 |         852.04375ms |         62.592542ms |         852.04375ms |                   0 |

- END -
| REQ TOTAL |             REQ BAS |             WORKERS |               WORKS |               SLEEP |                     |
|         4 |                   0 |                   2 |                   2 |                  1s |                     |
```

gen log file is

```
test login http post_2021-08-27 15:34.txt
```

2. custom function

```go
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
```

## params

1. meterun.Run

* `f` : function executed
* `workers`: the number of goroutine
* `works`: the number of executions per goroutine
* `sleep`: each execution of the rest time
* `title`: report title

2. post or get

* `workerId` current worker id 
* `url` : url
* `header`: net http header struct
* `content`: content