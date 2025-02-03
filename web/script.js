document.addEventListener('DOMContentLoaded', function() {
    if (document.getElementById('get-workout')) {
        document.getElementById('get-workout').addEventListener('click', function() {
            fetch('/workout/today')
                .then(response => response.json())
                .then(data => {
                    document.getElementById('workout-display').innerText = `${data.type} - ${data.duration} minutes`;
                })
                .catch(error => {
                    console.error('Error fetching workout of the day:', error);
                });
        });
    }

    if (document.getElementById('cardio-workout-log-form')) {
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
    }

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
        }
    }

    if (document.getElementById('weights-log-form')) {
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
    }

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
    if (document.getElementById('cardio-workouts-table')) {
        fetchLoggedCardioWorkouts();
    }
    if (document.getElementById('weights-workouts-table')) {
        fetchLoggedWeightsWorkouts();
    }
    if (document.getElementById('last-logged-workouts')) {
        fetchLastLoggedWorkouts();
    }

        // Admin panel functionality
        if (document.getElementById('admin-form')) {
            const adminForm = document.getElementById('admin-form');
            const tableNameSelect = document.getElementById('table-name');
            const operationSelect = document.getElementById('operation');
            const fieldsContainer = document.getElementById('fields-container');
            const viewResults = document.getElementById('view-results');
            const emptyTableButton = document.getElementById('empty-table-button');
    
            tableNameSelect.addEventListener('change', function() {
                updateFields();
                fetchTableData();
            });
            operationSelect.addEventListener('change', updateFields);
    
            adminForm.addEventListener('submit', function(event) {
                event.preventDefault();
                const formData = new FormData(adminForm);
                const tableName = formData.get('table-name');
                const operation = formData.get('operation');
                const data = Object.fromEntries(formData.entries());
    
                // Ensure the id field is sent as an integer for delete and update operations
                if (operation === 'delete' || operation === 'update') {
                    data.id = parseInt(data.id, 10);
                }
    
                // Ensure the duration field is sent as an integer
                if (data.duration) {
                    data.duration = parseInt(data.duration, 10);
                }
    
                // Ensure the distance field is sent as a float
                if (data.distance) {
                    data.distance = parseFloat(data.distance);
                }
    
                // Format the date to include seconds and timezone offset
                if (data.date) {
                    const date = new Date(data.date);
                    data.date = date.toISOString();
                }
    
                fetch(`/admin/${operation}/${tableName}`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(data)
                }).then(response => {
                    if (response.ok) {
                        console.log(`${operation} operation on ${tableName} table was successful!`);
                        alert(`${operation} operation on ${tableName} table was successful!`);
                        adminForm.reset();
                        fieldsContainer.innerHTML = '';
                        fetchTableData();
                    } else {
                        console.error(`Failed to ${operation} ${tableName}.`);
                        alert(`Failed to ${operation} ${tableName}.`);
                    }
                }).catch(error => {
                    console.error(`Error performing ${operation} on ${tableName}:`, error);
                    alert(`Error performing ${operation} on ${tableName}: ${error}`);
                });
            });
    
            emptyTableButton.addEventListener('click', function() {
                const tableName = tableNameSelect.value;
                if (confirm(`Are you sure you want to empty the ${tableName} table?`)) {
                    fetch(`/admin/empty/${tableName}`, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        }
                    }).then(response => {
                        if (response.ok) {
                            console.log(`Emptied ${tableName} table successfully!`);
                            alert(`Emptied ${tableName} table successfully!`);
                            fetchTableData();
                        } else {
                            console.error(`Failed to empty ${tableName} table.`);
                            alert(`Failed to empty ${tableName} table.`);
                        }
                    }).catch(error => {
                        console.error(`Error emptying ${tableName} table:`, error);
                        alert(`Error emptying ${tableName} table: ${error}`);
                    });
                }
            });
    
            function updateFields() {
                const tableName = tableNameSelect.value;
                const operation = operationSelect.value;
                fieldsContainer.innerHTML = '';
    
                if (operation === 'add' || operation === 'update') {
                    if (tableName === 'workouts') {
                        fieldsContainer.innerHTML = `
                            ${operation === 'update' ? '<label for="id">ID:</label><input type="number" id="id" name="id" required>' : ''}
                            <label for="type">Type:</label>
                            <input type="text" id="type" name="type" required>
                            <label for="duration">Duration (minutes):</label>
                            <input type="number" id="duration" name="duration" required>
                            <label for="distance">Distance (kms):</label>
                            <input type="number" id="distance" name="distance" required>
                            <label for="date">Date:</label>
                            <input type="datetime-local" id="date" name="date" required>
                        `;
                    } else if (tableName === 'weights_logs') {
                        fieldsContainer.innerHTML = `
                            ${operation === 'update' ? '<label for="id">ID:</label><input type="number" id="id" name="id" required>' : ''}
                            <label for="workout_type">Workout Type:</label>
                            <input type="text" id="workout_type" name="workout_type" required>
                            <label for="date">Date:</label>
                            <input type="datetime-local" id="date" name="date" required>
                        `;
                    } else if (tableName === 'exercises') {
                        fieldsContainer.innerHTML = `
                            ${operation === 'update' ? '<label for="id">ID:</label><input type="number" id="id" name="id" required>' : ''}
                            <label for="weights_log_id">Weights Log ID:</label>
                            <input type="number" id="weights_log_id" name="weights_log_id" required>
                            <label for="name">Name:</label>
                            <input type="text" id="name" name="name" required>
                            <label for="set1">Set 1:</label>
                            <input type="number" id="set1" name="set1" required>
                            <label for="set2">Set 2:</label>
                            <input type="number" id="set2" name="set2" required>
                            <label for="set3">Set 3:</label>
                            <input type="number" id="set3" name="set3" required>
                        `;
                    } else if (tableName === 'wods') {
                        fieldsContainer.innerHTML = `
                            ${operation === 'update' ? '<label for="id">ID:</label><input type="number" id="id" name="id" required>' : ''}
                            <label for="type">Type:</label>
                            <input type="text" id="type" name="type" required>
                            <label for="duration">Duration (minutes):</label>
                            <input type="number" id="duration" name="duration" required>
                            <label for="distance">Distance (kms):</label>
                            <input type="number" id="distance" name="distance" required>
                            <label for="date">Date:</label>
                            <input type="datetime-local" id="date" name="date" required>
                        `;
                    } else if (tableName === 'weight_workouts') {
                        fieldsContainer.innerHTML = `
                            ${operation === 'update' ? '<label for="id">ID:</label><input type="number" id="id" name="id" required>' : ''}
                            <label for="workout_type">Workout Type:</label>
                            <input type="text" id="workout_type" name="workout_type" required>
                            <label for="exercise">Exercise:</label>
                            <input type="text" id="exercise" name="exercise" required>
                        `;
                    }
                } else if (operation === 'delete') {
                    fieldsContainer.innerHTML = `
                        <label for="id">ID:</label>
                        <input type="number" id="id" name="id" required>
                    `;
                }
            }
    
            function fetchTableData() {
                const tableName = tableNameSelect.value;
                fetch(`/admin/view/${tableName}`)
                    .then(response => response.json())
                    .then(data => {
                        console.log(`Fetched data for ${tableName}:`, data);
                        viewResults.innerHTML = generateTableHTML(data);
                    })
                    .catch(error => {
                        console.error(`Error fetching data for ${tableName}:`, error);
                    });
            }
    
            function generateTableHTML(data) {
                if (!Array.isArray(data) || data.length === 0) {
                    return '<p>No data available</p>';
                }
    
                const table = document.createElement('table');
                const thead = document.createElement('thead');
                const tbody = document.createElement('tbody');
    
                // Generate table headers
                const headers = Object.keys(data[0]);
                const headerRow = document.createElement('tr');
                headers.forEach(header => {
                    const th = document.createElement('th');
                    th.textContent = header;
                    headerRow.appendChild(th);
                });
                thead.appendChild(headerRow);
    
                // Generate table rows
                data.forEach(item => {
                    const row = document.createElement('tr');
                    headers.forEach(header => {
                        const td = document.createElement('td');
                        td.textContent = item[header];
                        row.appendChild(td);
                    });
                    tbody.appendChild(row);
                });
    
                table.appendChild(thead);
                table.appendChild(tbody);
                return table.outerHTML;
            }
    
            // Fetch and display the default table data on page load
            fetchTableData();
        }
    });