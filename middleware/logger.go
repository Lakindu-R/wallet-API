package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("Request:", c.Request().Method, c.Path())
		return next(c)
	}
}
