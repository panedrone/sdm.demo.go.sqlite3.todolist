package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // SQLite3
	// _ "github.com/denisenkom/go-mssqldb" // SQL Server
	// _ "github.com/godror/godror"			// Oracle
	// only strings for MySQL (so far). see _prepareFetch below and related comments.
	// _ "github.com/go-sql-driver/mysql" // MySQL
	// _ "github.com/ziutek/mymysql/godrv" // MySQL
	// _ "github.com/lib/pq" // PostgeSQL
)

func initDb(ds *DataStore) {
	var err error
	ds.handle, err = sql.Open("sqlite3", "./todo-list.sqlite")
	if err != nil {
		panic(err)
	}
}
