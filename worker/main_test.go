package main

import (
	"io"
	"testing"

	"github.com/rs/zerolog"
)

func TestNewLoggerUsesConfiguredLevel(t *testing.T) {
	tests := []struct {
		name    string
		devMode bool
		level   zerolog.Level
		debug   bool
	}{
		{name: "dev info", devMode: true, level: zerolog.InfoLevel},
		{name: "dev debug", devMode: true, level: zerolog.DebugLevel, debug: true},
		{name: "json info", level: zerolog.InfoLevel},
		{name: "json debug", level: zerolog.DebugLevel, debug: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := newLogger(io.Discard, tt.devMode, tt.level)
			if got := logger.GetLevel() <= zerolog.DebugLevel; got != tt.debug {
				t.Errorf("Enabled(debug) = %t, want %t", got, tt.debug)
			}
			if logger.GetLevel() > zerolog.InfoLevel {
				t.Error("Enabled(info) = false, want true")
			}
		})
	}
}
