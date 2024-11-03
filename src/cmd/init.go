package main

import (
	"github.com/getsentry/sentry-go"

	"github.com/yiannis54/go-dns/src/internal/config"
)

func initSentry(cfg *config.EnvConfig) error {
	return sentry.Init(sentry.ClientOptions{
		Dsn:              cfg.SentryDsn,
		Release:          "0.0.1",
		Debug:            cfg.Env == "dev",
		AttachStacktrace: cfg.Env != "dev",
		EnableTracing:    cfg.Env != "dev",
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate:   0.4,
		ProfilesSampleRate: 0.3,
	})
}
