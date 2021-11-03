package main

import (
	"blog/pkg/form"
	"fmt"
	"net/http"
	"strconv"
)

// handle the home page
func (app *application) home(res http.ResponseWriter, req *http.Request) {
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
	id, err := strconv.Atoi(req.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(res)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		app.notFound(res)
		return
	}

	//flash := app.session.PopString(req, "flash")
	tmplData := &Templates{Snippet: snippet}

	app.render(res, req, "show.page.html", tmplData)
}

// handle to create a new record by the client with POST method
func (app *application) createSnippet(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		app.clientError(res, http.StatusBadRequest)
		return
	}

	title := req.PostForm.Get("title")
	content := req.PostForm.Get("content")
	expire := req.PostForm.Get("expire")

	form := form.NewForm(req.PostForm)
	form.Require("title", "content", "expire")
	form.MaxLength("title", 100)
	form.PermittedValues("expire", "1", "7", "30")
	if !form.Valid() {
		//fmt.Println(form)
		app.render(res, req, "snippetform.page.html", &Templates{Form: form})
		return
	}

	id, err := app.snippets.Insert(title, content, expire)
	if err != nil {
		app.errLog.Println(err)
		app.notFound(res)
		return
	}

	app.session.Put(req, "flash", "Snippet is created successfully!")

	http.Redirect(res, req, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "snippetform.page.html", &Templates{Form: &form.Form{}})
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.html", &Templates{Form: &form.Form{}})
}

func (app *application) userLoginForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "user post login data")
}

func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "handle user logout")
}

func (app *application) userSignUp(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.html", &Templates{Form: &form.Form{}})
	//fmt.Fprintln(w, "sign up page")
}

func (app *application) userSignUpForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "post user sign up data")
}

