package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

var cErr = make(chan error)

func ServeHTTP(r *mux.Router) {

	// api endpoints
	// user interface (i hope i can use frontend router)
	if err := http.ListenAndServe(":8080", r); err != nil {
		cErr <- err
	}
}

func CErr() <-chan error {
	return cErr
}
