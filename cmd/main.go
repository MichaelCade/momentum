package main

import (
	"log"
	"momentum/internal/database"
	"momentum/internal/routes"
	"net/http"
	"os"
)

func main() {
	// Set the DATABASE_URL environment variable for testing purposes
	os.Setenv("DATABASE_URL", "postgres://postgres:Passw0rd999!@192.168.169.104:5432/momentum?sslmode=disable")

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
