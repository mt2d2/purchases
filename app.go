package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type app struct {
	db *sqlx.DB
}

func newApp() *app {
	db, err := sqlx.Open("sqlite3", *db)
	if err != nil {
		log.Panicln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Panicln(err)
	}

	return &app{db}
}

func (app *app) destroy() {
	app.db.Close()
}
