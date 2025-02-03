document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('get-workout').addEventListener('click', function() {
        fetch('/workout/today')
            .then(response => response.json())
            .then(data => {
                document.getElementById('workout-display').innerText = `Workout of the Day: ${data.type} - ${data.duration} minutes`;
            })
            .catch(error => {
                console.error('Error fetching workout of the day:', error);
            });
    });

    document.getElementById('cardio-workout-log-form').addEventListener('submit', function(event) {
        event.preventDefault();
        const formData = new FormData(event.target);
        const workout = {
            type: formData.get('exercise-type'),
            duration: parseInt(formData.get('duration'), 10),
            distance: parseFloat(formData.get('distance'))
        };
        fetch('/workout/log/cardio', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(workout)
        }).then(response => {
            if (response.ok) {
                console.log('Cardio workout logged successfully!');
                fetchLoggedCardioWorkouts();
                resetForm('cardio-workout-log-form');
            } else {
                console.error('Failed to log cardio workout.');
            }
        }).catch(error => {
            console.error('Error logging cardio workout:', error);
        });
    });

    document.getElementById('workout-type').addEventListener('change', function(event) {
        const workoutType = event.target.value;
        console.log(`Selected workout type: ${workoutType}`);
        fetch(`/workout/weight-workouts?type=${workoutType}`)
            .then(response => response.json())
            .then(exercises => {
                console.log('Fetched exercises:', exercises);
                const exercisesContainer = document.getElementById('exercises');
                exercisesContainer.innerHTML = '';
                if (exercises.length === 0) {
                    exercisesContainer.innerHTML = '<p>No exercises found for this workout type.</p>';
                } else {
                    const table = document.createElement('table');
                    table.innerHTML = `
                        <thead>
                            <tr>
                                <th>Exercise</th>
                                <th>Set 1 Weight</th>
                                <th>Set 2 Weight</th>
                                <th>Set 3 Weight</th>
                            </tr>
                        </thead>
                        <tbody></tbody>
                    `;
                    const tbody = table.querySelector('tbody');
                    exercises.forEach(exercise => {
                        const row = document.createElement('tr');
                        row.innerHTML = `
                            <td>${exercise.exercise}</td>
                            <td><input type="number" name="${exercise.exercise}-set1" value="0"></td>
                            <td><input type="number" name="${exercise.exercise}-set2" value="0"></td>
                            <td><input type="number" name="${exercise.exercise}-set3" value="0"></td>
                        `;
                        tbody.appendChild(row);
                    });
                    exercisesContainer.appendChild(table);
                }
            })
            .catch(error => {
                console.error('Error fetching exercises:', error);
            });
    });

    document.getElementById('weights-log-form').addEventListener('submit', function(event) {
        event.preventDefault();
        const formData = new FormData(event.target);
        const workoutType = formData.get('workout-type');
        const exercises = Array.from(document.querySelectorAll('#exercises tbody tr')).map(row => {
            const exerciseName = row.querySelector('td').innerText;
            return {
                name: exerciseName,
                set1: parseInt(formData.get(`${exerciseName}-set1`) || 0, 10),
                set2: parseInt(formData.get(`${exerciseName}-set2`) || 0, 10),
                set3: parseInt(formData.get(`${exerciseName}-set3`) || 0, 10)
            };
        });
        const weightsLog = {
            workout_type: workoutType,
            exercises: exercises,
            date: new Date().toISOString() // Set the current date and time
        };
        console.log('Logging weights workout:', weightsLog);
        fetch('/workout/log/weights', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(weightsLog)
        }).then(response => {
            if (response.ok) {
                console.log('Weights workout logged successfully!');
                fetchLoggedWeightsWorkouts();
                fetchLastLoggedWorkouts();
                resetForm('weights-log-form');
                document.getElementById('exercises').innerHTML = ''; // Clear exercises table
            } else {
                console.error('Failed to log weights workout.');
            }
        }).catch(error => {
            console.error('Error logging weights workout:', error);
        });
    });

    function fetchLoggedCardioWorkouts() {
        console.log('Fetching logged cardio workouts...');
        fetch('/workout/logs/cardio')
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                console.log('Fetched logged cardio workouts:', data);
                const tableBody = document.getElementById('cardio-workouts-table').querySelector('tbody');
                tableBody.innerHTML = '';
                data.forEach(workout => {
                    const row = document.createElement('tr');
                    row.innerHTML = `
                        <td>${workout.type}</td>
                        <td>${workout.duration}</td>
                        <td>${workout.distance}</td>
                        <td>${new Date(workout.date).toLocaleString()}</td>
                    `;
                    tableBody.appendChild(row);
                });
            })
            .catch(error => {
                console.error('Error fetching logged cardio workouts:', error);
            });
    }

    function fetchLoggedWeightsWorkouts() {
        console.log('Fetching logged weights workouts...');
        fetch('/workout/logs/weights')
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                console.log('Fetched logged weights workouts:', data);
                const tableBody = document.getElementById('weights-workouts-table').querySelector('tbody');
                tableBody.innerHTML = '';
                data.forEach(workout => {
                    const row = document.createElement('tr');
                    row.innerHTML = `
                        <td>${workout.workout_type}</td>
                        <td>${new Date(workout.date).toLocaleString()}</td>
                    `;
                    tableBody.appendChild(row);
                });
            })
            .catch(error => {
                console.error('Error fetching logged weights workouts:', error);
            });
    }

    function fetchLastLoggedWorkouts() {
        fetchLastLoggedWorkout('push', 'push-workout-table');
        fetchLastLoggedWorkout('pull', 'pull-workout-table');
        fetchLastLoggedWorkout('legs', 'legs-workout-table');
    }
    
    function fetchLastLoggedWorkout(workoutType, tableId) {
        console.log(`Fetching last logged ${workoutType} workout...`);
        fetch(`/workout/last/weights?type=${workoutType}`)
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(weightsLog => {
                console.log(`Fetched last logged ${workoutType} workout:`, weightsLog);
                const tableBody = document.getElementById(tableId).querySelector('tbody');
                tableBody.innerHTML = '';
                if (weightsLog) {
                    if (weightsLog.exercises && weightsLog.exercises.length > 0) {
                        weightsLog.exercises.forEach(exercise => {
                            const row = document.createElement('tr');
                            row.innerHTML = `
                                <td>${exercise.name}</td>
                                <td>${exercise.set1}</td>
                                <td>${exercise.set2}</td>
                                <td>${exercise.set3}</td>
                                <td>${new Date(weightsLog.date).toLocaleString()}</td>
                            `;
                            tableBody.appendChild(row);
                        });
                    } else {
                        console.log(`No exercises found for ${workoutType} workout.`);
                        tableBody.innerHTML = '<tr><td colspan="5">No exercises found</td></tr>';
                    }
                } else {
                    tableBody.innerHTML = '<tr><td colspan="5">No workout found</td></tr>';
                }
            })
            .catch(error => {
                console.error(`Error fetching last logged ${workoutType} workout:`, error);
            });
    }

    function resetForm(formId) {
        document.getElementById(formId).reset();
    }

    // Fetch logged workouts and last logged workouts on page load
    fetchLoggedCardioWorkouts();
    fetchLoggedWeightsWorkouts();
    fetchLastLoggedWorkouts();
});