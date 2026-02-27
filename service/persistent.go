package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/regr76/timetravel/entity"
)

const PersistentTimeFormat = "20060102150405" // use a consistent time format for start and end times of records

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

func (s *PersistentRecordService) GetVersion(ctx context.Context, id int, version int) (entity.Record, error) {
	records := s.data[id]
	if len(records) == 0 {
		return nil, ErrRecordDoesNotExist
	}

	for _, record := range records {
		if record.Version == version {
			copied, ok := record.Copy().(*entity.PersistentRecord) // copy is necessary so modifations to the record don't change the stored record
			if !ok {
				return nil, errors.New("failed to cast record to PersistentRecord")
			}
			return copied, nil
		}
	}

	return nil, ErrVersionDoesNotExist
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

	// single read for the entire list of versions. May be optimized to read last row only.
	existingRecord := s.data[id]
	if len(existingRecord) != 0 {
		return ErrRecordAlreadyExists
	}

	newRecord := entity.PersistentRecord{
		ID:      record.GetID(),
		Version: 1,
		Start:   time.Now().UTC().Format(PersistentTimeFormat),
		End:     "",
		Data:    record.GetData(),
	}

	s.data[id] = append(s.data[id], newRecord)
	return nil
}

// UpdateRecord will update End of last version, and add a new record with incremented version.
func (s *PersistentRecordService) UpdateRecord(ctx context.Context, id int, updates map[string]*string) (entity.Record, error) {
	// single read for the entire list of versions
	existingRecord := s.data[id]
	lenSoFar := len(existingRecord)
	if lenSoFar == 0 {
		return nil, ErrRecordDoesNotExist
	}

	copyOfLastVersion, ok := existingRecord[lenSoFar-1].Copy().(*entity.PersistentRecord) // get the latest version of the record
	if !ok {
		return nil, errors.New("failed to cast record to PersistentRecord")
	}

	copyOfLastVersion.End = time.Now().UTC().Format(PersistentTimeFormat) // set end time for last version

	newData := copyOfLastVersion.Copy().(*entity.PersistentRecord).GetData() // copy data from last version for the new version

	for key, value := range updates {
		if value == nil { // deletion update
			delete(newData, key)
		} else {
			newData[key] = *value
		}
	}

	newVersion := entity.PersistentRecord{
		ID:      copyOfLastVersion.ID,
		Version: copyOfLastVersion.Version + 1,
		Start:   copyOfLastVersion.End,
		End:     "",
		Data:    newData,
	}

	s.data[id][lenSoFar-1] = *copyOfLastVersion
	s.data[id] = append(s.data[id], newVersion)

	return newVersion.Copy(), nil
}
