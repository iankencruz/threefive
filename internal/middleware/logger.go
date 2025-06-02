package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *statusRecorder) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *statusRecorder) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(rec, r)

			duration := time.Since(start)
			status := rec.status
			size := rec.size

			version := fmt.Sprintf("\033[1;33mHTTP/%d.%d\033[0m", r.ProtoMajor, r.ProtoMinor)
			ip := fmt.Sprintf("\033[38;5;208m%s\033[0m", r.RemoteAddr)
			statusStr := colorStatus(status)
			sizeStr := fmt.Sprintf("\033[94m%4dB\033[0m", size)
			durStr := fmt.Sprintf("\033[32m%8s\033[0m", formatDuration(duration))

			line := fmt.Sprintf(
				"%s \033[37mfrom %s\033[0m - %s %s in %s",
				version, ip, statusStr, sizeStr, durStr,
			)

			pathColored := fmt.Sprintf("\033[1;33m%s\033[0m", r.URL.Path) // bold yellow
			meta := fmt.Sprintf("%s%s",
				colorField("Method", padRight(colorMethod(r.Method), 6)),
				colorField("Path", padRight(pathColored, 32)),
			)

			logger.Info(line + " " + meta)
		})
	}
}

func colorStatus(status int) string {
	switch {
	case status >= 200 && status < 300:
		return fmt.Sprintf("\033[32m%3d\033[0m", status) // green
	case status >= 300 && status < 400:
		return fmt.Sprintf("\033[36m%3d\033[0m", status) // cyan
	case status >= 400 && status < 500:
		return fmt.Sprintf("\033[33m%3d\033[0m", status) // yellow
	default:
		return fmt.Sprintf("\033[31m%3d\033[0m", status) // red
	}
}

func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%.2fµs", float64(d)/float64(time.Microsecond))
	}
	return fmt.Sprintf("%.2fms", float64(d)/float64(time.Millisecond))
}

func colorMethod(method string) string {
	switch method {
	case "GET":
		return "\033[1;32mGET\033[0m"
	case "POST":
		return "\033[1;36mPOST\033[0m"
	case "PUT":
		return "\033[1;33mPUT\033[0m"
	case "PATCH":
		return "\033[1;35mPATCH\033[0m"
	case "DELETE":
		return "\033[1;31mDELETE\033[0m"
	default:
		return fmt.Sprintf("\033[1m%s\033[0m", method)
	}
}

func colorField(key, value string) string {
	return fmt.Sprintf("\033[37m%s: \033[0m%s ", key, value)
}

func padRight(s string, length int) string {
	visibleLen := visibleLength(s)
	if visibleLen < length {
		return s + spaces(length-visibleLen)
	}
	return s
}

func spaces(n int) string {
	return fmt.Sprintf("%*s", n, "")
}

func visibleLength(s string) int {
	// Strip ANSI escape codes to calculate visible width
	clean := ""
	skip := false
	for i := range len(s) {
		if s[i] == '\033' {
			skip = true
			continue
		}
		if skip && s[i] == 'm' {
			skip = false
			continue
		}
		if !skip {
			clean += string(s[i])
		}
	}
	return len(clean)
}
