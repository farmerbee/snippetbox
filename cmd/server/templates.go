package main

import (
	"blog/pkg/form"
	"blog/pkg/models"
	"html/template"
	"path"
	"path/filepath"
	"time"
)

// Templates store the data of HTML templates to be rendered
type Templates struct {
	Year     int
	Flash    string
	Form     *form.Form
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

// customize the style of time
func humanDate(t time.Time) string {
	time := t.Format("2006-01-02 15:04:06")
	return time
}

// register customized data processing functions with the templates
var tmplFuncs = template.FuncMap{
	"humanDate": humanDate,
}

// cache the templates in the memory
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := path.Base(page)

		//tmpl, err := template.ParseFiles(page)
		tmpl, err := template.New(name).Funcs(tmplFuncs).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		tmpl, err = tmpl.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}

		tmpl, err = tmpl.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}

		cache[name] = tmpl
	}

	return cache, nil
}
