package models

import (
	"database/sql"
	"log"
	"momentum/internal/database"
	"time"
)

// Workout represents a workout entry in the database.
type Workout struct {
	ID       int       `json:"id"`
	Type     string    `json:"type"`     // Type of workout (e.g., cardio, strength)
	Duration int       `json:"duration"` // Duration in minutes
	Distance float64   `json:"distance"` // Distance in kilometers (if applicable)
	Date     time.Time `json:"date"`     // Date of the workout
}

// WeightsLog represents a log entry for a weights workout.
type WeightsLog struct {
	ID          int        `json:"id"`
	WorkoutType string     `json:"workout_type" db:"workout_type"`
	Exercises   []Exercise `json:"exercises"`
	Date        time.Time  `json:"date"`
}

// Exercise represents an exercise entry in a weights log.
type Exercise struct {
	ID           int    `json:"id"`
	WeightsLogID int    `json:"weights_log_id" db:"weights_log_id"`
	Name         string `json:"name"`
	Set1         int    `json:"set1"`
	Set2         int    `json:"set2"`
	Set3         int    `json:"set3"`
}

// WOD represents a workout of the day entry in the database.
type WOD struct {
	ID       int       `json:"id"`
	Type     string    `json:"type"`     // Type of workout (e.g., cardio, strength)
	Duration int       `json:"duration"` // Duration in minutes
	Distance float64   `json:"distance"` // Distance in kilometers (if applicable)
	Date     time.Time `json:"date"`     // Date of the workout
}

// WeightWorkout represents a weight workout entry in the database.
type WeightWorkout struct {
	ID          int    `json:"id"`
	WorkoutType string `json:"workout_type" db:"workout_type"`
	Exercise    string `json:"exercise" db:"exercise"`
}

// FetchWorkoutOfTheDay retrieves a random workout of the day from the database.
func FetchWorkoutOfTheDay() (*WOD, error) {
	var wod WOD
	err := database.DB.Get(&wod, "SELECT * FROM wods ORDER BY RANDOM() LIMIT 1")
	if err != nil {
		log.Printf("Error fetching workout of the day: %v", err)
		return nil, err
	}
	return &wod, nil
}

// SaveWorkout saves a new workout to the database.
func SaveWorkout(workout Workout) error {
	_, err := database.DB.NamedExec(`INSERT INTO workouts (type, duration, distance, date) VALUES (:type, :duration, :distance, :date)`, &workout)
	return err
}

// SaveWeightsLog saves a new weights log to the database.
func SaveWeightsLog(weightsLog WeightsLog) error {
	weightsLog.Date = time.Now() // Set the current time and date
	tx := database.DB.MustBegin()
	var weightsLogID int
	err := tx.QueryRowx(`INSERT INTO weights_logs (workout_type, date) VALUES ($1, $2) RETURNING id`, weightsLog.WorkoutType, weightsLog.Date).Scan(&weightsLogID)
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, exercise := range weightsLog.Exercises {
		exercise.WeightsLogID = weightsLogID
		_, err := tx.NamedExec(`INSERT INTO exercises (weights_log_id, name, set1, set2, set3) VALUES (:weights_log_id, :name, :set1, :set2, :set3)`, &exercise)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// FetchLoggedCardioWorkouts retrieves all logged cardio workouts from the database.
func FetchLoggedCardioWorkouts() ([]Workout, error) {
	var workouts []Workout
	err := database.DB.Select(&workouts, "SELECT * FROM workouts WHERE type IN ('run', 'bike', 'row', 'walk', 'crosstrainer') ORDER BY date DESC")
	if err != nil {
		log.Printf("Error fetching logged cardio workouts: %v", err)
		return nil, err
	}
	return workouts, nil
}

// FetchLoggedWeightsWorkouts retrieves all logged weights workouts from the database.
func FetchLoggedWeightsWorkouts() ([]WeightsLog, error) {
	var weightsLogs []WeightsLog
	err := database.DB.Select(&weightsLogs, "SELECT * FROM weights_logs ORDER BY date DESC")
	if err != nil {
		log.Printf("Error fetching logged weights workouts: %v", err)
		return nil, err
	}
	return weightsLogs, nil
}

// FetchWeightWorkouts retrieves weight workouts from the database.
func FetchWeightWorkouts(workoutType string) ([]WeightWorkout, error) {
	var weightWorkouts []WeightWorkout
	err := database.DB.Select(&weightWorkouts, "SELECT * FROM weight_workouts WHERE workout_type=$1", workoutType)
	if err != nil {
		log.Printf("Error fetching weight workouts: %v", err)
		return nil, err
	}
	return weightWorkouts, nil
}

// FetchLastLoggedCardioWorkout retrieves the last logged cardio workout from the database.
func FetchLastLoggedCardioWorkout() (*Workout, error) {
	var workout Workout
	err := database.DB.Get(&workout, "SELECT * FROM workouts WHERE type IN ('run', 'bike', 'row', 'walk', 'crosstrainer') ORDER BY date DESC LIMIT 1")
	if err != nil {
		log.Printf("Error fetching last logged cardio workout: %v", err)
		return nil, err
	}
	return &workout, nil
}

// FetchLastLoggedWeightsWorkout retrieves the last logged weights workout for a specific type from the database.
func FetchLastLoggedWeightsWorkout(workoutType string) (*WeightsLog, error) {
	var weightsLog WeightsLog
	err := database.DB.Get(&weightsLog, "SELECT * FROM weights_logs WHERE workout_type=$1 ORDER BY date DESC LIMIT 1", workoutType)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No weights workout found for type: %s", workoutType)
			return nil, nil
		}
		log.Printf("Error fetching last logged weights workout: %v", err)
		return nil, err
	}

	log.Printf("Fetched weights log: %+v", weightsLog)

	var exercises []Exercise
	err = database.DB.Select(&exercises, "SELECT * FROM exercises WHERE weights_log_id=$1", weightsLog.ID)
	if err != nil {
		log.Printf("Error fetching exercises for weights log ID %d: %v", weightsLog.ID, err)
		return nil, err
	}
	log.Printf("Fetched exercises for weights log ID %d: %+v", weightsLog.ID, exercises)
	weightsLog.Exercises = exercises

	return &weightsLog, nil
}
