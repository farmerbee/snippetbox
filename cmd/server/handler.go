package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) home(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}

	snippets, err := app.snippets.Lasted()
	if err != nil {
		app.serverError(res, err)
	}

	//for _, s := range snippets {
	//	fmt.Fprintf(res, "%v\n", *s)
	//}

	temlsData := &Templates{
		Snippets: snippets,
	}

	app.render(res, req, "home.page.html", temlsData)

	//files := []string{
	//	"../../app/html/home.page.html",
	//	"../../app/html/base.layout.html",
	//	"../../app/html/footer.partial.html",
	//}
	////temp, err := template.ParseFiles(tempPath)
	//temp, err := template.ParseFiles(files...)
	//if err != nil {
	//	//app.errLog.Println(err)
	//	//res.WriteHeader(http.StatusInternalServerError)
	//	app.serverError(res, err)
	//	return
	//}
	//
	//err = temp.Execute(res, &temlsData)
	//if err != nil {
	//	//app.errLog.Println(err)
	//	//res.WriteHeader(http.StatusInternalServerError)
	//	app.serverError(res, err)
	//}

}

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
	//files := []string{
	//	"../../app/html/show.page.html",
	//	"../../app/html/base.layout.html",
	//	"../../app/html/footer.partial.html",
	//}
	//
	//tmpl, err := template.ParseFiles(files...)
	//if err != nil {
	//	app.serverError(res, err)
	//}
	//
	//err = tmpl.Execute(res, *snippet)
	//if err != nil {
	//	app.serverError(res, err)
	//}
}

func (app *application) createSnippet(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		res.Header().Set("allow", http.MethodGet)
		//res.WriteHeader(http.StatusMethodNotAllowed)
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
	//res.Write([]byte("you can create snippet later"))
}
