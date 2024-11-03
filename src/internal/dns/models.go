package dns

import "time"

type UpdateRecord struct {
	ID        uint64
	NewIP     string
	CreatedAt time.Time
}

type UpdateRecordList []UpdateRecord
