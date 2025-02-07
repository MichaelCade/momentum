document.addEventListener('DOMContentLoaded', function() {
    if (document.getElementById('cardio-workout-log-form')) {
        document.getElementById('cardio-workout-log-form').addEventListener('submit', function(event) {
            event.preventDefault();
            const formData = new FormData(event.target);
            const minutes = parseInt(formData.get('duration-minutes'), 10);
            const seconds = parseInt(formData.get('duration-seconds'), 10);
            const duration = minutes * 60 + seconds; // Convert to seconds
            const workout = {
                type: formData.get('exercise-type'),
                duration: duration,
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
});