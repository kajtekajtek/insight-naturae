/* internal/api/sensor_test.go - integration tests for the sensor 
	configuration endpoints */
package api

import (
	"testing"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"bytes"

	"github.com/stretchr/testify/assert"
	"github.com/kajtekajtek/insight-naturae/pkg/models"
)

// subscribe the sensor
func TestSubscribeSensor(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	token := SetupUser(t, db, r)

	// create the sensor payload
	requestBody := subscribeSensorRequestBody{
		SensorID: sensorID,
	}
	payload, _ := json.Marshal(requestBody)

	// add sensor to the user
	req, _ := http.NewRequest("POST", "/user/sensors", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Sensor added successfully")
}

// sensor subscription with an invalid payload
func TestSubscribeSensorInvalidPayload(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	token := SetupUser(t, db, r)

	// send an invalid payload as the sensor
	sensor := struct{
		password string
	}{
		password,
	}
	payload, _ := json.Marshal(sensor)

	// add sensor to the user
	req, _ := http.NewRequest("POST", "/user/sensors", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid input")
}

// sensor subscription with a missing field
func TestSubscribeSensorMissingField(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	token := SetupUser(t, db, r)

	// send a payload with a missing field
	sensor := struct{}{}
	payload, _ := json.Marshal(sensor)

	// add sensor to the user
	req, _ := http.NewRequest("POST", "/user/sensors", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid input")
}

// subscribe already subscribed sensor
func TestSubscribeSensorSubscribedSensor(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	token := SetupUser(t, db, r)

	// create the sensor payload
	sensor := subscribeSensorRequestBody{
		SensorID: sensorID,
	}
	payload, _ := json.Marshal(sensor)

	// subscribe the sensor
	req, _ := http.NewRequest("POST", "/user/sensors", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// subscribe the sensor again
	req, _ = http.NewRequest("POST", "/user/sensors", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer " + token)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "Sensor already exists")
}

// subscribe the sensor as an unauthorized user
func TestSubscribeSensorUnauthorized(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	// create the sensor payload
	sensor := subscribeSensorRequestBody{
		SensorID: sensorID,
	}
	payload, _ := json.Marshal(sensor)

	// subscribe the sensor
	req, _ := http.NewRequest("POST", "/user/sensors", bytes.NewBuffer(payload))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Missing Authorization Header")
}

// subscribe the sensor with a token passed in wrong format
func TestSubscribeSensorInvalidTokenHeaderFormat(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	token := SetupUser(t, db, r)
	
	// create the sensor payload
	sensor := subscribeSensorRequestBody{
		SensorID: sensorID,
	}
	payload, _ := json.Marshal(sensor)

	// subscribe the sensor
	req, _ := http.NewRequest("POST", "/user/sensors", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", token) // missing Bearer

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid Authorization Header")
}

// get sensors subscribed by the user
func TestGetSensors(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	token := SetupUser(t, db, r)

	// subscribe the sensor
	sensor := subscribeSensorRequestBody{
		SensorID: sensorID,
	}
	payload, _ := json.Marshal(sensor)

	req, _ := http.NewRequest("POST", "/user/sensors", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// get the subscribed sensors
	req, _ = http.NewRequest("GET", "/user/sensors", nil)
	req.Header.Set("Authorization", "Bearer " + token)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), sensorID)
}

// get subscribed sensors as an unauthorized user
func TestGetSensorsUnauthorized(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	// get the sensors
	req, _ := http.NewRequest("GET", "/user/sensors", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Missing Authorization Header")
}

// get subscribed sensors without any subscriptions
func TestGetSensorsNoSensors(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	token := SetupUser(t, db, r)

	// get the sensors
	req, _ := http.NewRequest("GET", "/user/sensors", nil)
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	t.Log(w.Body.String())	

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), "null")
}

// get subscribed sensors with an invalid token
func TestGetSensorsInvalidToken(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	// get the sensors
	req, _ := http.NewRequest("GET", "/user/sensors", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")
}

// get subscribed sensors with an invalid token header format
func TestGetSensorsInvalidTokenHeaderFormat(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	token := SetupUser(t, db, r)

	// get the sensors
	req, _ := http.NewRequest("GET", "/user/sensors", nil)
	req.Header.Set("Authorization", token) // missing Bearer

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid Authorization Header")
}

// unsubscribe the sensor from the user's list
func TestUnsubscribeSensor(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	token := SetupUser(t, db, r)

	// subscribe the sensor
	sensor := subscribeSensorRequestBody{
		SensorID: sensorID,
	}
	payload, _ := json.Marshal(sensor)

	req, _ := http.NewRequest("POST", "/user/sensors", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// unsubscribe the sensor
	req, _ = http.NewRequest("DELETE", "/user/sensors/" + sensorID, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer " + token)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Sensor removed successfully")

	// get the sensors
	req, _ = http.NewRequest("GET", "/user/sensors", nil)
	req.Header.Set("Authorization", "Bearer " + token)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), "null")
}

// unsubscribe the sensor with an invalid token
func TestUnsubscribeSensorInvalidToken(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	// unsubscribe the sensor
	req, _ := http.NewRequest("DELETE", "/user/sensors/" + sensorID, nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")
}

// unsubscribe the sensor as an unauthorized user
func TestUnsubscribeSensorUnauthorized(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	// unsubscribe the sensor
	req, _ := http.NewRequest("DELETE", "/user/sensors/" + sensorID, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Missing Authorization Header")
}

// unsubscribe a not subscribed sensor
func TestUnsubscribeSensorNotSubscribedSensor(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	token := SetupUser(t, db, r)

	// create the sensor payload
	sensor := subscribeSensorRequestBody{
		SensorID: sensorID,
	}
	payload, _ := json.Marshal(sensor)

	// unsubscribe the sensor
	req, _ := http.NewRequest("DELETE", "/user/sensors/" + sensorID, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Sensor not found")
}

// unsubscribe the sensor with an invalid payload
func TestUnsubscribeSensorInvalidPayload(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	token := SetupUser(t, db, r)

	// send an invalid payload as the sensor
	sensor := struct{
		password string
	}{
		password,
	}
	payload, _ := json.Marshal(sensor)

	// unsubscribe the sensor
	req, _ := http.NewRequest("DELETE", "/user/sensors/" + sensorID, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid input")
}

// unsubscribe the sensor with a missing field
func TestUnsubscribeSensorMissingField(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	token := SetupUser(t, db, r)

	// send a payload with a missing field
	sensor := models.SensorSubscription{
		Username: username,
	}
	payload, _ := json.Marshal(sensor)

	// unsubscribe the sensor
	req, _ := http.NewRequest("DELETE", "/user/sensors/" + sensorID, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid input")
}

// unsubscribe the sensor with an invalid token header format
func TestUnsubscribeSensorInvalidTokenHeaderFormat(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	token := SetupUser(t, db, r)

	// create the sensor payload
	sensor := subscribeSensorRequestBody{
		SensorID: sensorID,
	}
	payload, _ := json.Marshal(sensor)

	// unsubscribe the sensor
	req, _ := http.NewRequest("DELETE", "/user/sensors/" + sensorID, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", token) // missing Bearer

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid Authorization Header")
}