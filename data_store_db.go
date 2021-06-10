package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func initDb(ds *DataStore) (err error) {
	ds.db, err = sql.Open("sqlite3", "./todo-list.sqlite")
	return
}
