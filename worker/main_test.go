package main

import (
	"context"
	"io"
	"log/slog"
	"testing"
)

func TestNewLogHandlerUsesConfiguredLevel(t *testing.T) {
	tests := []struct {
		name    string
		devMode bool
		level   slog.Level
		debug   bool
	}{
		{name: "dev info", devMode: true, level: slog.LevelInfo},
		{name: "dev debug", devMode: true, level: slog.LevelDebug, debug: true},
		{name: "json info", level: slog.LevelInfo},
		{name: "json debug", level: slog.LevelDebug, debug: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := newLogHandler(io.Discard, tt.devMode, tt.level)
			if got := handler.Enabled(context.Background(), slog.LevelDebug); got != tt.debug {
				t.Errorf("Enabled(debug) = %t, want %t", got, tt.debug)
			}
			if !handler.Enabled(context.Background(), slog.LevelInfo) {
				t.Error("Enabled(info) = false, want true")
			}
		})
	}
}
