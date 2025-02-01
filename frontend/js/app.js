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

        document.getElementById("auth-container").style.display = "none";
        document.getElementById("dashboard-container").style.display = "block";

        connectWebSocket();

        alert("Successfully logged in.");
    } else {
        alert(data.error);
    }
}

// subscribe to a sensor
async function subscribeSensor() {
    const sensorID = document.getElementById('sensor-id').value;

    const response = await fetch(`${API_URL}/user/sensors`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ sensor_id: sensorID })
    });

    const data = await response.json();
    if (data.message) {
        alert(data.message);
        createChart(sensorID);
    } else {
        alert(data.error);
    }
}

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

        console.log('Received sensor data:', sensorData);

        updateChart(sensorData);
    };

    ws.onclose = () => {
        console.log('WebSocket closed. Reconnecting...');
        setTimeout(connectWebSocket, 5000);
    };
}

function createChart(sensorID) {
    if (charts[sensorID]) {
        return;
    }

    console.log('Creating chart for sensor:', sensorID)

    const chartContainer = document.getElementById('chart-container');
    const canvas = document.createElement('canvas');
    canvas.id = `char-${sensorID}`;
    chartContainer.appendChild(canvas);

    const ctx = canvas.getContext('2d');
    charts[sensorID] = new Chart(ctx, {
        type: 'line',
        data: {
            labels: [],
            datasets: [{
                label: `Sensor ${sensorID}`,
                data: [],
                borderColor: getRandomColor(),
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
            scales: {
                x: { type: 'category', title: { display: true, text: 'Time' } },
                y: { beginAtZero: false, title: { display: true, text: "Value" } }
            }
        }
    });
}

function updateChart(sensorData) {
    const sensorID = sensorData.sensor_id;
    const timestamp = new Date(sensorData.timestamp).toLocaleTimeString();

    console.log("Updating chart: ", sensorData)

    if (!charts[sensorID]) {
        console.log('Creating chart for sensor:', sensorID)
        createChart(sensorID);
    }

    const chart = charts[sensorID];
    chart.data.labels.push(timestamp);
    chart.data.datasets[0].data.push(sensorData.value);

    if (chart.data.labels.length > 20) {
        chart.data.labels.shift();
        chart.data.datasets[0].data.shift();
    }

    chart.update();
    console.log('Chart updated');
}

function getRandomColor() {
    return `hsl(${Math.random() * 360}, 100%, 50%)`;
}