package main

import (
	"errors"
	"fmt"
	"foresee/cmd/web/viewmodels"
	"foresee/internal/models"
	"foresee/internal/services"
	"foresee/internal/validator"
	"net/http"

	"github.com/google/uuid"
)

type signupForm struct {
	Username            string `form:"username"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type loginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type createMarketForm struct {
	Title               string `form:"title"`
	Description         string `form:"description"`
	Category            string `form:"category"`
	ResolverType        string `form:"resolver_type"`
	ExpiresAt           string `form:"expires_at"`
	validator.Validator `form:"-"`
}

type placeBetForm struct {
	OutcomeID string `form:"outcome_id"`
	Amount    int    `form:"amount"`
	validator.Validator
}

type resolveMarketForm struct {
	OutcomeID string `form:"outcome_id"`
	validator.Validator
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	markets, err := app.marketService.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	marketViews := make([]viewmodels.MarketView, 0, len(markets))
	for _, m := range markets {
		marketViews = append(
			marketViews,
			viewmodels.NewMarketView(*m, app.location),
		)
	}

	data := app.newTemplateData(r)
	data.Markets = marketViews

	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = signupForm{}
	app.render(w, http.StatusOK, "signup.html", data)
}

func (app *application) signupPost(w http.ResponseWriter, r *http.Request) {
	var form signupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.serverError(w, err)
		return
	}

	form.CheckField(validator.NotBlank(form.Username), "username", "The username cannot be empty")
	form.CheckField(validator.MinChars(form.Username, 4), "username", "This field must have at least 4 characters")
	form.CheckField(validator.NotBlank(form.Email), "email", "The email cannot be empty")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "The email must be a valid address")
	form.CheckField(validator.NotBlank(form.Password), "password", "The password cannot be empty")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "The password must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	err = app.users.Insert(form.Username, form.Email, form.Password)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrUsernameAlreadyExists):
			form.AddFieldError("username", "This username already exists, use a different one")

		case errors.Is(err, models.ErrEmailAlreadyExists):
			form.AddFieldError("email", "This email is already in use, use a different one")

		default:
			app.serverError(w, err)
		}

		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Account created successfully, please log in")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = loginForm{}
	app.render(w, http.StatusOK, "login.html", data)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {
	var form loginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.serverError(w, err)
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "The email cannot be empty")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "The email must be in a valid format")
	form.CheckField(validator.NotBlank(form.Password), "password", "The password cannot be empty")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.html", data)
			return
		}
	}
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id.String())

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) account(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	userID, err := app.getUserId(r)
	if err != nil {
		app.serverError(w, err)
	}

	userBetHistory, err := app.betService.GetUserBetHistory(userID)
	if err != nil {
		app.serverError(w, err)
	}

	marketsPendingResolution, err := app.marketService.PendingResolution(userID)
	if err != nil {
		app.serverError(w, err)
	}
	data.PendingResolutions = marketsPendingResolution

	data.BetHistory = userBetHistory
	app.render(w, http.StatusOK, "account.html", data)
}

func (app *application) createMarket(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = createMarketForm{}
	app.render(w, http.StatusOK, "create_market.html", data)
}

func (app *application) createMarketPost(w http.ResponseWriter, r *http.Request) {
	var form createMarketForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.serverError(w, err)
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "Market must have a name")
	form.CheckField(validator.MinChars(form.Title, 4), "title", "Title must be at least 4 characters long")
	form.CheckField(validator.NotBlank(form.Description), "description", "Description cannot be empty")
	form.CheckField(validator.NotBlank(form.Category), "category", "Description cannot be empty")
	form.CheckField(validator.PermittedValue(models.Category(form.Category), models.AllCategories()...), "category", "The category must be valid")
	form.CheckField(validator.NotBlank(form.ResolverType), "resolverType", "Description cannot be empty")
	form.CheckField(validator.PermittedValue(models.ResolverType(form.ResolverType), models.AllResolverTypes()...), "resolverType", "The resolver type must be valid")
	form.CheckField(validator.NotBlank(form.ExpiresAt), "expiresAt", "The expiry date must be fulfilled")
	form.CheckField(validator.IsValidDate(form.ExpiresAt), "expiresAt", "The expiry date must be valid and must not be in the past")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create_market.html", data)
		return
	}

	id, err := app.getUserId(r)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.marketService.Create(
		form.Title,
		form.Description,
		form.Category,
		form.ResolverType,
		form.ExpiresAt,
		id,
	)

	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, r, "/markets", http.StatusSeeOther)
}

func (app *application) viewMarket(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	m, err := app.marketService.Get(id)
	if err != nil {
		app.serverError(w, err)
	}

	data := app.newTemplateData(r)
	data.Market = viewmodels.NewMarketView(m, app.location)
	app.infoLog.Printf("%q", data.Market.Outcomes)
	data.Form = placeBetForm{}
	app.render(w, http.StatusOK, "detail_market.html", data)
}

func (app *application) dailyClaimPost(w http.ResponseWriter, r *http.Request) {
	id, err := app.getUserId(r)
	if err != nil {
		app.serverError(w, err)
		return
	}

	s := services.UserService{Users: app.users}
	redirectTo := r.Header.Get("Referer")
	if redirectTo == "" {
		redirectTo = "/"
	}
	err = s.ClaimDailyReward(id)
	if err != nil {

		if errors.Is(err, services.ErrDailyRewardNotAvailable) {
			app.sessionManager.Put(r.Context(), "flash", err.Error())
			http.Redirect(w, r, redirectTo, http.StatusSeeOther)
			return
		}

		app.serverError(w, err)
	}

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Your reward of %d has been added to your balance", services.DailyRewardAmmount))
	http.Redirect(w, r, redirectTo, http.StatusSeeOther)
}

func (app *application) createBetPost(w http.ResponseWriter, r *http.Request) {
	marketID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var form placeBetForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.serverError(w, err)
		return
	}

	form.Validator.CheckField(validator.NotBlank(form.OutcomeID), "outcome_id", "There has been an error with your request, please try again later")
	form.Validator.CheckField(validator.MinNumber(form.Amount, models.MinimumBetAmount), "amount", "The minimum bet is 100")

	redirectTo := r.Header.Get("Referer")
	if redirectTo == "" {
		redirectTo = "/"
	}

	if !form.Valid() {
		app.sessionManager.Put(r.Context(), "flash", "TODO FIX THIS ERRORS")
		http.Redirect(w, r, redirectTo, http.StatusSeeOther)
		return
	}

	outcomeID, err := uuid.Parse(form.OutcomeID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	userID, err := app.getUserId(r)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.betService.Place(userID, marketID, outcomeID, form.Amount)
	if err != nil {
		app.sessionManager.Put(r.Context(), "flash_error", err.Error())
		http.Redirect(w, r, redirectTo, http.StatusSeeOther)
		return
	}

	//TODO MAKE THIS MESSAGE PRETTIER / SAY THE AMOUNT AND THE MARKET AND OUTCOME
	app.sessionManager.Put(r.Context(), "flash", "Your bet has been submitted successfully")
	http.Redirect(w, r, redirectTo, http.StatusSeeOther)
}

func (app *application) resolveMarket(w http.ResponseWriter, r *http.Request) {
	marketID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	m, err := app.marketService.Get(marketID)
	if err != nil {
		app.serverError(w, err)
	}

	userID, err := app.getUserId(r)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if *m.ResolverRef != userID {
		app.clientError(w, http.StatusForbidden)
		return
	}

	data := app.newTemplateData(r)
	data.Market = viewmodels.NewMarketView(m, app.location)
	data.Form = resolveMarketForm{}
	app.render(w, http.StatusOK, "resolve_market.html", data)
}

func (app *application) resolveMarketPost(w http.ResponseWriter, r *http.Request) {
	marketID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	var form resolveMarketForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.serverError(w, err)
		return
	}

	form.Validator.CheckField(validator.NotBlank(form.OutcomeID), "outcome_id", "The outcome must not be empty")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		m, err := app.marketService.Get(marketID)
		if err != nil {
			app.serverError(w, err)
			return
		}
		data.Market = viewmodels.NewMarketView(m, app.location)
		app.render(w, http.StatusUnprocessableEntity, "resolve_market.html", data)
		return
	}

	outcomeID, err := uuid.Parse(form.OutcomeID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	userID, err := app.getUserId(r)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.marketService.ResolveMarket(marketID, userID, outcomeID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Thanks for resolving the market")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
