package routes

import (
	"momentum/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func InitializeRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/workout/today", handlers.GetWorkoutOfTheDay).Methods("GET")
	router.HandleFunc("/workout/log/cardio", handlers.LogCardioWorkout).Methods("POST")
	router.HandleFunc("/workout/log/weights", handlers.LogWeightsWorkout).Methods("POST")
	router.HandleFunc("/workout/logs/cardio", handlers.GetLoggedCardioWorkouts).Methods("GET")
	router.HandleFunc("/workout/logs/weights", handlers.GetLoggedWeightsWorkouts).Methods("GET")
	router.HandleFunc("/workout/weight-workouts", handlers.GetWeightWorkouts).Methods("GET")
	router.HandleFunc("/workout/last/cardio", handlers.GetLastLoggedCardioWorkout).Methods("GET")
	router.HandleFunc("/workout/last/weights", handlers.GetLastLoggedWeightsWorkout).Methods("GET")

	// Admin routes
	router.HandleFunc("/admin/add/{table}", handlers.AddRecord).Methods("POST")
	router.HandleFunc("/admin/update/{table}", handlers.UpdateRecord).Methods("POST")
	router.HandleFunc("/admin/delete/{table}", handlers.DeleteRecord).Methods("POST")
	router.HandleFunc("/admin/view/{table}", handlers.ViewRecords).Methods("GET")
	router.HandleFunc("/admin/empty/{table}", handlers.EmptyTable).Methods("POST")

	// Serve static files from the "web" directory
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web")))

	return router
}
