package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/regr76/timetravel/api/helpers"
	v1 "github.com/regr76/timetravel/api/v1"
	"github.com/regr76/timetravel/service"
)

type API struct {
	router  *mux.Router
	records service.RecordService
}

func NewAPI(records service.RecordService) *API {
	return &API{records: records, router: mux.NewRouter()}
}

func (a *API) Records() service.RecordService {
	return a.records
}

// generates all api routes
func (a *API) CreateRoutesV1(routes *mux.Router) {
	routes.Path("/records/{id}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v1.GetRecords(a, w, r)
	}).Methods("GET")

	routes.Path("/records/{id}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v1.PostRecords(a, w, r)
	}).Methods("POST")
}

func (a *API) SetupRouter() *mux.Router {
	service := service.NewInMemoryRecordService()
	api := NewAPI(&service)

	apiRoute := a.router.PathPrefix("/api/v1").Subrouter()
	apiRoute.Path("/health").HandlerFunc(HealthCheckHandler)

	api.CreateRoutesV1(apiRoute)

	return a.router
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	helpers.LogError(err)
}
