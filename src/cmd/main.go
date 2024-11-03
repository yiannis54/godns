package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudflare/cloudflare-go"
	"github.com/getsentry/sentry-go"
	_ "github.com/mattn/go-sqlite3"

	"github.com/yiannis54/go-dns/src/internal/config"
	"github.com/yiannis54/go-dns/src/internal/db"
	"github.com/yiannis54/go-dns/src/internal/dns"
	"github.com/yiannis54/go-dns/src/internal/middleware"
)

func main() {
	log.Fatal(run())
}

func run() error {
	cfg, err := config.LoadConfiguration()
	if err != nil {
		return fmt.Errorf("configuration.Init: %w", err)
	}

	gracefulStop, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer stop()

	if err = initSentry(cfg); err != nil {
		return fmt.Errorf("sentry.Init: %w", err)
	}

	db, err := db.InitDB(cfg.DBDsn)
	if err != nil {
		return fmt.Errorf("db.Init: %w", err)
	}

	cloudflareClient, err := cloudflare.NewWithAPIToken(cfg.CloudflareAPIToken)
	if err != nil {
		return fmt.Errorf("cloudflare.Init: %w", err)
	}

	dnsService, err := dns.NewService(cloudflareClient, cfg.Domain, cfg.Interval, db)
	if err != nil {
		return fmt.Errorf("dnsService.Init: %w", err)
	}

	goroutinesEndChan := make(chan struct{})
	updateTicker := initUpdateDNSTicker(dnsService, goroutinesEndChan)

	dnsHandler := dns.NewHandler(dnsService)
	http.HandleFunc("/records/",
		middleware.BasicAuthHandler(
			dnsHandler.GetRecords,
			"admin",
			cfg.Password,
			"Please enter your username and password for this site",
		),
	)

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
	}

	go func() {
		for {
			<-gracefulStop.Done()
			updateTicker.Stop()
			sentry.Flush(3)
			db.SQL.Close()
			server.Shutdown(context.TODO())
		}
	}()

	fmt.Printf("Serving on port :%d...\n", cfg.Port)
	return server.ListenAndServe()
}
