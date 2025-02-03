package routes

import (
	"momentum/internal/handlers"

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

	return router
}
