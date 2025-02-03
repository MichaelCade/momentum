package database

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB() {
	var err error
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatalln("DATABASE_URL environment variable is not set")
	}
	log.Println("Connecting to database with connection string:", connStr)
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalln("Error connecting to the database:", err)
	}

	createSchema()
}

// createSchema creates the necessary database schema if it doesn't exist
func createSchema() {
	schema := `
    CREATE TABLE IF NOT EXISTS workouts (
        id SERIAL PRIMARY KEY,
        type VARCHAR(50),
        duration INT,
        distance FLOAT,
        date TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS workout_logs (
        id SERIAL PRIMARY KEY,
        exercise VARCHAR(50),
        time FLOAT,
        distance FLOAT,
        weight FLOAT,
        reps INT,
        sets INT,
        date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS weights_logs (
        id SERIAL PRIMARY KEY,
        workout_type VARCHAR(50),
        date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS exercises (
        id SERIAL PRIMARY KEY,
        weights_log_id INT REFERENCES weights_logs(id),
        name VARCHAR(50),
        set1 INT,
        set2 INT,
        set3 INT
    );

    CREATE TABLE IF NOT EXISTS wods (
        id SERIAL PRIMARY KEY,
        type VARCHAR(50),
        duration INT,
        distance FLOAT,
        date TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS weight_workouts (
        id SERIAL PRIMARY KEY,
        workout_type VARCHAR(50),
        exercise VARCHAR(50)
    );
    `
	DB.MustExec(schema)

	// Insert initial data if the wods table is empty
	var count int
	err := DB.Get(&count, "SELECT COUNT(*) FROM wods")
	if err != nil {
		log.Fatalln("Error checking wods table:", err)
	}

	if count == 0 {
		initialData := `
        INSERT INTO wods (type, duration, distance, date) VALUES
        ('walk', 60, 5.0, NOW()),
        ('run', 30, 5.0, NOW()),
        ('run', 30, 5.0, NOW()),
        ('crosstrainer', 30, 5.0, NOW()),
        ('row', 30, 5.0, NOW()),
        ('row', 30, 5.0, NOW()),
        ('bike', 60, 20.0, NOW());
        `
		DB.MustExec(initialData)
	}

	// Insert initial data for weight workouts if the table is empty
	err = DB.Get(&count, "SELECT COUNT(*) FROM weight_workouts")
	if err != nil {
		log.Fatalln("Error checking weight_workouts table:", err)
	}

	if count == 0 {
		initialData := `
        INSERT INTO weight_workouts (workout_type, exercise) VALUES
        ('push', 'Flat Dumbbells'),
        ('push', 'Flat Flys'),
        ('push', 'Seated Dumbbell front raises'),
        ('push', 'Seated Dumbbell side raises'),
        ('push', 'Seated Dumbbell shoulder press'),
        ('push', 'Tricep Pushdowns'),
        ('push', 'Incline Smith'),
        ('push', 'Close Grip Incline Smith'),
        ('push', 'Overhead Rope (Cables)'),
        ('push', 'Assisted Dips/Dip machine'),
        ('pull', 'Deadlifts'),
        ('pull', 'Bent Over Rows (Underhand)'),
        ('pull', 'Shrugs (Barbell or dumbbell)'),
        ('pull', 'Lat Pulldown'),
        ('pull', 'Upright Rows (Barbell or Rope)'),
        ('pull', 'Rear Delt Raises (Dumbbell)'),
        ('pull', 'Single Preacher Dumbbell Curls'),
        ('pull', 'EZ Bar Standing Curls'),
        ('pull', 'Double Dumbbell Hammer Curls'),
        ('legs', 'Barbell Squat'),
        ('legs', 'Straight leg deadlifts'),
        ('legs', 'Front squat (added)'),
        ('legs', 'Leg Press'),
        ('legs', 'Calf Raises on Leg Press'),
        ('legs', 'Leg Extensions'),
        ('legs', 'Hamstring curls (Machine)'),
        ('legs', 'Dumbbell lunges'),
        ('legs', 'Ab/Crunch Machine'),
        ('legs', 'Captains Chair Leg or Knee Raises');
        `
		DB.MustExec(initialData)
	}

	// Insert sample data into the workouts table if it is empty
	err = DB.Get(&count, "SELECT COUNT(*) FROM workouts")
	if err != nil {
		log.Fatalln("Error checking workouts table:", err)
	}

	if count == 0 {
		sampleData := `
        INSERT INTO workouts (type, duration, distance, date) VALUES
        ('run', 30, 5.0, NOW()),
        ('bike', 60, 20.0, NOW()),
        ('row', 30, 5.0, NOW());
        `
		DB.MustExec(sampleData)
	}
}

// GetDB returns the database connection
func GetDB() *sqlx.DB {
	return DB
}

// CloseDB closes the database connection
func CloseDB() {
	if err := DB.Close(); err != nil {
		log.Fatalf("Error closing database: %v", err)
	}
}
