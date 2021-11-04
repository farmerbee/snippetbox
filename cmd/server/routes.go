package main

import (
	"net/http"

	"github.com/justinas/alice"

	"github.com/bmizerany/pat"
)

func (app *application) route() http.Handler {
	standardMiddleware := alice.New(app.RecoverPanic, app.LogRequest, SecureHeader)
	dynamicMiddleware := alice.New(app.session.Enable)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.RequireAuthenticate).ThenFunc(app.createSnippetForm))
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.RequireAuthenticate).ThenFunc(app.createSnippet))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.userLogin))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.userLoginForm))
	mux.Post("/user/logout", dynamicMiddleware.ThenFunc(app.userLogout))
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.userSignUp))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.userSignUpForm))

	// handle static files
	fileServer := http.FileServer(http.Dir("../../app/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
