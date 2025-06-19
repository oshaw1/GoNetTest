console.log('scheduleForm.js loaded');

let isEditMode = false;
let currentTaskId = null;

// Make functions globally available
window.showModal = function() {
    console.log('showModal called');
    isEditMode = false;
    currentTaskId = null;
    
    // Set up form for creation
    const form = document.querySelector('.task-form');
    if (form) {
        form.setAttribute('hx-post', '/schedule/create');
        form.setAttribute('hx-target', '#schedule-list');
        form.setAttribute('hx-trigger', 'submit');
        form.removeAttribute('hx-put');
        
        // Reset form
        form.reset();
        updateFieldVisibility();
    }
    
    // Update modal title
    const modalTitle = document.querySelector('#task-modal h3');
    if (modalTitle) {
        modalTitle.textContent = 'Add New Task';
    }
    
    document.getElementById('task-modal').style.display = 'block';
};

window.editTask = function(taskId) {
    console.log('editTask called with ID:', taskId);
    
    const scheduleItem = document.getElementById('task-' + taskId);
    if (!scheduleItem) {
        console.error('Could not find schedule item for ID:', taskId);
        alert('Error: Could not find task data');
        return;
    }
    
    editTaskFromElement(scheduleItem);
};

window.editTaskFromElement = function(scheduleElement) {
    console.log('editTaskFromElement called', scheduleElement);
    
    // Extract data from data attributes
    const taskData = {
        id: scheduleElement.dataset.taskId,
        name: scheduleElement.dataset.taskName,
        test_type: scheduleElement.dataset.testType,
        chart_type: scheduleElement.dataset.chartType,
        recent_days: parseInt(scheduleElement.dataset.recentDays) || 0,
        datetime: scheduleElement.dataset.datetime,
        recurring: scheduleElement.dataset.recurring === 'true',
        interval: scheduleElement.dataset.interval,
        active: scheduleElement.dataset.active === 'true'
    };
    
    console.log('Extracted task data:', taskData);
    
    isEditMode = true;
    currentTaskId = taskData.id;
    
    // Update modal title
    const modalTitle = document.querySelector('#task-modal h3');
    if (modalTitle) {
        modalTitle.textContent = 'Edit Task';
    }
    
    // Populate form with the extracted data
    populateForm(taskData);
    
    // Set up form for editing
    const form = document.querySelector('.task-form');
    if (form) {
        form.setAttribute('hx-put', '/schedule/edit?id=' + taskData.id);
        form.setAttribute('hx-target', '#schedule-list');
        form.setAttribute('hx-trigger', 'submit');
        form.removeAttribute('hx-post');
    }
    
    document.getElementById('task-modal').style.display = 'block';
};

window.closeModal = function() {
    console.log('closeModal called');
    document.getElementById('task-modal').style.display = 'none';
    // Reset form
    const form = document.querySelector('.task-form');
    if (form) {
        form.reset();
        updateFieldVisibility();
    }
    isEditMode = false;
    currentTaskId = null;
};

function populateForm(task) {
    console.log('Populating form with task:', task);
    
    const nameInput = document.getElementById('name');
    if (nameInput) {
        nameInput.value = task.name || '';
    }
    
    // Set task type radio button
    if (task.test_type) {
        const testRadio = document.querySelector('input[name="task_type"][value="test"]');
        if (testRadio) {
            testRadio.checked = true;
        }
        const testTypeSelect = document.getElementById('test_type');
        if (testTypeSelect) {
            testTypeSelect.value = task.test_type;
        }
    } else if (task.chart_type) {
        const chartRadio = document.querySelector('input[name="task_type"][value="chart"]');
        if (chartRadio) {
            chartRadio.checked = true;
        }
        const chartTypeSelect = document.getElementById('chart_type');
        if (chartTypeSelect) {
            chartTypeSelect.value = task.chart_type;
        }
        
        // Handle historic checkbox and recent days
        if (task.recent_days && task.recent_days > 0) {
            const historicCheckbox = document.getElementById('historic');
            if (historicCheckbox) {
                historicCheckbox.checked = true;
            }
            const recentDaysInput = document.getElementById('recent_days');
            if (recentDaysInput) {
                recentDaysInput.value = task.recent_days;
            }
        }
    }
    
    // Set datetime for datetime-local input
    if (task.datetime) {
        const datetimeInput = document.getElementById('datetime');
        if (datetimeInput) {
            datetimeInput.value = task.datetime;
        }
    }
    
    const recurringCheckbox = document.getElementById('recurring');
    if (recurringCheckbox) {
        recurringCheckbox.checked = task.recurring || false;
    }
    
    if (task.interval) {
        const intervalSelect = document.getElementById('interval');
        if (intervalSelect) {
            intervalSelect.value = task.interval;
        }
    }
    
    const activeCheckbox = document.getElementById('active');
    if (activeCheckbox) {
        activeCheckbox.checked = task.active !== undefined ? task.active : true;
    }
    
    // Update field visibility based on populated data
    updateFieldVisibility();
    
    // Show/hide interval group based on recurring checkbox
    const intervalGroup = document.getElementById('interval-group');
    if (intervalGroup) {
        intervalGroup.style.display = task.recurring ? 'block' : 'none';
    }
    
    // Show/hide recent days based on historic checkbox
    const recentDaysField = document.querySelector('.recent-days-field');
    if (recentDaysField) {
        recentDaysField.style.display = (task.recent_days && task.recent_days > 0) ? 'block' : 'none';
    }
}

