package api

import (
	"fmt"
	"log/slog"
	"net/url"

	"github.com/imroc/req/v3"
)

type ServerConfig struct {
	URL    *url.URL
	DryRun bool
	APIKey string
}

type ErrorResponse struct {
	Success string `json:"success"`
	Message string `json:"error"`
}

func (err ErrorResponse) Error() string {
	return err.Message
}

type SlogAdapter struct {
}

// Errorf logs error messages using slog.Error.
func (s *SlogAdapter) Errorf(format string, v ...any) {
	slog.Error(fmt.Sprintf(format, v...))
}

// Warnf logs warning messages using slog.Warn.
func (s *SlogAdapter) Warnf(format string, v ...any) {
	slog.Warn(fmt.Sprintf(format, v...))
}

// Debugf logs debug messages using slog.Debug.
func (s *SlogAdapter) Debugf(format string, v ...any) {
	slog.Debug(fmt.Sprintf(format, v...))
}

func NewClient(config ServerConfig) *req.Client {
	return req.C().
		SetCommonHeader("x-api-key", config.APIKey).
		SetBaseURL(config.URL.String()).
		SetLogger(&SlogAdapter{}).
		EnableDebugLog()
}
