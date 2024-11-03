package dns

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrDBNil = errors.New("DB is nil")

type repository struct {
	db *sql.DB
}

func newRepository(db *sql.DB) (*repository, error) {
	if db == nil {
		return nil, ErrDBNil
	}
	return &repository{
		db: db,
	}, nil
}

func (r *repository) validate() error {
	if r.db == nil {
		return ErrDBNil
	}
	return nil
}

func (r *repository) saveUpdateRecord(newIP string) error {
	if err := r.validate(); err != nil {
		return err
	}

	dbCtx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	query := `INSERT INTO dns_updates (new_ip, created_at) values ($1, datetime('now'))`
	_, err := r.db.ExecContext(dbCtx, query, newIP)

	return err
}

func (r *repository) getDNSRecordUpdates() (UpdateRecordList, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	dbCtx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()

	query := `SELECT * FROM dns_updates order by created_at DESC`
	rows, err := r.db.QueryContext(dbCtx, query)
	if err != nil {
		return nil, err
	}

	recordList := make(UpdateRecordList, 0)
	defer rows.Close()
	for rows.Next() {
		var id uint64
		var ip string
		var createdAt time.Time

		err = rows.Scan(&id, &ip, &createdAt)
		if err != nil {
			return nil, err
		}

		recordList = append(recordList, UpdateRecord{
			ID:        id,
			NewIP:     ip,
			CreatedAt: createdAt,
		})
	}

	return recordList, rows.Err()
}
