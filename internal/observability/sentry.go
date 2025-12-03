package observability

import (
	"time"

	"github.com/getsentry/sentry-go"
)

// InitSentry configures sentry when DSN provided.
func InitSentry(dsn string) error {
	if dsn == "" {
		return nil
	}
	return sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		TracesSampleRate: 0.1,
	})
}

// Flush ensures buffered events delivered.
func Flush(timeout time.Duration) {
	sentry.Flush(timeout)
}

// CaptureError reports non-nil error to sentry.
func CaptureError(err error) {
	if err == nil {
		return
	}
	sentry.CaptureException(err)
}
