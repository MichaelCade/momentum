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
                const durationMinutes = Math.floor(workout.duration / 60);
                const durationSeconds = workout.duration % 60;
                const row = document.createElement('tr');
                row.innerHTML = `
                    <td>${workout.type}</td>
                    <td>${durationMinutes}m ${durationSeconds}s</td>
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