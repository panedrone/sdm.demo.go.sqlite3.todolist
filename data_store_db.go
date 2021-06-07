package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // SQLite3
)

func initDb(ds *DataStore) {
	var err error
	ds.db, err = sql.Open("sqlite3", "./todo-list.sqlite")
	if err != nil {
		panic(err)
	}
}
