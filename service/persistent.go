package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/regr76/timetravel/dbutils"
	"github.com/regr76/timetravel/entity"
)

const PersistentTimeFormat = "20060102150405" // use a consistent time format for start and end times of records

// PersistentRecordService is an in-memory implementation of RecordService.
type PersistentRecordService struct {
	db *sql.DB
}

func NewPersistentRecordService(db *sql.DB) PersistentRecordService {
	return PersistentRecordService{
		db: db,
	}
}

// GetRecord will retrieve record with latest version.
func (s *PersistentRecordService) GetRecord(ctx context.Context, id int) (entity.Record, error) {
	recordStr, err := dbutils.ReadLatestVersion(s.db, id)
	if err != nil {
		return nil, err
	}
	if recordStr == "" {
		return nil, ErrRecordDoesNotExist
	}

	// unmarshal data from json string to map[string]string
	output := &entity.PersistentRecord{}
	errUnMar := json.Unmarshal([]byte(recordStr), output)
	if errUnMar != nil {
		return nil, errUnMar
	}

	return output, nil
}

func (s *PersistentRecordService) GetVersion(ctx context.Context, id int, version int) (entity.Record, error) {
	recordStr, err := dbutils.ReadOneVersion(s.db, id, version)
	if err != nil {
		return nil, err
	}
	if recordStr == "" {
		return nil, ErrRecordDoesNotExist
	}

	// unmarshal data from json string to map[string]string
	output := &entity.PersistentRecord{}
	errUnMar := json.Unmarshal([]byte(recordStr), output)
	if errUnMar != nil {
		return nil, errUnMar
	}

	return output, nil
}

// ListRecords will retrieve record containing all versions.
func (s *PersistentRecordService) ListRecords(ctx context.Context, id int) (entity.VersionedRecords, error) {
	recordsStr, err := dbutils.ReadAllVersions(s.db, id)

	if err != nil {
		return nil, err
	}
	if len(recordsStr) == 0 {
		return nil, ErrRecordDoesNotExist
	}

	// unmarshal data from json string to map[string]string
	output := &entity.PersistentRecords{}
	for _, recordStr := range recordsStr {
		var record entity.PersistentRecord
		errUnMar := json.Unmarshal([]byte(recordStr), &record)
		if errUnMar != nil {
			return nil, errUnMar
		}
		output.Records = append(output.Records, record)
	}

	return output.Copy(), nil
}

// CreateRecord will create record with version 1.
func (s *PersistentRecordService) CreateRecord(ctx context.Context, record entity.Record) error {
	id := record.GetID()
	if id <= 0 {
		return ErrRecordIDInvalid
	}

	formattedData := `{}`
	data := record.GetData()
	if data != nil {
		formattedData = `{`
		for key, value := range data {
			formattedData += fmt.Sprintf(`"%s":"%s",`, key, value)
		}
		formattedData = formattedData[:len(formattedData)-1] + `}` // remove trailing comma and add closing bracket
	}

	return dbutils.WriteVersion(s.db, id, 1, time.Now().UTC().Format(PersistentTimeFormat), "", formattedData)
}

// UpdateRecord will update End of last version, and add a new record with incremented version.
func (s *PersistentRecordService) UpdateRecord(ctx context.Context, id int, updates map[string]*string) (entity.Record, error) {
	var version int
	copyOfLastVersion := &entity.PersistentRecord{}
	// first retrieve the record to see if an existing version exists
	recordStr, err := dbutils.ReadLatestVersion(s.db, id)

	if recordStr == "" || err != nil { // record does not exist, create new record with version 1
		version = 1
	} else { // record exists, need to update the End time of the last version and add a new version with updated data

		// unmarshal data from json string to map[string]string
		errUnMar := json.Unmarshal([]byte(recordStr), copyOfLastVersion)
		if errUnMar != nil {
			return nil, errUnMar
		}

		version = copyOfLastVersion.Version

		copyOfLastVersion.End = time.Now().UTC().Format(PersistentTimeFormat) // set end time for last version

		errWr := dbutils.UpdateVersion(
			s.db,
			copyOfLastVersion.GetID(),
			version,
			copyOfLastVersion.End,
		)
		if errWr != nil {
			return nil, errWr
		}

		version += 1 // increment version for the new version to be created
	}

	// now create the new version with updated data
	newData := copyOfLastVersion.Copy().(*entity.PersistentRecord).GetData() // copy data from last version for the new version
	if newData == nil {
		newData = map[string]string{}
	}
	for key, value := range updates {
		if value == nil { // deletion update
			delete(newData, key)
		} else {
			if value != nil {
				newData[key] = *value
			}
		}
	}

	formattedData := `{}`
	if newData != nil || len(newData) > 0 {
		formattedData = `{`
		for key, value := range newData {
			formattedData += fmt.Sprintf(`"%s":"%s",`, key, value)
		}
		formattedData = formattedData[:len(formattedData)-1] + `}` // remove trailing comma and add closing bracket
	}
	newVersion := &entity.PersistentRecord{
		ID:      id,
		Version: version,
		Start:   time.Now().UTC().Format(PersistentTimeFormat),
		End:     "",
		Data:    newData,
	}
	errWr := dbutils.WriteVersion(
		s.db,
		newVersion.GetID(),
		newVersion.Version,
		newVersion.Start,
		newVersion.End,
		formattedData,
	)
	if errWr != nil {
		return nil, errWr
	}

	return newVersion, nil
}

func (s *PersistentRecordService) ExportAllRecords(ctx context.Context) ([]string, error) {
	recordsStr, err := dbutils.ReadAllRows(s.db)
	if err != nil {
		return nil, err
	}
	return recordsStr, nil
}
