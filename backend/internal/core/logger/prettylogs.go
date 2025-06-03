package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

type Handler struct {
	writer    io.Writer
	attrs     []slog.Attr
	groupName string
	level     slog.Level
}

func NewHandler(w io.Writer, level slog.Level) slog.Handler {
	return &Handler{
		writer: w,
		level:  level,
	}
}

func (h *Handler) Enabled(_ context.Context, lvl slog.Level) bool {
	return lvl >= h.level
}

func (h *Handler) Handle(_ context.Context, r slog.Record) error {
	b := &strings.Builder{}

	// Time
	fmt.Fprintf(b, "\033[90m%s\033[0m ", r.Time.Format("15:04:05"))

	// Level with color
	levelColor := map[slog.Level]string{
		slog.LevelDebug: "\033[36m",
		slog.LevelInfo:  "\033[32m",
		slog.LevelWarn:  "\033[33m",
		slog.LevelError: "\033[31m",
	}
	levelStr := strings.ToUpper(r.Level.String())
	fmt.Fprintf(b, "%s%-5s\033[0m ", levelColor[r.Level], levelStr)

	// Message
	fmt.Fprintf(b, "%s ", r.Message)

	// Attrs from log call
	r.Attrs(func(a slog.Attr) bool {
		if a.Value.Kind() == slog.KindGroup {
			fmt.Fprintf(b, "\n  \033[90m%s:\033[0m {\n", a.Key)
			for _, nested := range a.Value.Group() {
				fmt.Fprintf(b, "    \033[90m%s:\033[0m %v\n", nested.Key, nested.Value)
			}
			fmt.Fprint(b, "  }\n")
		} else {
			fmt.Fprintf(b, "\033[90m%s=\033[0m%v ", a.Key, a.Value)
		}
		return true
	})

	// Global attrs (e.g. logger.With(...))
	for _, a := range h.attrs {
		fmt.Fprintf(b, "\033[90m%s=\033[0m%s ", a.Key, a.Value.String())
	}

	fmt.Fprint(h.writer, b.String()+"\n")
	return nil
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{
		writer:    h.writer,
		level:     h.level,
		attrs:     append(h.attrs, attrs...),
		groupName: h.groupName,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		writer:    h.writer,
		level:     h.level,
		attrs:     h.attrs,
		groupName: name,
	}
}
