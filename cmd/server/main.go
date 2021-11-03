package main

import (
	"blog/pkg/models/mysql"
	"database/sql"
	"flag"
	"github.com/golangcollege/sessions"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// application holds dependencies used(shared) by handlers
type application struct {
	errLog        *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	users         *mysql.UserModel
	templateCache map[string]*template.Template
	session       *sessions.Session
}

func main() {
	// handle different log
	infoLog := log.New(os.Stdout, "INFO\t", log.LstdFlags)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// user-defined host:port
	address := flag.String("addr", ":4000", "http network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true",
		"Mysql data source name(username:password@/dbname")
	secret := flag.String("secret", "LbJmYFguVTWxJ48AkdZE6Lva5hUWmK16", "Secret key")
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

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	app := &application{
		errLog:        errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		users:         &mysql.UserModel{DB: db},
		templateCache: tmplCache,
		session:       session,
	}

	// customize a http.Server
	srv := &http.Server{
		Addr:     *address,
		Handler:  app.route(),
		ErrorLog: errorLog,
	}
	//err = srv.ListenAndServe()
	cw, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	tlsPath := path.Join(cw, "..", "..", "tls")
	err = srv.ListenAndServeTLS(path.Join(tlsPath, "cert.pem"), path.Join(tlsPath, "key.pem"))
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
