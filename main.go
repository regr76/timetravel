package main

import (
	"log"
	"net/http"
	"time"

	"github.com/regr76/timetravel/api"
)

func main() {
	app := api.NewAPI(nil)
	router := app.SetupRouter()

	address := "127.0.0.1:8000"
	srv := &http.Server{
		Handler:      router,
		Addr:         address,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("listening on %s", address)
	log.Fatal(srv.ListenAndServe())
}
