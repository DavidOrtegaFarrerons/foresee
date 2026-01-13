package main

import (
	"errors"
	"flag"
	"foresee/internal/models"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type application struct {
	infoLog        *log.Logger
	errorLog       *log.Logger
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	users          *models.UserModel
	markets        *models.MarketModel
	sessionManager *scs.SessionManager
	location       *time.Location
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("db", "postgres://user:pass@postgres:5432/foresee?sslmode=disable", "Database connection string")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)
	tc, err := newTemplateCache()
	if err != nil {
		panic(err)
	}

	location, err := time.LoadLocation("Europe/Madrid")
	if err != nil {
		panic(err)
	}

	db, err := openDb(*dsn)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/migrations",
		"postgres",
		driver,
	)

	sesssionManager := scs.New()
	sesssionManager.Store = postgresstore.New(db)
	sesssionManager.Lifetime = 12 * time.Hour
	sesssionManager.Cookie.Secure = true

	if err != nil {
		log.Fatal(err)
	}

	if m == nil {
		log.Fatal("migration driver could not be connected")
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("No migrations have been applied")
		} else {
			log.Fatal(err)
		}
	}

	app := application{
		infoLog:       infoLog,
		errorLog:      errorLog,
		templateCache: tc,
		formDecoder:   form.NewDecoder(),
		users: &models.UserModel{
			DB: db,
		},
		markets: &models.MarketModel{
			DB: db,
		},
		sessionManager: sesssionManager,
		location:       location,
	}

	log.Printf("Starting server on %s", *addr)
	err = http.ListenAndServe(*addr, app.routes())
	log.Fatal(err)
}
