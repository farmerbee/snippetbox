package main

import (
	"blog/pkg/models/mysql"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
)

type application struct {
	errLog   *log.Logger
	infoLog  *log.Logger
	snippets *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// handle different log
	infoLog := log.New(os.Stdout, "INFO\t", log.LstdFlags)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// user-defined host:port
	address := flag.String("addr", ":4000", "http network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true",
		"Mysql data source name(username:password@/dbname")
	flag.Parse()
	infoLog.Printf("starting server at %s", *address)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()


	// templates cached in the memory
	tmplCache, err := newTemplateCache(path.Join("..", "..", "app", "html"))
	if err != nil {
		errorLog.Println(err)
	}


	app := &application{
		errLog:   errorLog,
		infoLog:  infoLog,
		snippets: &mysql.SnippetModel{DB: db},
		templateCache: tmplCache,
	}

	// customize a http.Server
	srv := &http.Server{
		Addr:     *address,
		Handler:  app.route(),
		ErrorLog: errorLog,
	}
	err = srv.ListenAndServe()
	if err != nil {
		errorLog.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
