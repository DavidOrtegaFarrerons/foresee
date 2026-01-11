package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
)

type application struct {
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	tc, err := newTemplateCache()
	if err != nil {
		panic(err)
	}
	app := application{
		templateCache: tc,
	}
	log.Printf("Starting server on %s", *addr)
	err = http.ListenAndServe(*addr, app.routes())
	log.Fatal(err)
}
