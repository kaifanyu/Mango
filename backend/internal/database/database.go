package database

import (
	"database/sql"
	"fmt"
	"log"

	// "github.com/mattn/go-sqlite3"
	"github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func ConnectDB() {
	var err error
	cfg := mysql.Config{
		User:   "root",
		Passwd: "para",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "mango",
	}

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Database connected")
	createTables()
}

// Create tables for authentication & progress tracking
func createTables() {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) UNIQUE NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	progressTable := `
	CREATE TABLE IF NOT EXISTS user_progress (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT NOT NULL,
		video_id VARCHAR(100) NOT NULL,
		progress INT DEFAULT 0,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	_, err := db.Exec(userTable)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	_, err = db.Exec(progressTable)
	if err != nil {
		log.Fatalf("Failed to create user_progress table: %v", err)
	}
}

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}

func ValidatePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func ValidateUser(username, password string) bool {
	var hashedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username=?", username).Scan(&hashedPassword)
	if err != nil {
		return false
	}

	return ValidatePassword(hashedPassword, password)
}

func RegisterUser(username, email, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		log.Fatal("Failed to Hash Password")
	}

	//insert the hashed password into the database
	query := `INSERT into users (username, email, password) VALUES (?, ?, ?)`
	_, err = db.Exec(query, username, email, hashedPassword)
	if err != nil {
		log.Fatal("Failed to insert user into database")
	}

	log.Println("User created successfully")
	return nil
}
