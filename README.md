# Momentum

## Overview
This project is a web application built with Go that serves as a workout generator and logging system. It allows users to retrieve a workout of the day, log their workouts, and track various types of weight training plans.

## Project Structure
```
momentum
├── cmd
│   └── main.go          # Entry point of the application
├── internal
│   ├── handlers
│   │   └── workout.go   # HTTP request handlers for workouts
│   ├── models
│   │   └── workout.go    # Workout data model
│   ├── routes
│   │   └── routes.go    # Application routes
│   └── database
│       └── database.go  # Database connection and queries
├── web
│   └── index.html       # Frontend HTML file
|   └── script.js
|   └── style.css   
├── go.mod               # Module definition and dependencies
└── README.md            # Project documentation
```

## Features
- **Random Workout Generator**: Fetches a workout of the day from the database.
- **Workout Logging**: Allows users to log workouts with details such as exercise type, duration, and distance.
- **Weights Tracking**: Supports various weight training plans including Push, Pull, and Legs routines.

## Setup Instructions
1. Clone the repository:
   ```
   git clone <repository-url>
   cd momentum
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Run the application:
   ```
   go run cmd/main.go
   ```

4. Open your browser and navigate to `http://localhost:8080` to access the web app.

## Usage
- Access the workout of the day via the main page.
- Log your workouts through the provided interface.
- Track your weight training plans and view your workout history.

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

## License
This project is licensed under the MIT License.