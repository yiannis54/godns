package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type EnvConfig struct {
	Env       string
	Port      int
	DBDsn     string
	SentryDsn string
	Password  string
	Interval  time.Duration
	Domain    string

	CloudflareAPIToken string
}

func LoadConfiguration() (*EnvConfig, error) {
	port, err1 := strconv.Atoi(os.Getenv("SRV_PORT"))
	interval, err2 := strconv.Atoi(os.Getenv("UPDATE_INTERVAL"))
	if errs := errors.Join(err1, err2); errs != nil {
		return nil, errs
	}

	return &EnvConfig{
		Env:                os.Getenv("ENV"),
		Port:               port,
		SentryDsn:          os.Getenv("SENTRY_DSN"),
		DBDsn:              os.Getenv("DB_DSN"),
		Password:           os.Getenv("PASSWORD"),
		Interval:           time.Duration(interval) * time.Hour,
		Domain:             os.Getenv("DOMAIN"),
		CloudflareAPIToken: os.Getenv("CLOUDFLARE_API_TOKEN"),
	}, nil
}
