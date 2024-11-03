package dns

import (
	"context"
	"errors"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/yiannis54/go-dns/src/internal/db"
)

type Service struct {
	domain     string
	Interval   time.Duration
	client     *cloudflare.API
	repository *repository
}

func NewService(client *cloudflare.API, domain string, interval time.Duration, db *db.DB) (*Service, error) {
	repo, err := newRepository(db.SQL)
	if err != nil {
		return nil, err
	}

	return &Service{
		client:     client,
		domain:     domain,
		Interval:   interval,
		repository: repo,
	}, nil
}

func (s *Service) getRecordID() (*cloudflare.DNSRecord, string, error) {
	zoneID, err := s.client.ZoneIDByName(s.domain)
	if err != nil {
		return nil, "", err
	}

	records, _, err := s.client.ListDNSRecords(
		context.Background(),
		cloudflare.ZoneIdentifier(zoneID),
		cloudflare.ListDNSRecordsParams{
			Type: "A",
		},
	)
	if err != nil {
		return nil, "", err
	}

	if len(records) > 1 {
		return nil, "", errors.New("more than one A records exist")
	}

	return &records[0], zoneID, nil
}

func (s *Service) UpdateRecord() error {
	newIP, err1 := getIP()
	record, zoneID, err2 := s.getRecordID()
	if errs := errors.Join(err1, err2); errs != nil {
		return errs
	}

	if newIP == record.Content {
		// IP did not change, no need to update.
		return nil
	}

	params := cloudflare.UpdateDNSRecordParams{
		ID:      record.ID,
		Content: newIP,
	}

	_, err := s.client.UpdateDNSRecord(
		context.Background(),
		cloudflare.ZoneIdentifier(zoneID),
		params,
	)
	if err != nil {
		return err
	}

	return s.recordUpdate(newIP)
}

func (s *Service) recordUpdate(newIP string) error {
	return s.repository.saveUpdateRecord(newIP)
}
