package main

import (
	"log"
	"net/http"
	"time"

	"github.com/regr76/timetravel/api"
)

func main() {
	filename := "timetravel.db"
	log.Printf("initializing database with file %s", filename)

	db, err := initDB(filename)
	if err != nil {
		log.Fatal(err)
	}
	// close and check the error
	defer func() {
		if cerr := db.Close(); cerr != nil {
			log.Printf("db close: %v", cerr)
		}
	}()
	// Verify the connection
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	app := api.NewAPI(nil, nil, db)
	router := app.SetupRouter(db)

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
