package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// handle the home page
func (app *application) home(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}

	snippets, err := app.snippets.Lasted()
	if err != nil {
		app.serverError(res, err)
	}

	temlsData := &Templates{
		Snippets: snippets,
	}

	app.render(res, req, "home.page.html", temlsData)
}

// handle to show detail of specific record by id
func (app *application) showSnippet(res http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(req.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(res)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		app.notFound(res)
		return
	}
	tmplData := &Templates{Snippet: snippet}

	app.render(res, req, "show.page.html", tmplData)
}

// handle to create a new record by the client with POST method
func (app *application) createSnippet(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.Header().Set("allow", http.MethodGet)
		app.clientError(res, http.StatusMethodNotAllowed)
		return
	}

	title := "mock"
	content := "this is mock data"
	expire := 0
	id, err := app.snippets.Insert(title, content, strconv.Itoa(expire))
	if err != nil {
		app.notFound(res)
		return
	}

	http.Redirect(res, req, fmt.Sprintf("/showsnippet?id=%d", id), http.StatusSeeOther)
}
