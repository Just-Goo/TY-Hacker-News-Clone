package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	
	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	appName string
	server  server
	debug   bool
	errLog  *log.Logger
	infoLog *log.Logger
	view    *jet.Set
	session *scs.SessionManager
}

type server struct {
	host string
	port string
	url  string
}

func main() {

	server := server{
		host: "localhost",
		port: "8080",
		url:  "http://localhost:8080",
	}

	db, err := openDB("root:@tcp(localhost:3306)/apidb")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Init app
	app := &application{
		server:  server,
		appName: "Hacker News Clone",
		debug:   true,
		infoLog: log.New(os.Stdout, "INFO\t", log.Ltime|log.Ldate|log.Lshortfile),
		errLog:  log.New(os.Stderr, "ERROR\t", log.Ltime|log.Ldate|log.Llongfile),
	}

	// Init jet template
	if app.debug {
		app.view = jet.NewSet(
			jet.NewOSFileSystemLoader("../../views"),
			jet.InDevelopmentMode(), // remove in production
		)
	} else {
		app.view = jet.NewSet(
			jet.NewOSFileSystemLoader("../../views"),
		)
	}

	// Init session mananger
	app.session = scs.New()
	app.session.Lifetime = 24 * time.Hour
	app.session.Cookie.Persist = true
	app.session.Cookie.Name = app.appName
	app.session.Cookie.Domain = app.server.host
	app.session.Cookie.SameSite = http.SameSiteStrictMode
	app.session.Store = mysqlstore.New(db)

	if err := app.listenAndServe(); err != nil {
		log.Fatal(err)
	}

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
