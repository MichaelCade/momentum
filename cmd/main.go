package main

import (
	"log"
	"momentum/internal/database"
	"momentum/internal/routes"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Check if the DATABASE_URL environment variable is set
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatalln("DATABASE_URL environment variable is not set")
	}

	database.InitDB() // Initialize the database connection

	router := routes.InitializeRoutes() // Initialize routes using gorilla/mux

	// Serve static files from the "web" directory
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))

	router.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "web/style.css")
	})

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
