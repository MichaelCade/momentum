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
	Duration float64   `json:"duration"` // Duration in seconds
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
	err := database.DB.Select(&workouts, "SELECT * FROM workouts WHERE LOWER(type) IN ('run', 'bike', 'row', 'walk', 'crosstrainer') ORDER BY date DESC")
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

// Admin Section - Add, Update, Delete

// AddWorkout adds a new workout to the database
func AddWorkout(workout Workout) error {
	_, err := database.DB.NamedExec(`INSERT INTO workouts (type, duration, distance, date) VALUES (:type, :duration, :distance, :date)`, &workout)
	return err
}

// UpdateWorkout updates an existing workout in the database
func UpdateWorkout(workout Workout) error {
	_, err := database.DB.NamedExec(`UPDATE workouts SET type=:type, duration=:duration, distance=:distance, date=:date WHERE id=:id`, &workout)
	return err
}

// DeleteWorkout deletes a workout from the database
func DeleteWorkout(id int) error {
	_, err := database.DB.Exec(`DELETE FROM workouts WHERE id=$1`, id)
	return err
}

// AddWeightsLog adds a new weights log to the database
func AddWeightsLog(weightsLog WeightsLog) error {
	weightsLog.Date = time.Now() // Set the current time and date
	_, err := database.DB.NamedExec(`INSERT INTO weights_logs (workout_type, date) VALUES (:workout_type, :date)`, &weightsLog)
	return err
}

// UpdateWeightsLog updates an existing weights log in the database
func UpdateWeightsLog(weightsLog WeightsLog) error {
	_, err := database.DB.NamedExec(`UPDATE weights_logs SET workout_type=:workout_type, date=:date WHERE id=:id`, &weightsLog)
	return err
}

func DeleteWeightsLog(id int) error {
	// Delete related records in the exercises table
	_, err := database.DB.Exec("DELETE FROM exercises WHERE weights_log_id = $1", id)
	if err != nil {
		log.Printf("Error deleting related exercises: %v", err)
		return err
	}

	// Delete the record in the weights_logs table
	_, err = database.DB.Exec("DELETE FROM weights_logs WHERE id = $1", id)
	if err != nil {
		log.Printf("Error deleting record from weights_logs table: %v", err)
		return err
	}
	return nil
}

// AddExercise adds a new exercise to the database
func AddExercise(exercise Exercise) error {
	_, err := database.DB.NamedExec(`INSERT INTO exercises (weights_log_id, name, set1, set2, set3) VALUES (:weights_log_id, :name, :set1, :set2, :set3)`, &exercise)
	return err
}

// UpdateExercise updates an existing exercise in the database
func UpdateExercise(exercise Exercise) error {
	_, err := database.DB.NamedExec(`UPDATE exercises SET weights_log_id=:weights_log_id, name=:name, set1=:set1, set2=:set2, set3=:set3 WHERE id=:id`, &exercise)
	return err
}

// DeleteExercise deletes an exercise from the database
func DeleteExercise(id int) error {
	_, err := database.DB.Exec(`DELETE FROM exercises WHERE id=$1`, id)
	return err
}

// AddWOD adds a new WOD to the database
func AddWOD(wod WOD) error {
	_, err := database.DB.NamedExec(`INSERT INTO wods (type, duration, distance, date) VALUES (:type, :duration, :distance, :date)`, &wod)
	return err
}

// UpdateWOD updates an existing WOD in the database
func UpdateWOD(wod WOD) error {
	_, err := database.DB.NamedExec(`UPDATE wods SET type=:type, duration=:duration, distance=:distance, date=:date WHERE id=:id`, &wod)
	return err
}

// DeleteWOD deletes a WOD from the database
func DeleteWOD(id int) error {
	_, err := database.DB.Exec(`DELETE FROM wods WHERE id=$1`, id)
	return err
}

