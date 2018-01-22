package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"

	"github.com/mfranczy/prs/pkg/api"
	"github.com/mfranczy/prs/pkg/ui"
	"github.com/mfranczy/prs/pkg/db"
	"github.com/mfranczy/prs/pkg/web"
)

func main() {
	// set logger

	// add options for daemon

	db := db.NewPostgresDB()
	defer db.Close()

	r := mux.NewRouter()
	api.Init(r, api.Service{DB: db}) // add logger here?
	ui.Init(r)

	// pass handlers to http
	go web.ServeHTTP(r)

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	select {
	case <-term:
		log.Println("SIGTERM received, shutting down...")
	case err := <-web.CErr():
		log.Fatal("Unable to start http service: ", err)
	}

}
