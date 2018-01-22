package ui

import (
	"net/http"
	"html/template"

	"github.com/gorilla/mux"
)

func layoutRender(rw http.ResponseWriter, req *http.Request) {
	t := template.New("layoutPage")
	t, _ = t.ParseFiles("../../pkg/ui/layouts/default.tmpl")
	t.Execute(rw, nil)
}

func Init(r *mux.Router) {
	r.HandleFunc("/", layoutRender)
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("../../pkg/ui/layouts/"))),
	)
}
