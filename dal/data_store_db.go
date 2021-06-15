package dal

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var ds = &DataStore{}

func initDb(ds *DataStore) (err error) {
	ds.db, err = sql.Open("sqlite3", "./todo-list.sqlite")
	return
}

func OpenDB() error {
	return ds.Open()
}

func CloseDB() error {
	return ds.Close()
}

func NewTasksDao() *TasksDao {
	return &TasksDao{Ds: ds}
}

func NewGroupsDao() *GroupsDao {
	return &GroupsDao{Ds: ds}
}
