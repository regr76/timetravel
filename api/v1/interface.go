package v1

import "github.com/regr76/timetravel/service"

type Storage interface {
	Records() service.RecordService
}
