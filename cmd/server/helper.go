package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serverError(res http.ResponseWriter, err error) {
	// record the error message  and the stack trace
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errLog.Output(2, trace)
	app.errLog.Println(trace)

	// response an error message to the remote
	http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(res http.ResponseWriter, statusCode int) {
	http.Error(res, http.StatusText(statusCode), statusCode)
}

func (app *application) notFound(res http.ResponseWriter) {
	app.clientError(res, http.StatusNotFound)
}

// render and send the html file with specified templates and data
func (app *application) render(res http.ResponseWriter, req *http.Request, pageName string, tmplData *Templates) {
	tmpl, ok := app.templateCache[pageName]
	if !ok {
		app.serverError(res, fmt.Errorf("the page %s is not found", pageName))
		return
	}

	buf := new(bytes.Buffer)
	err := tmpl.Execute(buf, app.addTmplYear(tmplData, req))
	if err != nil {
		app.serverError(res, err)
		return
	}

	buf.WriteTo(res)
}


// generate the year data before rendering a page
func (app *application) addTmplYear(tmplData *Templates, req *http.Request) *Templates {
	if tmplData == nil {
		tmplData = &Templates{}
	}
	tmplData.Year = time.Now().Year()

	return tmplData
}