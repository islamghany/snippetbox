package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"islamghany/snippetbox/pkg/models/mysql"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *mysql.SnippetModel
}

func main() {

	ADDR := os.Getenv("ADDR")
	DNS := os.Getenv("DNS")

	if len(ADDR) < 1 {
		ADDR = ":4000"
	}
	addr := flag.String("addr", ADDR, "Http Network address")
	dsn := flag.String("dns", DNS, "MySQL data source name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &mysql.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infoLog.Println("Starting server on port", *addr)
	err = srv.ListenAndServe()

	// will call os.Exit(1)
	errorLog.Fatal(err)
}

func openDB(dns string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
