package v2

import "github.com/regr76/timetravel/service"

type Storage interface {
	InMemRecords() service.RecordService
	PersistentRecords() service.RecordService
}
