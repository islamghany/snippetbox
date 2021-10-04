package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHTTPRequest)
	mux.HandleFunc("/", app.homeHandler)
	mux.HandleFunc("/snippet", app.snippetHandler)
	mux.HandleFunc("/snippet/create", app.createSnippetHandler)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
