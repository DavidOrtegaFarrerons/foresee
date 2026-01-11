package main

import (
	"html/template"
	"path/filepath"
)

type templateData struct {
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		//Other html base or partial files to be included
		/*patterns := []string{
			"ui/html/partials/"
		}*/

		ts := template.New(name)

		cache[name] = ts
	}

	return cache, nil
}
