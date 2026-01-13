package main

import (
	"foresee/cmd/web/viewmodels"
	"foresee/internal/models"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
)

type templateData struct {
	IsAuthenticated  bool
	Balance          int
	Flash            string
	Form             any
	MarketCategories []models.Category
	ResolverTypes    []models.ResolverType

	Markets []viewmodels.MarketView
	Market  viewmodels.MarketView
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	data := &templateData{
		Flash:            app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated:  isAuthenticated(r),
		MarketCategories: models.AllCategories(),
		ResolverTypes:    models.AllResolverTypes(),
		Balance:          0,
	}

	if !data.IsAuthenticated {
		return data
	}

	idStr := app.sessionManager.GetString(r.Context(), "authenticatedUserID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return data
	}
	balance, err := app.users.Balance(id)
	if err != nil {
		return data
	}

	data.Balance = balance

	return data
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
