document.addEventListener('DOMContentLoaded', function() {
    fetchLoggedCardioWorkouts();
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
            visualizeCardioWorkouts(data);
        })
        .catch(error => {
            console.error('Error fetching logged cardio workouts:', error);
        });
}

function visualizeCardioWorkouts(data) {
    const dates = data.map(workout => new Date(workout.date).toLocaleDateString());
    const distances = data.map(workout => workout.distance);
    const types = data.map(workout => workout.type);

    // Cardio Trend Chart
    const ctxTrend = document.getElementById('cardioTrendChart').getContext('2d');
    new Chart(ctxTrend, {
        type: 'line',
        data: {
            labels: dates,
            datasets: [{
                label: 'Distance (kms)',
                data: distances,
                borderColor: 'rgba(75, 192, 192, 1)',
                backgroundColor: 'rgba(75, 192, 192, 0.2)',
                fill: true,
            }]
        },
        options: {
            scales: {
                x: {
                    title: {
                        display: true,
                        text: 'Date'
                    }
                },
                y: {
                    title: {
                        display: true,
                        text: 'Distance (kms)'
                    }
                }
            }
        }
    });

    // Total Distance Chart
    const totalDistanceByType = types.reduce((acc, type, index) => {
        acc[type] = (acc[type] || 0) + distances[index];
        return acc;
    }, {});

    const ctxTotal = document.getElementById('totalDistanceChart').getContext('2d');
    new Chart(ctxTotal, {
        type: 'bar',
        data: {
            labels: Object.keys(totalDistanceByType),
            datasets: [{
                label: 'Total Distance (kms)',
                data: Object.values(totalDistanceByType),
                backgroundColor: 'rgba(153, 102, 255, 0.2)',
                borderColor: 'rgba(153, 102, 255, 1)',
                borderWidth: 1
            }]
        },
        options: {
            scales: {
                x: {
                    title: {
                        display: true,
                        text: 'Exercise Type'
                    }
                },
                y: {
                    title: {
                        display: true,
                        text: 'Total Distance (kms)'
                    }
                }
            }
        }
    });

    // Combined Total Distance Pie Chart
    const ctxCombinedTotal = document.getElementById('combinedTotalDistanceChart').getContext('2d');
    new Chart(ctxCombinedTotal, {
        type: 'pie',
        data: {
            labels: Object.keys(totalDistanceByType),
            datasets: [{
                label: 'Distance Covered by Type (kms)',
                data: Object.values(totalDistanceByType),
                backgroundColor: [
                    'rgba(255, 99, 132, 0.2)',
                    'rgba(54, 162, 235, 0.2)',
                    'rgba(255, 206, 86, 0.2)',
                    'rgba(75, 192, 192, 0.2)',
                    'rgba(153, 102, 255, 0.2)',
                    'rgba(255, 159, 64, 0.2)'
                ],
                borderColor: [
                    'rgba(255, 99, 132, 1)',
                    'rgba(54, 162, 235, 1)',
                    'rgba(255, 206, 86, 1)',
                    'rgba(75, 192, 192, 1)',
                    'rgba(153, 102, 255, 1)',
                    'rgba(255, 159, 64, 1)'
                ],
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
            plugins: {
                legend: {
                    position: 'top',
                },
                title: {
                    display: true,
                    text: 'Total Distance Covered by Type'
                }
            }
        }
    });
}