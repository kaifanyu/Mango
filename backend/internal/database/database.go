package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func ConnectDB() {
	var err error
	cfg := mysql.Config{
		User:   "root",
		Passwd: "para",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "mango",
	}

	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	pingErr := DB.Ping()
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

	_, err := DB.Exec(userTable)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	_, err = DB.Exec(progressTable)
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

func ValidateImageHash(imageHash string) bool {
	var storedHash string

	err := DB.QueryRow("SELECT hash FROM registration").Scan(&storedHash)
	if err != nil {
		log.Print("Error storing hash", err)
		return false
	}

	return imageHash == storedHash
}

func ValidateUser(username, password string) bool {
	var hashedPassword string
	err := DB.QueryRow("SELECT password FROM users WHERE username=?", username).Scan(&hashedPassword)
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
	_, err = DB.Exec(query, username, email, hashedPassword)
	if err != nil {
		log.Fatal("Failed to insert user into database")
	}

	log.Println("User created successfully")
	return nil
}

func ExecDisplayQuery(query string, args ...interface{}) []map[string]interface{} {
	rows, err := DB.Query(query, args...)
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Printf("Error getting columns: %v", err)
		return nil
	}

	var results []map[string]interface{}

	// For each row
	for rows.Next() {
		// Create a slice of interface{} to hold the values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		// Set up pointers to each interface{} value
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}

		// Create a map for this row
		entry := make(map[string]interface{})

		for i, col := range columns {
			val := values[i]

			// Handle nil values
			if val == nil {
				entry[col] = nil
				continue
			}

			// Handle different data types from database
			switch v := val.(type) {
			case []byte:
				// Convert []byte to string for text/varchar fields
				entry[col] = string(v)
			default:
				// For all other types, use as is
				entry[col] = v
			}
		}

		results = append(results, entry)
	}

	return results
}

func GetUserID(username string) (int, error) {
	var userID int
	query := ` SELECT id FROM users WHERE username = ?`
	err := DB.QueryRow(query, username).Scan(&userID)
	return userID, err
}
