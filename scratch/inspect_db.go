package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	connStr := os.Getenv("DATABASE_URL")
	fmt.Printf("Connecting with: %s\n", connStr)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Inspecting 'exercises' table schema:")
	rows, err := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = 'exercises'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name, dataType string
		rows.Scan(&name, &dataType)
		fmt.Printf("Column: %s, Type: %s\n", name, dataType)
	}

	fmt.Println("\nChecking if 'exercise_sets' table exists:")
	var exists bool
	err = db.Get(&exists, "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'exercise_sets')")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("exercise_sets exists: %v\n", exists)

	fmt.Println("\nCounting records:")
	var weightLogsCount, exercisesCount, setsCount, workoutsCount, wodsCount int
	db.Get(&weightLogsCount, "SELECT COUNT(*) FROM weights_logs")
	db.Get(&exercisesCount, "SELECT COUNT(*) FROM exercises")
	db.Get(&workoutsCount, "SELECT COUNT(*) FROM workouts")
	db.Get(&wodsCount, "SELECT COUNT(*) FROM wods")
	if exists {
		db.Get(&setsCount, "SELECT COUNT(*) FROM exercise_sets")
	}

	fmt.Printf("weights_logs: %d\n", weightLogsCount)
	fmt.Printf("exercises: %d\n", exercisesCount)
	fmt.Printf("exercise_sets: %d\n", setsCount)
	fmt.Printf("workouts (cardio): %d\n", workoutsCount)
	fmt.Printf("wods: %d\n", wodsCount)
}
