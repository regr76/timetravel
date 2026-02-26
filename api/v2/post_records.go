package v2

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/regr76/timetravel/api/helpers"
	"github.com/regr76/timetravel/entity"
	"github.com/regr76/timetravel/service"
)

// POST /records/{id}
// if the record exists, the record is updated.
// if the record doesn't exist, the record is created.
func PostRecords(a service.Storage, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	idNumber, err := strconv.ParseInt(id, 10, 32)

	if err != nil || idNumber <= 0 {
		err := helpers.WriteError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		helpers.LogError(err)
		return
	}

	var body map[string]*string
	err = json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		err := helpers.WriteError(w, "invalid input; could not parse json", http.StatusBadRequest)
		helpers.LogError(err)
		return
	}

	// first retrieve the record
	var record *entity.InMemoryRecord
	_, err = a.PersistentRecords().GetRecord(
		ctx,
		int(idNumber),
	)

	if !errors.Is(err, service.ErrRecordDoesNotExist) { // record exists
		temp, err := a.PersistentRecords().UpdateRecord(ctx, int(idNumber), body)
		record = temp.(*entity.InMemoryRecord)

		if err != nil {
			errInWriting := helpers.WriteError(w, helpers.ErrInternal.Error(), http.StatusInternalServerError)
			helpers.LogError(err)
			helpers.LogError(errInWriting)
			return
		}

	} else { // record does not exist

		// exclude the delete updates
		recordMap := map[string]string{}
		for key, value := range body {
			if value != nil {
				recordMap[key] = *value
			}
		}

		record = &entity.InMemoryRecord{
			ID:   int(idNumber),
			Data: recordMap,
		}
		err = a.PersistentRecords().CreateRecord(ctx, record)
	}

	if err != nil {
		errInWriting := helpers.WriteError(w, helpers.ErrInternal.Error(), http.StatusInternalServerError)
		helpers.LogError(err)
		helpers.LogError(errInWriting)
		return
	}

	err = helpers.WriteJSON(w, record, http.StatusOK)
	helpers.LogError(err)
}
