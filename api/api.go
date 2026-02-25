package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/regr76/timetravel/service"
)

type API struct {
	router  *mux.Router
	records service.RecordService
}

func NewAPI(records service.RecordService) *API {
	return &API{records: records, router: mux.NewRouter()}
}

// generates all api routes
func (a *API) CreateRoutes(routes *mux.Router) {
	routes.Path("/records/{id}").HandlerFunc(a.GetRecords).Methods("GET")
	routes.Path("/records/{id}").HandlerFunc(a.PostRecords).Methods("POST")
}

func (a *API) SetupRouter() *mux.Router {
	service := service.NewInMemoryRecordService()
	api := NewAPI(&service)

	apiRoute := a.router.PathPrefix("/api/v1").Subrouter()
	apiRoute.Path("/health").HandlerFunc(HealthCheckHandler)
	api.CreateRoutes(apiRoute)

	return a.router
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	logError(err)
}
