package handlers

import (
	"encoding/json"
	"log"
	"momentum/internal/models"
	"net/http"
)

// GetWorkoutOfTheDay handles the request to get the workout of the day
func GetWorkoutOfTheDay(w http.ResponseWriter, r *http.Request) {
	wod, err := models.FetchWorkoutOfTheDay()
	if err != nil {
		log.Printf("Error fetching workout of the day: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wod)
}

// LogCardioWorkout handles the request to log a cardio workout
func LogCardioWorkout(w http.ResponseWriter, r *http.Request) {
	var workout models.Workout
	if err := json.NewDecoder(r.Body).Decode(&workout); err != nil {
		log.Printf("Error decoding cardio workout log request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Received cardio workout log request: %+v", workout)
	if err := models.SaveWorkout(workout); err != nil {
		log.Printf("Error saving cardio workout: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// LogWeightsWorkout handles the request to log a weights workout
func LogWeightsWorkout(w http.ResponseWriter, r *http.Request) {
	var weightsLog models.WeightsLog
	if err := json.NewDecoder(r.Body).Decode(&weightsLog); err != nil {
		log.Printf("Error decoding weights log request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Received weights log request: %+v", weightsLog)
	if err := models.SaveWeightsLog(weightsLog); err != nil {
		log.Printf("Error saving weights log: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetLoggedCardioWorkouts handles the request to get all logged cardio workouts
func GetLoggedCardioWorkouts(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to fetch logged cardio workouts")
	workouts, err := models.FetchLoggedCardioWorkouts()
	if err != nil {
		log.Printf("Error fetching logged cardio workouts: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Fetched logged cardio workouts: %+v", workouts)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workouts)
}

// GetLoggedWeightsWorkouts handles the request to get all logged weights workouts
func GetLoggedWeightsWorkouts(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to fetch logged weights workouts")
	workouts, err := models.FetchLoggedWeightsWorkouts()
	if err != nil {
		log.Printf("Error fetching logged weights workouts: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Fetched logged weights workouts: %+v", workouts)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workouts)
}

// GetWeightWorkouts handles the request to get weight workouts
func GetWeightWorkouts(w http.ResponseWriter, r *http.Request) {
	workoutType := r.URL.Query().Get("type")
	if workoutType == "" {
		http.Error(w, "Missing workout type", http.StatusBadRequest)
		return
	}
	log.Printf("Fetching weight workouts for type: %s", workoutType)
	weightWorkouts, err := models.FetchWeightWorkouts(workoutType)
	if err != nil {
		log.Printf("Error fetching weight workouts: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Fetched weight workouts: %+v", weightWorkouts)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weightWorkouts)
}

// GetLastLoggedCardioWorkout handles the request to get the last logged cardio workout
func GetLastLoggedCardioWorkout(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching last logged cardio workout")
	workout, err := models.FetchLastLoggedCardioWorkout()
	if err != nil {
		log.Printf("Error fetching last logged cardio workout: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Fetched last logged cardio workout: %+v", workout)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

// GetLastLoggedWeightsWorkout handles the request to get the last logged weights workout
func GetLastLoggedWeightsWorkout(w http.ResponseWriter, r *http.Request) {
	workoutType := r.URL.Query().Get("type")
	if workoutType == "" {
		http.Error(w, "Missing workout type", http.StatusBadRequest)
		return
	}
	log.Printf("Fetching last logged weights workout for type: %s", workoutType)
	weightsLog, err := models.FetchLastLoggedWeightsWorkout(workoutType)
	if err != nil {
		log.Printf("Error fetching last logged weights workout: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if weightsLog == nil {
		http.Error(w, "No weights workout found", http.StatusNotFound)
		return
	}
	log.Printf("Fetched last logged weights workout: %+v", weightsLog)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weightsLog)
}
