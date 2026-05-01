package logging

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// FieldLogger is a slog-backed logger that accepts logging.Field values.
type FieldLogger interface {
	DebugContext(ctx context.Context, msg string, fields ...Field)
	InfoContext(ctx context.Context, msg string, fields ...Field)
	WarnContext(ctx context.Context, msg string, fields ...Field)
	ErrorContext(ctx context.Context, msg string, fields ...Field)
	WithFields(fields ...Field) FieldLogger
	Underlying() *slog.Logger
}

// SlogOptions configures NewSlogLogger.
type SlogOptions struct {
	Level     slog.Leveler
	Output    io.Writer
	AddSource bool
}

type fieldSlogLogger struct{ l *slog.Logger }

// NewSlogLogger creates a slog-backed FieldLogger.
func NewSlogLogger(opts SlogOptions) FieldLogger {
	if opts.Level == nil {
		opts.Level = slog.LevelInfo
	}
	if opts.Output == nil {
		opts.Output = os.Stdout
	}
	h := slog.NewJSONHandler(opts.Output, &slog.HandlerOptions{AddSource: opts.AddSource, Level: opts.Level})
	return fieldSlogLogger{l: slog.New(h)}
}

// FromSlogFields wraps an existing slog.Logger as a FieldLogger.
func FromSlogFields(l *slog.Logger) FieldLogger {
	if l == nil {
		return NewSlogLogger(SlogOptions{})
	}
	return fieldSlogLogger{l: l}
}

func (s fieldSlogLogger) DebugContext(ctx context.Context, msg string, fields ...Field) {
	s.l.DebugContext(ctx, msg, attrs(fields)...)
}
func (s fieldSlogLogger) InfoContext(ctx context.Context, msg string, fields ...Field) {
	s.l.InfoContext(ctx, msg, attrs(fields)...)
}
func (s fieldSlogLogger) WarnContext(ctx context.Context, msg string, fields ...Field) {
	s.l.WarnContext(ctx, msg, attrs(fields)...)
}
func (s fieldSlogLogger) ErrorContext(ctx context.Context, msg string, fields ...Field) {
	s.l.ErrorContext(ctx, msg, attrs(fields)...)
}
func (s fieldSlogLogger) WithFields(fields ...Field) FieldLogger {
	return fieldSlogLogger{l: s.l.With(attrs(fields)...)}
}
func (s fieldSlogLogger) Underlying() *slog.Logger { return s.l }

func attrs(fields []Field) []any {
	out := make([]any, 0, len(fields))
	for _, f := range fields {
		out = append(out, slog.Any(f.Key, f.Value))
	}
	return out
}

func String(key, value string) Field      { return F(key, value) }
func Int(key string, value int) Field     { return F(key, value) }
func Int64(key string, value int64) Field { return F(key, value) }
func Bool(key string, value bool) Field   { return F(key, value) }
func Any(key string, value any) Field     { return F(key, value) }
func Err(err error) Field                 { return F("error", err) }
