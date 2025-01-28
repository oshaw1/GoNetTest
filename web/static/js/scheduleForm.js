document.addEventListener('DOMContentLoaded', function() {
    // Handle form submission
    const form = document.querySelector('.task-form');
    if (form) {
        form.addEventListener('htmx:beforeRequest', function(event) {
            event.preventDefault();
            
            const formData = new FormData(form);
            
            const requestData = {
                name: formData.get('name'),
                datetime: new Date(formData.get('datetime')).toISOString(),
                recurring: formData.get('recurring') === 'on',
                interval: formData.get('interval') || 'daily',
                active: formData.get('active') === 'on'
            };

            const taskType = formData.get('task_type');
            if (taskType === 'test') {
                requestData.test_type = formData.get('test_type');
            } else {
                requestData.chart_type = formData.get('chart_type');
                if (formData.get('historic') === 'on' && formData.get('recent_days')) {
                    requestData.recent_days = parseInt(formData.get('recent_days'));
                }
            }

            fetch('/schedule/create', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestData)
            })
            .then(response => response.json())
            .then(() => {
                // Trigger refresh of schedule list
                htmx.trigger('#schedule-list', 'taskChanged');
                closeModal();
            })
            .catch(error => console.error('Error:', error));
        });
    }

    const recurringCheckbox = document.getElementById('recurring');
    if (recurringCheckbox) {
        recurringCheckbox.addEventListener('change', function() {
            const intervalGroup = document.getElementById('interval-group');
            intervalGroup.style.display = this.checked ? 'block' : 'none';
        });
    }

    // Task type radio buttons
    const taskTypeRadios = document.querySelectorAll('input[name="task_type"]');
    taskTypeRadios.forEach(radio => {
        radio.addEventListener('change', updateFieldVisibility);
    });

    // Historic checkbox
    const historicCheckbox = document.getElementById('historic');
    if (historicCheckbox) {
        historicCheckbox.addEventListener('change', function() {
            const recentDaysField = document.querySelector('.recent-days-field');
            recentDaysField.style.display = this.checked ? 'block' : 'none';
        });
    }

    updateFieldVisibility();
});

function updateFieldVisibility() {
    const selectedType = document.querySelector('input[name="task_type"]:checked').value;
    const testFields = document.querySelectorAll('.test-field');
    const chartFields = document.querySelectorAll('.chart-field');
    
    if (selectedType === 'test') {
        testFields.forEach(field => field.style.display = 'block');
        chartFields.forEach(field => field.style.display = 'none');
        document.querySelector('.recent-days-field').style.display = 'none';
    } else {
        testFields.forEach(field => field.style.display = 'none');
        chartFields.forEach(field => field.style.display = 'block');
        // Show recent days only if historic is checked
        const historicChecked = document.getElementById('historic').checked;
        document.querySelector('.recent-days-field').style.display = 
            historicChecked ? 'block' : 'none';
    }
}

function showModal() {
    document.getElementById('task-modal').style.display = 'block';
}

function closeModal() {
    document.getElementById('task-modal').style.display = 'none';
    // Reset form
    document.querySelector('.task-form').reset();
    updateFieldVisibility();
}

window.onclick = function(event) {
    const modal = document.getElementById('task-modal');
    if (event.target == modal) {
        closeModal();
    }
}