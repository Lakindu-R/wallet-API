package middleware

import (
	"net/http"
	"sync"
	"time"

	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

type client struct {
	count     int
	timestamp time.Time
}

var clients = make(map[string]*client)
var mu sync.Mutex

func getLimit() int {
	val := os.Getenv("RATE_LIMIT")
	if val == "" {
		return 5
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return 5
	}
	return n
}

var requestLimit = getLimit()
var window = time.Minute

func RateLimiter(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		ip := c.RealIP()

		mu.Lock()
		defer mu.Unlock()

		cl, exists := clients[ip]

		if !exists {
			clients[ip] = &client{count: 1, timestamp: time.Now()}
		} else {

			// reset after time window
			if time.Since(cl.timestamp) > window {
				cl.count = 1
				cl.timestamp = time.Now()
			} else {
				cl.count++

				if cl.count > requestLimit {
					return c.JSON(http.StatusTooManyRequests, map[string]string{
						"error": "Exceeded Rate Limit.",
					})
				}
			}
		}

		return next(c)
	}
}
