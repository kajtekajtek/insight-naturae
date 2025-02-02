const API_URL = 'http://localhost:8080';
let token = null;
let ws = null;
let charts = {};

// user registration
async function register() {
    const username = document.getElementById('register-username').value;
    const password = document.getElementById('register-password').value;

    const response = await fetch(`${API_URL}/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
    });

    const data = await response.json();
    alert(data.message || data.error);
}

// user login
async function login() {
    const username = document.getElementById('login-username').value;
    const password = document.getElementById('login-password').value;

    const response = await fetch(`${API_URL}/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
    });

    const data = await response.json();
    if (data.token) {
        token = data.token;

        localStorage.setItem('token', token);

        initDashboard();

        await getUserSensors();

        alert("Successfully logged in.");
    } else {
        alert(data.error);
    }
}

// user logout
function logout() {
    localStorage.removeItem('token');

    token = null;

    ws?.close();

    document.getElementById("auth-container").style.display = "block";
    document.getElementById("dashboard-container").style.display = "none";

    alert("Successfully logged out.");
}

// check if user is logged in
function checkLoginStatus() {
    token = localStorage.getItem('token');
    if (token) {
        initDashboard();
    }
}

// subscribe to a sensor
async function subscribeSensor() {
    const sensorID = document.getElementById('sensor-id').value;

    const response = await fetch(`${API_URL}/user/sensors/${sensorID}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
    });

    const data = await response.json();
    if (data.message) {
        alert(data.message);
        createChart(sensorID);
    } else {
        alert(data.error);
    }
}

// unsubscribe from a sensor
async function unsubscribeSensor(sensorID) {
    const response = await fetch(`${API_URL}/user/sensors/${sensorID}`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
    });

    const data = await response.json();
    if (data.message) {
        alert(data.message);

        removeChart(sensorID);

        document.getElementById(`sensor-item-${sensorID}`).remove();
    } else {
        alert(data.error);
    }
}

// create a chart for a sensor
async function createChart(sensorID) {
    if (charts[sensorID]) {
        return;
    }

    console.log('Creating chart for sensor:', sensorID)

    // create a canvas element for the chart
    const chartContainer = document.getElementById('chart-container');
    const canvas = document.createElement('canvas');
    canvas.id = `chart-${sensorID}`;
    chartContainer.appendChild(canvas);

    // get historical sensor data
    const historicalData = await getSensorData(sensorID);
    const labels = historicalData.map(r => new Date(r.timestamp)
        .toLocaleTimeString());
    const values = historicalData.map(r => r.value);

    // create the chart
    const ctx = canvas.getContext('2d');
    charts[sensorID] = new Chart(ctx, {
        type: 'line',
        data: {
            labels: labels,
            datasets: [{
                label: `Sensor ${sensorID}`,
                data: values,
                borderColor: getRandomColor(),
                borderWidth: 2,
                fill: false
            }]
        },
        options: {
            responsive: true,
            // maintain the aspect ratio
            scales: {
                x: { type: 'category', title: { display: true, text: 'Time' } },
                y: { beginAtZero: false, title: { display: true, text: "Value" } }
            },
            // remove points on the line
            elements: {
                point: {
                    radius: 0
                }
            }
        }
    });

    initDashboard();
}

// remove chart from the dashboard
function removeChart(sensorID) {
    if (charts[sensorID]) {
        delete charts[sensorID];
        document.getElementById(`chart-${sensorID}`).remove();
    }
}

// update the chart with new sensor data
function updateChart(sensorData) {
    const sensorID = sensorData.sensor_id;
    const timestamp = new Date(sensorData.timestamp).toLocaleTimeString();

    if (!charts[sensorID]) {
        createChart(sensorID);
    }

    const chart = charts[sensorID];
    chart.data.labels.push(timestamp);
    chart.data.datasets[0].data.push(sensorData.value);

    // shift the labels and data if more than 20
    if (chart.data.labels.length > 20) {
        chart.data.labels.shift();
        chart.data.datasets[0].data.shift();
    }

    chart.update();
}

// connect to the WebSocket
function connectWebSocket() {
    if (!token) {
        console.error('Missing token');
        return;
    }

    const wsURL = `${API_URL}/ws?token=${token}`;
    ws = new WebSocket(wsURL);

    // on connection open
    ws.onopen = () => {
        console.log('Connected to WebSocket');
    };

    // on message received
    ws.onmessage = (event) => {
        if (event.data === 'ping') {
            console.log('Received ping')
            return;
        }
        // parse the sensor data
        const sensorData = JSON.parse(event.data);

        updateChart(sensorData);
    };

    ws.onclose = () => {
        console.log('WebSocket closed. Reconnecting...');
        setTimeout(connectWebSocket, 5000);
    };
}

// initialize the dashboard
async function initDashboard() {
    document.getElementById("auth-container").style.display = "none";
    document.getElementById("dashboard-container").style.display = "block";

    // sensor list element to display subscribed sensors
    const sensorList = document.getElementById('sensor-list');
    sensorList.innerHTML = "";    

    // get subscribed sensors and create charts
    await getUserSensors().then(data => {
        if (Array.isArray(data)) {
            data.forEach(sensor => {
                // create a chart for each sensor
                createChart(sensor.sensor_id);

                // check if the unsubscribe button is already created
                if (document.getElementById(`sensor-item-${sensor.sensor_id}`)) {
                    return;
                }                

                /* add sensor to the subscribed sensors list with 
                    an unsubscribe button */
                const sensorItem = document.createElement('div');
                sensorItem.id = `sensor-item-${sensor.sensor_id}`;
                sensorItem.innerHTML = `
                    <span>Sensor ${sensor.sensor_id}</span>
                    <button class="unsubscribe-btn" data-sensor-id="${sensor.sensor_id}">unsubscribe</button>`;

                sensorList.appendChild(sensorItem);
            });

            // add event listener to the unsubscribe button
            document.querySelectorAll('.unsubscribe-btn').forEach(btn => {
                btn.addEventListener('click', () => {
                    const sensorID = btn.getAttribute('data-sensor-id');
                    unsubscribeSensor(sensorID);
                });
            });
        }
    })

    connectWebSocket();
}

// get random color in HSL format
function getRandomColor() {
    return `hsl(${Math.random() * 360}, 100%, 50%)`;
}

// get user subscribed sensors
async function getUserSensors() {
    const response = await fetch(`${API_URL}/user/sensors`, {
        method: "GET",
        headers: { 'Authorization': `Bearer ${token}` }
    });

    const data = await response.json();
    return data;
}

// get historical sensor data
async function getSensorData(sensorID) {
    const response = await fetch(`${API_URL}/user/sensors/${sensorID}/data`, {
        method: "GET",
        headers: { 'Authorization': `Bearer ${token}` }
    });

    const data = await response.json();
    return data;
}

// check if the user is logged in
window.onload = checkLoginStatus;