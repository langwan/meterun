# meterun

golang stress test tool http post or custom function, support the log file

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

```

output

```
- START - 2021-08-27 15:35 - test login http post
|       SEC |                 QPS |            MAX TIME |            MIN TIME |                 P90 |                 BAD |
|         1 |                   2 |        2.029433083s |        2.028006333s |        2.029433083s |                   0 |

|       SEC |                 QPS |            MAX TIME |            MIN TIME |                 P90 |                 BAD |
|         2 |                   2 |        655.149959ms |        649.918625ms |        655.149959ms |                   0 |

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
	meterun.Run(func() bool {
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

2. post

* `url` : url
* `header`: net http Header struct
* `contentType`: content type
* `timeout`: http Timeout