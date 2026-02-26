package service

import (
	"context"
	"errors"

	"github.com/regr76/timetravel/entity"
)

var ErrRecordDoesNotExist = errors.New("record with that id does not exist")
var ErrRecordIDInvalid = errors.New("record id must >= 0")
var ErrRecordAlreadyExists = errors.New("record already exists")

// InMemoryRecordService is an in-memory implementation of RecordService.
type InMemoryRecordService struct {
	data map[int]entity.InMemoryRecord
}

func NewInMemoryRecordService() InMemoryRecordService {
	return InMemoryRecordService{
		data: map[int]entity.InMemoryRecord{},
	}
}

func (s *InMemoryRecordService) GetRecord(ctx context.Context, id int) (entity.Record, error) {
	record := s.data[id]
	if record.GetID() == 0 {
		return nil, ErrRecordDoesNotExist
	}

	copied, ok := record.Copy().(*entity.InMemoryRecord) // copy is necessary so modifations to the record don't change the stored record
	if !ok {
		return nil, errors.New("failed to cast record to InMemoryRecord")
	}
	return copied, nil
}

func (s *InMemoryRecordService) CreateRecord(ctx context.Context, record entity.Record) error {
	id := record.GetID()
	if id <= 0 {
		return ErrRecordIDInvalid
	}

	existingRecord := s.data[id]
	if existingRecord.GetID() != 0 {
		return ErrRecordAlreadyExists
	}

	s.data[id] = *(record.(*entity.InMemoryRecord))
	return nil
}

func (s *InMemoryRecordService) UpdateRecord(ctx context.Context, id int, updates map[string]*string) (entity.Record, error) {
	entry := s.data[id]
	if entry.GetID() == 0 {
		return nil, ErrRecordDoesNotExist
	}

	for key, value := range updates {
		if value == nil { // deletion update
			delete(entry.GetData(), key)
		} else {
			entry.GetData()[key] = *value
		}
	}

	return entry.Copy(), nil
}
