package v2

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/regr76/timetravel/api/helpers"
	"github.com/regr76/timetravel/service"
)

// GET /records/{id}
// GetRecord retrieves the record.
func GetVersion(a service.Storage, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	version := mux.Vars(r)["version"]
	idNumber, err1 := strconv.ParseInt(id, 10, 32)
	versionNumber, err2 := strconv.ParseInt(version, 10, 32)

	if err1 != nil || idNumber <= 0 {
		err := helpers.WriteError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		helpers.LogError(err)
		return
	}

	if err2 != nil || versionNumber <= 0 {
		err := helpers.WriteError(w, "invalid version; version must be a positive number", http.StatusBadRequest)
		helpers.LogError(err)
		return
	}

	record, err := a.PersistentRecords().GetVersion(
		ctx,
		int(idNumber),
		int(versionNumber),
	)
	if err != nil {
		err := helpers.WriteError(w, fmt.Sprintf("record of id %v version %v does not exist", idNumber, versionNumber), http.StatusBadRequest)
		helpers.LogError(err)
		return
	}

	err = helpers.WriteJSON(w, record, http.StatusOK)
	helpers.LogError(err)
}