// AddWeightWorkout adds a new weight workout to the database
func AddWeightWorkout(weightWorkout WeightWorkout) error {
	_, err := database.DB.NamedExec(`INSERT INTO weight_workouts (workout_type, exercise) VALUES (:workout_type, :exercise)`, &weightWorkout)
	return err
}

// UpdateWeightWorkout updates an existing weight workout in the database
func UpdateWeightWorkout(weightWorkout WeightWorkout) error {
	_, err := database.DB.NamedExec(`UPDATE weight_workouts SET workout_type=:workout_type, exercise=:exercise WHERE id=:id`, &weightWorkout)
	return err
}

// DeleteWeightWorkout deletes a weight workout from the database
func DeleteWeightWorkout(id int) error {
	_, err := database.DB.Exec(`DELETE FROM weight_workouts WHERE id=$1`, id)
	return err
}

// ViewWorkouts retrieves all workouts from the database
func ViewWorkouts() ([]Workout, error) {
	var workouts []Workout
	err := database.DB.Select(&workouts, "SELECT * FROM workouts")
	if err != nil {
		log.Printf("Error viewing workouts: %v", err)
		return nil, err
	}
	return workouts, nil
}

// ViewWeightsLogs retrieves all weights logs from the database
func ViewWeightsLogs() ([]WeightsLog, error) {
	var weightsLogs []WeightsLog
	err := database.DB.Select(&weightsLogs, "SELECT * FROM weights_logs")
	if err != nil {
		log.Printf("Error viewing weights logs: %v", err)
		return nil, err
	}
	return weightsLogs, nil
}

// ViewExercises retrieves all exercises from the database
func ViewExercises() ([]Exercise, error) {
	var exercises []Exercise
	err := database.DB.Select(&exercises, "SELECT * FROM exercises")
	if err != nil {
		log.Printf("Error viewing exercises: %v", err)
		return nil, err
	}
	return exercises, nil
}

// ViewWODs retrieves all WODs from the database
func ViewWODs() ([]WOD, error) {
	var wods []WOD
	err := database.DB.Select(&wods, "SELECT * FROM wods")
	if err != nil {
		log.Printf("Error viewing WODs: %v", err)
		return nil, err
	}
	return wods, nil
}

// ViewWeightWorkouts retrieves all weight workouts from the database
func ViewWeightWorkouts() ([]WeightWorkout, error) {
	var weightWorkouts []WeightWorkout
	err := database.DB.Select(&weightWorkouts, "SELECT * FROM weight_workouts")
	if err != nil {
		log.Printf("Error viewing weight workouts: %v", err)
		return nil, err
	}
	return weightWorkouts, nil
}

func EmptyWorkouts() error {
	_, err := database.DB.Exec("DELETE FROM workouts")
	if err != nil {
		log.Printf("Error emptying workouts table: %v", err)
		return err
	}
	return nil
}

func EmptyWeightsLogs() error {
	// Delete related records in the exercises table
	_, err := database.DB.Exec("DELETE FROM exercises WHERE weights_log_id IN (SELECT id FROM weights_logs)")
	if err != nil {
		log.Printf("Error deleting related exercises: %v", err)
		return err
	}

	// Delete records in the weights_logs table
	_, err = database.DB.Exec("DELETE FROM weights_logs")
	if err != nil {
		log.Printf("Error emptying weights_logs table: %v", err)
		return err
	}
	return nil
}

func EmptyExercises() error {
	_, err := database.DB.Exec("DELETE FROM exercises")
	if err != nil {
		log.Printf("Error emptying exercises table: %v", err)
		return err
	}
	return nil
}

func EmptyWODs() error {
	_, err := database.DB.Exec("DELETE FROM wods")
	if err != nil {
		log.Printf("Error emptying wods table: %v", err)
		return err
	}
	return nil
}

func EmptyWeightWorkouts() error {
	_, err := database.DB.Exec("DELETE FROM weight_workouts")
	if err != nil {
		log.Printf("Error emptying weight_workouts table: %v", err)
		return err
	}
	return nil
}
