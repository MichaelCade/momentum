document.addEventListener('DOMContentLoaded', function() {
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
                    // Fetch updated data for history page
                    fetchLoggedCardioWorkouts();
                    fetchLoggedWeightsWorkouts();
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
                        // Fetch updated data for history page
                        fetchLoggedCardioWorkouts();
                        fetchLoggedWeightsWorkouts();
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
                    if (header === 'duration') {
                        const durationMinutes = Math.floor(item[header] / 60);
                        const durationSeconds = item[header] % 60;
                        td.textContent = `${durationMinutes}m ${durationSeconds}s`;
                    } else {
                        td.textContent = item[header];
                    }
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