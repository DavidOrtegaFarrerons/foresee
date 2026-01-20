package main

import (
	"foresee/cmd/web/viewmodels"
	"foresee/internal/models"
	"foresee/internal/services"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

type templateData struct {
	IsAuthenticated     bool
	Balance             int
	CanClaimDailyReward bool
	Flash               string
	FlashError          string
	Form                any
	MarketCategories    []models.Category
	ResolverTypes       []models.ResolverType
	BetHistory          []models.BetHistoryRow

	Markets            []viewmodels.MarketView
	Market             viewmodels.MarketView
	PendingResolutions []models.Market
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	data := &templateData{
		Flash:            app.sessionManager.PopString(r.Context(), "flash"),
		FlashError:       app.sessionManager.PopString(r.Context(), "flash_error"),
		IsAuthenticated:  isAuthenticated(r),
		MarketCategories: models.AllCategories(),
		ResolverTypes:    models.AllResolverTypes(),
		Balance:          0,
	}

	if !data.IsAuthenticated {
		return data
	}

	id, err := app.getUserId(r)
	if err != nil {
		return data
	}

	balance, lastClaimedAt, err := app.users.GetTemplateInfo(id)
	if err != nil {
		return data
	}

	data.Balance = balance
	data.CanClaimDailyReward = services.CanClaimReward(time.Now(), lastClaimedAt.Time)

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
