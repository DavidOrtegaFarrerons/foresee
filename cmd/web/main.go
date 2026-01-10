package main

import (
	"flag"
	"log"
	"net/http"
)

type application struct{}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	app := application{}
	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, app.routes())
	log.Fatal(err)
}
