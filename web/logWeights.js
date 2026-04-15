document.addEventListener('DOMContentLoaded', function() {
    if (document.getElementById('log-weights-workout')) {
        // Fetch and display exercises for the default selected workout type
        const defaultWorkoutType = document.getElementById('workout-type').value;
        fetchExercises(defaultWorkoutType);

        document.getElementById('workout-type').addEventListener('change', function(event) {
            const workoutType = event.target.value;
            fetchExercises(workoutType);
        });

        function fetchExercises(workoutType) {
            console.log(`Fetching exercises for workout type: ${workoutType}`);
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
                                    <th>Set 1 (kg x reps)</th>
                                    <th>Set 2 (kg x reps)</th>
                                    <th>Set 3 (kg x reps)</th>
                                </tr>
                            </thead>
                            <tbody></tbody>
                        `;
                        const tbody = table.querySelector('tbody');
                        exercises.forEach(exercise => {
                            const row = document.createElement('tr');
                            row.innerHTML = `
                                <td>${exercise.exercise}</td>
                                <td>
                                    <input type="number" name="${exercise.exercise}-set1-weight" placeholder="kg" style="width: 60px;"> x 
                                    <input type="number" name="${exercise.exercise}-set1-reps" placeholder="reps" style="width: 60px;">
                                </td>
                                <td>
                                    <input type="number" name="${exercise.exercise}-set2-weight" placeholder="kg" style="width: 60px;"> x 
                                    <input type="number" name="${exercise.exercise}-set2-reps" placeholder="reps" style="width: 60px;">
                                </td>
                                <td>
                                    <input type="number" name="${exercise.exercise}-set3-weight" placeholder="kg" style="width: 60px;"> x 
                                    <input type="number" name="${exercise.exercise}-set3-reps" placeholder="reps" style="width: 60px;">
                                </td>
                            `;
                            tbody.appendChild(row);
                        });
                        exercisesContainer.appendChild(table);
                    }
                })
                .catch(error => {
                    console.error('Error fetching exercises:', error);
                });
        }
    }

    if (document.getElementById('weights-log-form')) {
        document.getElementById('weights-log-form').addEventListener('submit', function(event) {
            event.preventDefault();
            const formData = new FormData(event.target);
            const workoutType = formData.get('workout-type');
            const exercises = Array.from(document.querySelectorAll('#exercises tbody tr')).map(row => {
                const exerciseName = row.querySelector('td').innerText;
                const sets = [];
                for (let i = 1; i <= 3; i++) {
                    const weightVal = formData.get(`${exerciseName}-set${i}-weight`);
                    const repsVal = formData.get(`${exerciseName}-set${i}-reps`);
                    if (weightVal !== "" || repsVal !== "") {
                        sets.push({
                            weight: parseFloat(weightVal || 0),
                            reps: parseInt(repsVal || 0, 10)
                        });
                    }
                }
                return {
                    name: exerciseName,
                    sets: sets
                };
            }).filter(ex => ex.sets.length > 0);

            const weightsLog = {
                workout_type: workoutType,
                exercises: exercises,
                date: new Date().toISOString()
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
    }
});