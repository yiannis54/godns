package main

import (
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/yiannis54/go-dns/src/internal/dns"
)

func initUpdateDNSTicker(dnsService *dns.Service, goroutinesEndChan chan struct{}) *time.Ticker {
	ticker := time.NewTicker(dnsService.Interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				if err := dnsService.UpdateRecord(); err != nil {
					sentry.CaptureException(err)
				}
			case <-goroutinesEndChan:
				ticker.Stop()
				return
			}
		}
	}()
	return ticker
}
