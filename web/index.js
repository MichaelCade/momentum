document.addEventListener('DOMContentLoaded', function() {
    const getWorkoutButton = document.getElementById('get-workout');
    const workoutDisplay = document.getElementById('workout-display');

    getWorkoutButton.addEventListener('click', function() {
        fetch('/workout/today')
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(workout => {
                workoutDisplay.innerHTML = `
                    <p>Type: ${workout.type}</p>
                    <p>Duration: ${workout.duration} minutes</p>
                    <p>Distance: ${workout.distance} kms</p>
                `;
            })
            .catch(error => {
                console.error('Error fetching workout of the day:', error);
                workoutDisplay.innerHTML = '<p>Error fetching workout of the day. Please try again later.</p>';
            });
    });

     // Fetch logged cardio workouts and visualize them on page load
     fetchLoggedCardioWorkouts();
});