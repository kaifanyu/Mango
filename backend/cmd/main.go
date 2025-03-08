package main

import (
	"log"
	"net/http"

	"fmt"
	"mango/internal/database"
	"mango/internal/routes"
)

func main() {
	// connect to database
	database.ConnectDB()

	// set up routes
	router := routes.SetUpRoutes()

	port := 8080
	fmt.Println("Server is up and running...")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), router))
}
