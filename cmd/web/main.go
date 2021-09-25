package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {

	ADDR := os.Getenv("ADDR")
	if len(ADDR) < 1 {
		ADDR = ":4000"
	}
	addr := flag.String("addr", ADDR, "Http Network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/", app.homeHandler)
	mux.HandleFunc("/snippet", app.snippetHandler)
	mux.HandleFunc("/snippet/create", app.createSnippetHandler)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}
	infoLog.Println("Starting server on port", *addr)
	err := srv.ListenAndServe()

	// will call os.Exit(1)
	errorLog.Fatal(err)
}
