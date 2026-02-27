package v2

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/regr76/timetravel/api/helpers"
	"github.com/regr76/timetravel/service"
)

// GET /records/{id}/version/{version}
// ListRecords retrieves the record including all versions
func ListRecords(a service.Storage, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	idNumber, err := strconv.ParseInt(id, 10, 32)
	if err != nil || idNumber <= 0 {
		err := helpers.WriteError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		helpers.LogError(err)
		return
	}

	records, err := a.PersistentRecords().ListRecords(
		ctx,
		int(idNumber),
	)
	if err != nil {
		err := helpers.WriteError(w, fmt.Sprintf("record of id %v does not exist", idNumber), http.StatusBadRequest)
		helpers.LogError(err)
		return
	}

	err = helpers.WriteJSON(w, records, http.StatusOK)
	helpers.LogError(err)
}
