package middleware

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

const (
	cReset  = "\033[0m"
	cGray   = "\033[90m"
	cYellow = "\033[1;33m"
	cOrange = "\033[38;5;208m"
	cBlue   = "\033[94m"
	cGreen  = "\033[32m"
	cWhite  = "\033[37m"
)

func CustomRequestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:       true,
		LogMethod:       true,
		LogURI:          true,
		LogRemoteIP:     true,
		LogLatency:      true,
		LogResponseSize: true,
		LogProtocol:     true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			// 1. Base Info
			timeStr := fmt.Sprintf("%s%s%s", cGray, time.Now().Format("15:04:05"), cReset)
			proto := fmt.Sprintf("%s%-10s%s", cYellow, v.Protocol, cReset)
			ip := fmt.Sprintf("%sfrom %s%-15s%s", cGray, cOrange, v.RemoteIP, cReset)

			// 2. Create Segments (Key and Value joined with no internal padding)
			sizeStr := fmt.Sprintf("%sSize: %s%dB%s", cGray, cBlue, v.ResponseSize, cReset)
			durStr := fmt.Sprintf("%sDuration: %s%s%s", cGray, cGreen, formatDuration(v.Latency), cReset)
			statusStr := fmt.Sprintf("%sStatus: %s", cGray, colorStatus(fmt.Sprintf("%d", v.Status)))
			methStr := fmt.Sprintf("%sMethod: %s%s%s", cGray, cWhite, v.Method, cReset)
			pathStr := fmt.Sprintf("%sPath: %s%s%s", cGray, cYellow, v.URI, cReset)

			// 3. Print with fixed column widths for the segments
			// We use %-N[type] to pad the entire colored string block.
			// Note: Because ANSI codes take up "length" but are invisible,
			// we add the escape sequence length to the padding value.
			fmt.Printf("%s  %s  %s    %-25s %-30s %-22s %-22s %s\n",
				timeStr,
				proto,
				ip,
				sizeStr,
				durStr,
				statusStr,
				methStr,
				pathStr,
			)

			return nil
		},
	})
}

func colorStatus(status string) string {
	color := cWhite
	switch status[0] {
	case '2':
		color = cGreen
	case '3':
		color = "\033[36m"
	case '4':
		color = cYellow
	case '5':
		color = cRed
	}
	return fmt.Sprintf("%s%s%s", color, status, cReset)
}

const cRed = "\033[31m"

func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%.2fµs", float64(d)/float64(time.Microsecond))
	}
	return fmt.Sprintf("%.2fms", float64(d)/float64(time.Millisecond))
}