function updateFieldVisibility() {
    const selectedRadio = document.querySelector('input[name="task_type"]:checked');
    if (!selectedRadio) return;
    
    const selectedType = selectedRadio.value;
    const testFields = document.querySelectorAll('.test-field');
    const chartFields = document.querySelectorAll('.chart-field');
    const recentDaysField = document.querySelector('.recent-days-field');
    
    if (selectedType === 'test') {
        testFields.forEach(field => field.style.display = 'block');
        chartFields.forEach(field => field.style.display = 'none');
        if (recentDaysField) {
            recentDaysField.style.display = 'none';
        }
    } else {
        testFields.forEach(field => field.style.display = 'none');
        chartFields.forEach(field => field.style.display = 'block');
        // Show recent days only if historic is checked
        const historicCheckbox = document.getElementById('historic');
        const historicChecked = historicCheckbox ? historicCheckbox.checked : false;
        if (recentDaysField) {
            recentDaysField.style.display = historicChecked ? 'block' : 'none';
        }
    }
}

document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM Content Loaded');
    
    // Use event delegation for form submission since form might not exist yet
    document.body.addEventListener('htmx:beforeRequest', function(event) {
        const form = event.detail.elt;
        
        if (form && form.classList.contains('task-form')) {
            console.log('Intercepting task form submission');
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

            // Determine if this is create or edit based on form attributes
            const isEdit = form.hasAttribute('hx-put');
            const url = isEdit ? form.getAttribute('hx-put') : '/schedule/create';
            const method = isEdit ? 'PUT' : 'POST';

            console.log(`${method} request to ${url}:`, requestData);

            fetch(url, {
                method: method,
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestData)
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP ${response.status}: ${response.statusText}`);
                }
                // Both create and edit return JSON
                return response.json();
            })
            .then((data) => {
                console.log('Operation successful:', data);
                
                // Close modal and trigger refresh for both create and edit
                closeModal();
                htmx.trigger(document.body, 'taskChanged');
            })
            .catch(error => {
                console.error('Error:', error);
                alert('Error saving task: ' + error.message);
            });
        }
    });

    // Set up event listeners for form elements
    const recurringCheckbox = document.getElementById('recurring');
    if (recurringCheckbox) {
        recurringCheckbox.addEventListener('change', function() {
            const intervalGroup = document.getElementById('interval-group');
            if (intervalGroup) {
                intervalGroup.style.display = this.checked ? 'block' : 'none';
            }
        });
    }

    const taskTypeRadios = document.querySelectorAll('input[name="task_type"]');
    taskTypeRadios.forEach(radio => {
        radio.addEventListener('change', updateFieldVisibility);
    });

    const historicCheckbox = document.getElementById('historic');
    if (historicCheckbox) {
        historicCheckbox.addEventListener('change', function() {
            const recentDaysField = document.querySelector('.recent-days-field');
            if (recentDaysField) {
                recentDaysField.style.display = this.checked ? 'block' : 'none';
            }
        });
    }

    updateFieldVisibility();
});

// Handle successful form submission via HTMX (backup handler)
document.body.addEventListener('htmx:afterRequest', function(evt) {
    console.log('HTMX afterRequest:', evt.detail);
    if (evt.detail.successful && 
        (evt.detail.pathInfo.requestPath === '/schedule/create' || 
         evt.detail.pathInfo.requestPath.startsWith('/schedule/edit'))) {
        closeModal();
    }
});

// Close modal when clicking outside
window.onclick = function(event) {
    const modal = document.getElementById('task-modal');
    if (event.target === modal) {
        closeModal();
    }
};