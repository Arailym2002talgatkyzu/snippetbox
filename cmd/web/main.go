package main

import (
	"arai/pkg/models/pg"
	"context"
	"flag"
	"fmt"
	"github.com/golangcollege/sessions"
	"github.com/jackc/pgx/v4/pgxpool"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      *pg.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")

	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	username := flag.String("username", "postgres", "username")
	password := flag.String("password", "650464", "password")
	host := flag.String("host", "localhost", "host")
	port := flag.String("port", "5432", "port")
	dbname := flag.String("dbname", "snippetbox", "dbname")
	connString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", *username, *password, *host, *port, *dbname)
	dsn := flag.String("dsn", connString, "PostgreSQL data source name")
	conn, err := pgxpool.Connect(context.Background(), *dsn)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close()
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &pg.SnippetModel{DB: conn},
		templateCache: templateCache,
	}

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = server.ListenAndServe()
	errorLog.Fatal(err)

}
