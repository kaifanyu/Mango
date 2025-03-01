package main

import (
	"fmt"
	"net/http"
	"log"

	"backend/internal/routes"
	"backend/internal/database"
)

func main() {
	// connect to database
	database.ConnectDB()
}
