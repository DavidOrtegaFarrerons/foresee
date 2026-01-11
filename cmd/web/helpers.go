package main

import (
	"bytes"
	"net/http"
)

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		panic("aaa")
		return
	}

	buf := new(bytes.Buffer)
	err := ts.Execute(buf, data)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}
