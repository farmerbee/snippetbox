package main

import (
	"blog/pkg/form"
	"blog/pkg/models"
	"errors"
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
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	var id int
	if id, err = app.users.Authenticate(email, password); err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form := form.NewForm(nil)
			form.Errors.Add("generic", "Email address or password is invaid")
			app.render(w, r, "login.page.html", &Templates{
				Form: form,
			})
			return
		}
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "authenticatedId", id)

	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func (app *application) userLogout(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedId")
	app.session.Put(r, "flash", "You've logout successfully")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) userSignUp(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.html", &Templates{Form: &form.Form{}})
}

func (app *application) userSignUpForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
	}

	formData := form.NewForm(r.PostForm)
	userName := r.PostForm.Get("username")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	formData.Require("username", "email", "password")
	formData.MaxLength("username", 255)
	formData.MaxLength("password", 255)
	formData.MinLength("password", 8)
	formData.MatchPattern("email", form.EmailRX)

	if !formData.Valid() {
		app.render(w, r, "signup.page.html", &Templates{
			Form: formData,
		})
	}

	if err = app.users.Insert(userName, email, password); err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			formData.Errors.Add("email", "This email is already existed")
			app.render(w, r, "signup.page.html", &Templates{
				Form: formData,
			})
		} else {
			app.serverError(w, err)
		}

		return
	}

	app.session.Put(r, "flash", "Sign up successfully! Now you can login.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
