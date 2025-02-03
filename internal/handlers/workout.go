package handlers

import (
	"encoding/json"
	"log"
	"momentum/internal/models"
	"net/http"

	"github.com/gorilla/mux"
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

// Admin Section - Add, Update, Delete, View
// AddRecord handles the request to add a record to a table
func AddRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]

	var err error
	switch table {
	case "workouts":
		var workout models.Workout
		if err = json.NewDecoder(r.Body).Decode(&workout); err == nil {
			err = models.AddWorkout(workout)
		}
	case "weights_logs":
		var weightsLog models.WeightsLog
		if err = json.NewDecoder(r.Body).Decode(&weightsLog); err == nil {
			err = models.AddWeightsLog(weightsLog)
		}
	case "exercises":
		var exercise models.Exercise
		if err = json.NewDecoder(r.Body).Decode(&exercise); err == nil {
			err = models.AddExercise(exercise)
		}
	case "wods":
		var wod models.WOD
		if err = json.NewDecoder(r.Body).Decode(&wod); err == nil {
			err = models.AddWOD(wod)
		}
	case "weight_workouts":
		var weightWorkout models.WeightWorkout
		if err = json.NewDecoder(r.Body).Decode(&weightWorkout); err == nil {
			err = models.AddWeightWorkout(weightWorkout)
		}
	default:
		http.Error(w, "Invalid table name", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("Error adding record to %s: %v", table, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UpdateRecord handles the request to update a record in a table
func UpdateRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]

	var err error
	switch table {
	case "workouts":
		var workout models.Workout
		if err = json.NewDecoder(r.Body).Decode(&workout); err == nil {
			log.Printf("Received update request for workouts: %+v", workout)
			err = models.UpdateWorkout(workout)
		}
	case "weights_logs":
		var weightsLog models.WeightsLog
		if err = json.NewDecoder(r.Body).Decode(&weightsLog); err == nil {
			log.Printf("Received update request for weights_logs: %+v", weightsLog)
			err = models.UpdateWeightsLog(weightsLog)
		}
	case "exercises":
		var exercise models.Exercise
		if err = json.NewDecoder(r.Body).Decode(&exercise); err == nil {
			log.Printf("Received update request for exercises: %+v", exercise)
			err = models.UpdateExercise(exercise)
		}
	case "wods":
		var wod models.WOD
		if err = json.NewDecoder(r.Body).Decode(&wod); err == nil {
			log.Printf("Received update request for wods: %+v", wod)
			err = models.UpdateWOD(wod)
		}
	case "weight_workouts":
		var weightWorkout models.WeightWorkout
		if err = json.NewDecoder(r.Body).Decode(&weightWorkout); err == nil {
			log.Printf("Received update request for weight_workouts: %+v", weightWorkout)
			err = models.UpdateWeightWorkout(weightWorkout)
		}
	default:
		http.Error(w, "Invalid table name", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("Error updating record in %s: %v", table, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteRecord handles the request to delete a record from a table
func DeleteRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]

	var id struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
		log.Printf("Error decoding delete request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Received delete request for table %s with ID %d", table, id.ID)

	var err error
	switch table {
	case "workouts":
		err = models.DeleteWorkout(id.ID)
	case "weights_logs":
		err = models.DeleteWeightsLog(id.ID)
	case "exercises":
		err = models.DeleteExercise(id.ID)
	case "wods":
		err = models.DeleteWOD(id.ID)
	case "weight_workouts":
		err = models.DeleteWeightWorkout(id.ID)
	default:
		http.Error(w, "Invalid table name", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("Error deleting record from %s: %v", table, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ViewRecords handles the request to view records from a table
func ViewRecords(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]

	var records interface{}
	var err error
	switch table {
	case "workouts":
		records, err = models.ViewWorkouts()
	case "weights_logs":
		records, err = models.ViewWeightsLogs()
	case "exercises":
		records, err = models.ViewExercises()
	case "wods":
		records, err = models.ViewWODs()
	case "weight_workouts":
		records, err = models.ViewWeightWorkouts()
	default:
		http.Error(w, "Invalid table name", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("Error viewing records from %s: %v", table, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(records)
}

// EmptyTable handles the request to empty a table
func EmptyTable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	table := vars["table"]

	var err error
	switch table {
	case "workouts":
		err = models.EmptyWorkouts()
	case "weights_logs":
		err = models.EmptyWeightsLogs()
	case "exercises":
		err = models.EmptyExercises()
	case "wods":
		err = models.EmptyWODs()
	case "weight_workouts":
		err = models.EmptyWeightWorkouts()
	default:
		http.Error(w, "Invalid table name", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Printf("Error emptying table %s: %v", table, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
