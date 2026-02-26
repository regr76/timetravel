package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/regr76/timetravel/api/helpers"
	v1 "github.com/regr76/timetravel/api/v1"
	v2 "github.com/regr76/timetravel/api/v2"
	"github.com/regr76/timetravel/service"
)

type API struct {
	router         *mux.Router
	inMemRecords   service.RecordService
	persistRecords service.RecordService
}

func NewAPI(inMemRecords, persistRecords service.RecordService) *API {
	return &API{inMemRecords: inMemRecords, persistRecords: persistRecords, router: mux.NewRouter()}
}

func (a *API) InMemRecords() service.RecordService {
	return a.inMemRecords
}

func (a *API) PersistentRecords() service.RecordService {
	return a.persistRecords
}

// generates all api routes for V1 and adds them to the router
func (a *API) CreateRoutesV1(routes *mux.Router) {
	routes.Path("/records/{id}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v1.GetRecords(a, w, r)
	}).Methods("GET")

	routes.Path("/records/{id}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v1.PostRecords(a, w, r)
	}).Methods("POST")
}

// generates all api routes for V2 and adds them to the router
func (a *API) CreateRoutesV2(routes *mux.Router) {
	routes.Path("/records/{id}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v2.GetRecords(a, w, r)
	}).Methods("GET")

	routes.Path("/records/{id}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v2.PostRecords(a, w, r)
	}).Methods("POST")
}

func (a *API) SetupRouter() *mux.Router {
	inMemService := service.NewInMemoryRecordService()
	persistService := service.NewPersistentRecordService()
	api := NewAPI(&inMemService, &persistService)

	apiRoute1 := a.router.PathPrefix("/api/v1").Subrouter()
	apiRoute1.Path("/health").HandlerFunc(HealthCheckHandler)

	apiRoute2 := a.router.PathPrefix("/api/v2").Subrouter()
	apiRoute2.Path("/health").HandlerFunc(HealthCheckHandler)

	api.CreateRoutesV1(apiRoute1)
	api.CreateRoutesV2(apiRoute2)

	return a.router
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	helpers.LogError(err)
}
