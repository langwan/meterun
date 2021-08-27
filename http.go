package meterun

import (
	netHttp "net/http"
	"strings"
	"time"
)

func Post(url string, header netHttp.Header, content string, contentType string, timeout time.Duration) (int, error) {
	c := netHttp.Client{Timeout: timeout}

	req, err := netHttp.NewRequest("POST", url, strings.NewReader(content))
	req.Header = header

	do, err := c.Do(req)
	if err != nil {
		return -1, err
	}
	return do.StatusCode, nil
}
