package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorWhite  = "\033[97m"
)

// ColorHandler is a custom handler that adds colors to log output
type ColorHandler struct {
	opts   slog.HandlerOptions
	out    io.Writer
	mu     *sync.Mutex
	groups []string
	attrs  []slog.Attr
}

// NewColorHandler creates a new ColorHandler
func NewColorHandler(out io.Writer, opts *slog.HandlerOptions) *ColorHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &ColorHandler{
		out:  out,
		opts: *opts,
		mu:   &sync.Mutex{},
	}
}

// Enabled reports whether the handler handles records at the given level
func (h *ColorHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

// Handle handles the Record
//
//nolint:gocritic // slog.Handler interface requires value receiver
func (h *ColorHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Format time
	timeStr := r.Time.Format("15:04:05")

	// Get level color and text
	levelColor, levelText := h.levelColorAndText(r.Level)

	// Write time and level
	if _, err := fmt.Fprintf(h.out, "%s%s%s %s%-5s%s ",
		colorGray, timeStr, colorReset,
		levelColor, levelText, colorReset); err != nil {
		return err
	}

	// Write message
	if _, err := fmt.Fprintf(h.out, "%s", r.Message); err != nil {
		return err
	}

	// Write attributes
	if len(h.attrs) > 0 || r.NumAttrs() > 0 {
		if _, err := fmt.Fprintf(h.out, " %s", colorGray); err != nil {
			return err
		}

		// Write handler's attributes
		for _, attr := range h.attrs {
			if err := h.writeAttr(attr); err != nil {
				return err
			}
		}

		// Write record's attributes
		var attrErr error
		r.Attrs(func(attr slog.Attr) bool {
			if err := h.writeAttr(attr); err != nil {
				attrErr = err
				return false
			}
			return true
		})
		if attrErr != nil {
			return attrErr
		}

		if _, err := fmt.Fprintf(h.out, "%s", colorReset); err != nil {
			return err
		}
	}

	_, err := fmt.Fprintln(h.out)
	return err
}

// WithAttrs returns a new Handler with additional attributes
func (h *ColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h2 := *h
	h2.attrs = append(h2.attrs, attrs...)
	return &h2
}

// WithGroup returns a new Handler with the given group name
func (h *ColorHandler) WithGroup(name string) slog.Handler {
	h2 := *h
	h2.groups = append(h2.groups, name)
	return &h2
}

func (h *ColorHandler) levelColorAndText(level slog.Level) (color, text string) {
	switch level {
	case slog.LevelDebug:
		return colorCyan, "DEBUG"
	case slog.LevelInfo:
		return colorGreen, "INFO"
	case slog.LevelWarn:
		return colorYellow, "WARN"
	case slog.LevelError:
		return colorRed, "ERROR"
	default:
		return colorWhite, level.String()
	}
}

func (h *ColorHandler) writeAttr(attr slog.Attr) error {
	// Special handling for certain attribute keys
	valueColor := colorWhite
	if attr.Key == "error" {
		valueColor = colorRed
	}

	// Format based on value type
	switch v := attr.Value.Any().(type) {
	case string:
		_, err := fmt.Fprintf(h.out, " %s=%s%q%s", attr.Key, valueColor, v, colorGray)
		return err
	case time.Time:
		_, err := fmt.Fprintf(h.out, " %s=%s%s%s", attr.Key, valueColor, v.Format(time.RFC3339), colorGray)
		return err
	case error:
		_, err := fmt.Fprintf(h.out, " %s=%s%q%s", attr.Key, colorRed, v.Error(), colorGray)
		return err
	default:
		_, err := fmt.Fprintf(h.out, " %s=%s%v%s", attr.Key, valueColor, v, colorGray)
		return err
	}
}
