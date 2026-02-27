package api

import (
	"database/sql"
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
	persistRecords service.VersionedRecordService
	db             *sql.DB
}

func NewAPI(inMemRecords service.RecordService, persistRecords service.VersionedRecordService, db *sql.DB) *API {
	return &API{
		inMemRecords:   inMemRecords,
		persistRecords: persistRecords,
		router:         mux.NewRouter(),
		db:             db,
	}
}

func (a *API) InMemRecords() service.RecordService {
	return a.inMemRecords
}

func (a *API) PersistentRecords() service.VersionedRecordService {
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
	routes.Path("/records/{id}/list").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v2.ListRecords(a, w, r)
	}).Methods("GET")

	routes.Path("/records/{id}/versions/{version}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v2.GetVersion(a, w, r)
	}).Methods("GET")

	routes.Path("/records/{id}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v2.GetRecords(a, w, r)
	}).Methods("GET")

	routes.Path("/records/{id}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v2.PostRecords(a, w, r)
	}).Methods("POST")
}

func (a *API) SetupRouter(db *sql.DB) *mux.Router {
	inMemService := service.NewInMemoryRecordService()
	persistService := service.NewPersistentRecordService(db)
	api := NewAPI(&inMemService, &persistService, db)

	apiRoute1 := a.router.PathPrefix("/api/v1").Subrouter()
	apiRoute2 := a.router.PathPrefix("/api/v2").Subrouter()

	api.CreateRoutesV1(apiRoute1)
	api.CreateRoutesV2(apiRoute2)

	a.router.Path("/health").HandlerFunc(HealthCheckHandler)

	return a.router
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	helpers.LogError(err)
}
