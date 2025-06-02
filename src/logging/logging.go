package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

const (
	reset       = "\033[0m"
	black       = "\033[30m"
	red         = "\033[31m"
	green       = "\033[32m"
	yellow      = "\033[33m"
	blue        = "\033[34m"
	magenta     = "\033[35m"
	cyan        = "\033[36m"
	white       = "\033[37m"
	boldBlack   = "\033[1;30m"
	boldRed     = "\033[1;31m"
	boldGreen   = "\033[1;32m"
	boldYellow  = "\033[1;33m"
	boldBlue    = "\033[1;34m"
	boldMagenta = "\033[1;35m"
	boldCyan    = "\033[1;36m"
	boldWhite   = "\033[1;37m"
)

// Setup configures the global logger based on environment
func SetupLogger() {
	// Pretty text format for development
	handler := NewPrettyHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	// Set as default logger
	slog.SetDefault(slog.New(handler))
}

// PrettyHandler formats logs in a more human-readable format for development
type PrettyHandler struct {
	h     slog.Handler
	w     io.Writer
	attrs []slog.Attr
}

// NewPrettyHandler creates a new PrettyHandler
func NewPrettyHandler(w io.Writer, opts *slog.HandlerOptions) *PrettyHandler {
	return &PrettyHandler{
		h: slog.NewTextHandler(w, opts),
		w: w,
	}
}

func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		h:     h.h.WithAttrs(attrs),
		w:     h.w,
		attrs: append(h.attrs, attrs...),
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return &PrettyHandler{
		h: h.h.WithGroup(name),
		w: h.w,
	}
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	timeStr := r.Time.Format("15:04:05.000")
	levelStr := r.Level.String()[:4] // Get first 4 chars: INFO, WARN, etc.

	// Color codes for different log levels
	var colorCode string
	switch r.Level {
	case slog.LevelDebug:
		colorCode = "\033[36m" // Cyan
	case slog.LevelInfo:
		colorCode = "\033[32m" // Green
	case slog.LevelWarn:
		colorCode = "\033[33m" // Yellow
	case slog.LevelError:
		colorCode = "\033[31m" // Red
	default:
		colorCode = "\033[0m" // Default/Reset
	}
	resetColor := "\033[0m"

	// Format: [TIME] [COLORED_LEVEL] MESSAGE key=value key=value
	_, err := fmt.Fprintf(h.w, "[%s] [%s%s%s] %s",
		timeStr,
		colorCode, levelStr, resetColor,
		r.Message)
	if err != nil {
		return err
	}

	// Add attributes if there are any
	hasAttrs := len(h.attrs) > 0 || r.NumAttrs() > 0
	if hasAttrs {
		if _, err := fmt.Fprint(h.w, " "); err != nil {
			return err
		}
	}

	// Print attributes from the handler
	for _, attr := range h.attrs {
		_, err := fmt.Fprintf(h.w, "%s=%v ", attr.Key, attr.Value.Any())
		if err != nil {
			return err
		}
	}

	// Print attributes from the record
	r.Attrs(func(attr slog.Attr) bool {
		_, err := fmt.Fprintf(h.w, "%s=%v ", attr.Key, attr.Value.Any())
		// We can't return the error inside this callback, so we'll just continue
		// if there's an error and handle the incomplete output
		return err == nil
	})

	// Add newline at the end
	_, err = fmt.Fprintln(h.w)
	return err
}

// Info logs formatted info messages in Printf style
func Info(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	slog.Info(message)
}

// Error logs structured error information
func Error(msg string, err error, keyValues ...any) {
	// Combine the error with any additional key-values
	allAttrs := make([]any, 0, len(keyValues)+2)
	allAttrs = append(allAttrs, "error", err)
	allAttrs = append(allAttrs, keyValues...)

	slog.Error(msg, allAttrs...)
}

// Fatal logs an error and exits the program
func Fatal(msg string, err error, keyValues ...any) {
	formattedMsg := fmt.Sprintf("%s[FATAL]%s %s", boldRed, reset, msg)

	// Combine the error with any additional key-values
	allAttrs := make([]any, 0, len(keyValues)+2)
	allAttrs = append(allAttrs, "error", err)
	allAttrs = append(allAttrs, keyValues...)

	slog.Error(formattedMsg, allAttrs...)
	os.Exit(1)
}

// Warn logs a warning message
func Warn(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	slog.Warn(message)
}

// InfoWithFields logs formatted info messages with optional key-value pairs
func InfoWithFields(format string, args []any, keyValues ...any) {
	message := fmt.Sprintf(format, args...)

	if len(keyValues) > 0 {
		slog.Info(message, keyValues...)
	} else {
		slog.Info(message)
	}
}
