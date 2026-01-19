package main

import (
	mw "github.com/iankencruz/threefive/internal/middleware"
	"github.com/labstack/echo/v5"
)

// ANSI Color Codes
const (
	cReset  = "\033[0m"
	cGray   = "\033[38;5;245m"
	cYellow = "\033[1;33m"
	cOrange = "\033[38;5;208m"
	cBlue   = "\033[94m"
	cGreen  = "\033[32m"
	cWhite  = "\033[37m"
)

func main() {

	e := echo.New()

	e.Use(mw.CustomRequestLogger())

	e.GET("/", func(c *echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("Failed to start server", "Error", err)
	}

}
