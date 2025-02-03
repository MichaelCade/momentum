package main

import (
	"log"
	"momentum/internal/database"
	"momentum/internal/routes"
	"net/http"
	"os"
)

func main() {
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

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
