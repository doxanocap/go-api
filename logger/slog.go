package logger

import (
	"auth-api/config"
	"context"
	"fmt"
	"log/slog"
	"os"
)

const (
	keyPayload = "payload"
	keyModule  = "module"
)

type Logger struct {
	log    *slog.Logger
	attrs  []slog.Attr
	module string
}

func InitSlogLogger(cfg *config.Config) *Logger {
	var handler slog.Handler
	switch cfg.ENV {
	case config.EnvProduction:
		handler = slog.NewJSONHandler(os.Stdout, nil)
	default:
		handler = slog.NewTextHandler(os.Stdout, nil)
	}

	return &Logger{
		log: slog.New(handler),
	}
}

func (l *Logger) WithModule(module string) *Logger {
	return &Logger{
		log:    l.log,
		module: fmt.Sprintf("%s[%s]", l.module, module),
	}
}

func (l *Logger) WithAttrs(attrs ...slog.Attr) *Logger {
	return &Logger{
		attrs:  attrs,
		log:    l.log,
		module: l.module,
	}
}

func (l *Logger) Log(level slog.Level, msg string, attrs ...slog.Attr) {
	l.attrs = append(l.attrs, attrs...)
	l.log.LogAttrs(context.Background(), level, msg,
		[]slog.Attr{
			slog.String(keyModule, l.module),
			{
				Key:   keyPayload,
				Value: slog.GroupValue(l.attrs...),
			},
		}...,
	)
}

func (l *Logger) Info(msg string, args ...slog.Attr) {
	l.Log(slog.LevelInfo, msg, args...)
}

func (l *Logger) Debug(msg string, args ...slog.Attr) {
	l.Log(slog.LevelDebug, msg, args...)
}

func (l *Logger) Warn(msg string, args ...slog.Attr) {
	l.Log(slog.LevelWarn, msg, args...)
}

func (l *Logger) Error(msg string, args ...slog.Attr) {
	l.Log(slog.LevelError, msg, args...)
}
