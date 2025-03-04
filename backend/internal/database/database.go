package database

import (
	"database/sql"
)

var db *sql.DB

func ConnectDB() {
	var err error
	db, err = sql.Open("sqlite3", "data.db")
}

func ValidateUser() {

}
