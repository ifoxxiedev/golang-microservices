package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

type Config struct {
	Repo   data.Repository
	Client *http.Client
}

var counts int64
var maxCounts int64 = 10

func main() {
	log.Println("Starting the authentication service")
	// TODO Connect to DB
	db := connectToDB()
	if db == nil {
		log.Panic("Failed to connect to database")
	}

	// set up config
	app := Config{
		Client: &http.Client{},
	}
	app.setupRepo(db)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

func openDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	for {
		db, err := openDb(dsn)
		if err != nil {
			log.Println("Postgres is not ready yet...", err)
			counts += 1
		} else {
			log.Println("Postgres is ready!")
			return db
		}

		if counts > maxCounts {
			log.Printf("Postgres is not ready after 10 attempts. Exiting... %v\n", err)
			return nil
		}

		log.Printf("Backing off for 5 seconds before trying again...")
		time.Sleep(5 * time.Second)
		continue
	}
}

func (app *Config) setupRepo(conn *sql.DB) {
	repo := data.NewPostgresRepository(conn)
	app.Repo = repo
}
