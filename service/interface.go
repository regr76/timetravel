package service

import (
	"context"

	"github.com/regr76/timetravel/entity"
)

// Implements method to get, create, and update record data.
type RecordService interface {

	// GetRecord will retrieve an record (latest version).
	GetRecord(ctx context.Context, id int) (entity.Record, error)

	// CreateRecord will insert a new record (first version).
	//
	// If it a record with that id already exists it will fail (or simply append new version).
	CreateRecord(ctx context.Context, record entity.Record) error

	// UpdateRecord will change the internal `Map` values of the record (from latest version) if it exists.
	// if the update[key] is null it will delete that key from the record's Map.
	//
	// UpdateRecord will error if id <= 0 or the record does not exist with that id.
	UpdateRecord(ctx context.Context, id int, updates map[string]*string) (entity.Record, error)
}

type VersionedRecordService interface {
	RecordService
	GetVersion(ctx context.Context, id int, version int) (entity.Record, error)
	ListRecords(ctx context.Context, id int) (entity.VersionedRecords, error)
	ExportAllRecords(ctx context.Context) ([]string, error)
}

type Storage interface {
	InMemRecords() RecordService
	PersistentRecords() VersionedRecordService
}
