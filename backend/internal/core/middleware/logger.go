package middleware

import (
	"fmt"
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

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		duration := time.Since(start)
		status := rec.status
		size := rec.size

		// Time prefix (HH:mm:ss)
		timeStr := fmt.Sprintf("\033[38;5;245m%8s\033[0m", time.Now().Format("15:04:05"))

		// Protocol version
		versionStr := fmt.Sprintf("HTTP/%d.%d", r.ProtoMajor, r.ProtoMinor)
		version := fmt.Sprintf("\033[1;33m%-8s\033[0m", versionStr)

		// Remote IP
		ip := fmt.Sprintf("\033[38;5;208m%-15s\033[0m", r.RemoteAddr)

		// Status
		statusPadded := fmt.Sprintf("%-3d", status)
		statusStr := colorStatus(statusPadded)

		// Size
		sizeStr := fmt.Sprintf("\033[94m%6dB\033[0m", size)

		// Duration
		durStr := fmt.Sprintf("\033[32m%10s\033[0m", formatDuration(duration))

		// Path and method
		methodColored := fmt.Sprintf("%-6s", colorMethod(r.Method))
		pathColored := fmt.Sprintf("\033[1;33m%-32s\033[0m", r.URL.Path)

		meta := fmt.Sprintf("%s%s",
			colorField("Method", methodColored),
			colorField("Path", pathColored),
		)

		line := fmt.Sprintf(
			"%s %s \033[37mfrom %s\033[0m - %s %s in %s",
			timeStr, version, ip, statusStr, sizeStr, durStr,
		)

		fmt.Println(line + " " + meta)
	})
}

func colorStatus(status string) string {
	switch status[:1] {
	case "2":
		return fmt.Sprintf("\033[32m%s\033[0m", status) // Green for 2xx
	case "3":
		return fmt.Sprintf("\033[36m%s\033[0m", status) // Cyan for 3xx
	case "4":
		return fmt.Sprintf("\033[33m%s\033[0m", status) // Yellow for 4xx
	case "5":
		return fmt.Sprintf("\033[31m%s\033[0m", status) // Red for 5xx
	default:
		return fmt.Sprintf("\033[37m%s\033[0m", status) // White (default)
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
