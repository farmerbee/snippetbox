package main

import "net/http"

func (app *application) route() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/showsnippet", app.showSnippet)
	mux.HandleFunc("/createsnippet", app.createSnippet)


	// handle static files
	fileServer := http.FileServer(http.Dir("../../app/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))



	return mux
}
