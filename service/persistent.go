package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/regr76/timetravel/entity"
)

// PersistentRecordService is an in-memory implementation of RecordService.
type PersistentRecordService struct {
	data map[int][]entity.PersistentRecord
	db   *sql.DB
}

func NewPersistentRecordService(db *sql.DB) PersistentRecordService {
	return PersistentRecordService{
		data: map[int][]entity.PersistentRecord{},
		db:   db,
	}
}

// GetRecord will retrieve record with latest version.
func (s *PersistentRecordService) GetRecord(ctx context.Context, id int) (entity.Record, error) {
	records := s.data[id]
	if len(records) == 0 {
		return nil, ErrRecordDoesNotExist
	}

	record := records[len(records)-1] // get the latest version of the record

	copied, ok := record.Copy().(*entity.PersistentRecord) // copy is necessary so modifations to the record don't change the stored record
	if !ok {
		return nil, errors.New("failed to cast record to PersistentRecord")
	}
	return copied, nil
}

// ListRecords will retrieve record containing all versions.
func (s *PersistentRecordService) ListRecords(ctx context.Context, id int) (entity.VersionedRecords, error) {
	records := s.data[id]
	if len(records) == 0 {
		return nil, ErrRecordDoesNotExist
	}

	output := &entity.PersistentRecords{
		Records: records,
	}

	return output.Copy(), nil
}

// CreateRecord will create record with version 1.
func (s *PersistentRecordService) CreateRecord(ctx context.Context, record entity.Record) error {
	id := record.GetID()
	if id <= 0 {
		return ErrRecordIDInvalid
	}

	existingRecord := s.data[id]
	if len(existingRecord) != 0 {
		return ErrRecordAlreadyExists
	}

	newRecord := entity.PersistentRecord{
		ID:      record.GetID(),
		Version: 1,
		Start:   time.Now().UTC().Format("20060102150405"),
		End:     "",
		Data:    record.GetData(),
	}

	s.data[id] = append(s.data[id], newRecord)
	return nil
}

// UpdateRecord will update End of last version, and add a new record with incremented version.
func (s *PersistentRecordService) UpdateRecord(ctx context.Context, id int, updates map[string]*string) (entity.Record, error) {
	existingRecord := s.data[id]
	if len(existingRecord) == 0 {
		return nil, ErrRecordDoesNotExist
	}

	existingRecord[len(existingRecord)-1].End = time.Now().UTC().Format("20060102150405") // set end time for last version

	lastEntry := existingRecord[len(existingRecord)-1] // get the latest version of the record

	for key, value := range updates {
		if value == nil { // deletion update
			delete(lastEntry.GetData(), key)
		} else {
			lastEntry.GetData()[key] = *value
		}
	}

	newRecord := entity.PersistentRecord{
		ID:      lastEntry.ID,
		Version: lastEntry.Version + 1,
		Start:   lastEntry.End,
		End:     "",
		Data:    lastEntry.GetData(),
	}

	s.data[id] = append(s.data[id], newRecord)
	return newRecord.Copy(), nil
}
