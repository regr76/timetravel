package v2

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/regr76/timetravel/api/helpers"
	"github.com/regr76/timetravel/entity"
	"github.com/regr76/timetravel/service"
)

// POST /records/{id}
// if the record exists, the record is updated with a new version.
// if the record doesn't exist, the record is created.
func UpdateRecords(a service.Storage, w http.ResponseWriter, r *http.Request) {
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

	temp, err := a.PersistentRecords().UpdateRecord(ctx, int(idNumber), body)
	if err != nil {
		errInWriting := helpers.WriteError(w, helpers.ErrInternal.Error(), http.StatusInternalServerError)
		helpers.LogError(err)
		helpers.LogError(errInWriting)
		return
	}

	record := temp.(*entity.PersistentRecord)

	err = helpers.WriteJSON(w, record, http.StatusOK)
	helpers.LogError(err)
}
