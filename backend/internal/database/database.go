package database

import (
	"database/sql"
	"fmt"
	"github.com/mattin/go-sqlite3"
)

var db *sql.DB

func ConnectDB() {
	var err error
	db, err = sql.Open("sqlite3", "data.db")
	
}