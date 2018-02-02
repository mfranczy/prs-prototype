package ui

import (
	"net/http"
	"html/template"
	"runtime"
	"path"
	"path/filepath"

	"github.com/gorilla/mux"
)

var basePath string

func init() {
	_, file, _, _ := runtime.Caller(0)
	basePath = path.Join(filepath.Dir(file), "layouts")
}

func layoutRender(rw http.ResponseWriter, req *http.Request) {
	t := template.New("layoutPage")
	t, _ = t.ParseFiles(path.Join(basePath, "default.tmpl"))
	t.Execute(rw, nil)
}

func Init(r *mux.Router) {
	r.HandleFunc("/", layoutRender)
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir(basePath))),
	)
}
