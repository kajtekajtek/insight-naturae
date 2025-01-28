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

// add a sensor to the user
func TestAddSensor(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupRouter(db)

	token := SetupUser(t, db, r)

	// create the sensor payload
	sensor := models.UserSensor{
		Username: username,
		SensorID: sensorID,
	}
	payload, _ := json.Marshal(sensor)

	// add sensor to the user
	req, _ := http.NewRequest("POST", "/user/sensors", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Sensor added successfully")
}

// add a sensor with an invalid payload
func TestAddSensorInvalidPayload(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupRouter(db)

	token := SetupUser(t, db, r)

	// send an invalid payload as a sensor
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

// add a sensor with a missing field
func TestAddSensorMissingField(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupRouter(db)

	token := SetupUser(t, db, r)

	// send a payload with a missing field
	sensor := models.UserSensor{
		Username: username,
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

// add a sensor with an invalid token
func TestAddSensorExistingSensor(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupRouter(db)

	token := SetupUser(t, db, r)

	// create the sensor payload
	sensor := models.UserSensor{
		Username: username,
		SensorID: sensorID,
	}
	payload, _ := json.Marshal(sensor)

	// add sensor to the user
	req, _ := http.NewRequest("POST", "/user/sensors", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// add the same sensor again
	req, _ = http.NewRequest("POST", "/user/sensors", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer " + token)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "Sensor already exists")
}

// add a sensor as an unauthorized user
func TestAddSensorUnauthorized(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupRouter(db)

	// create the sensor payload
	sensor := models.UserSensor{
		Username: username,
		SensorID: sensorID,
	}
	payload, _ := json.Marshal(sensor)

	// add sensor to the user
	req, _ := http.NewRequest("POST", "/user/sensors", bytes.NewBuffer(payload))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Missing Authorization Header")
}

// add a sensor with a token passed in wrong format
func TestAddSensorInvalidTokenHeaderFormat(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupRouter(db)

	token := SetupUser(t, db, r)
	
	// create the sensor payload
	sensor := models.UserSensor{
		Username: username,
		SensorID: sensorID,
	}
	payload, _ := json.Marshal(sensor)

	// add sensor to the user
	req, _ := http.NewRequest("POST", "/user/sensors", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", token) // missing Bearer

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid Authorization Header")
}

// get sensors followed by the user
func TestGetSensors(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupRouter(db)

	token := SetupUser(t, db, r)

	// add a sensor to the user
	sensor := models.UserSensor{
		Username: username,
		SensorID: sensorID,
	}
	payload, _ := json.Marshal(sensor)

	req, _ := http.NewRequest("POST", "/user/sensors", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// get the sensors
	req, _ = http.NewRequest("GET", "/user/sensors", nil)
	req.Header.Set("Authorization", "Bearer " + token)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), sensorID)
}

// get sensors as an unauthorized user
func TestGetSensorsUnauthorized(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupRouter(db)

	// get the sensors
	req, _ := http.NewRequest("GET", "/user/sensors", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Missing Authorization Header")
}

// get sensors without any sensors
func TestGetSensorsNoSensors(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupRouter(db)

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

// get sensors with an invalid token
func TestGetSensorsInvalidToken(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupRouter(db)

	// get the sensors
	req, _ := http.NewRequest("GET", "/user/sensors", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")
}

// get sensors with an invalid token header format
func TestGetSensorsInvalidTokenHeaderFormat(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupRouter(db)

	token := SetupUser(t, db, r)

	// get the sensors
	req, _ := http.NewRequest("GET", "/user/sensors", nil)
	req.Header.Set("Authorization", token) // missing Bearer

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid Authorization Header")
}

// remove a sensor from the user's list
func TestRemoveSensor(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupRouter(db)

	token := SetupUser(t, db, r)

	// add a sensor to the user
	sensor := models.UserSensor{
		Username: username,
		SensorID: sensorID,
	}
	payload, _ := json.Marshal(sensor)

	req, _ := http.NewRequest("POST", "/user/sensors", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// remove the sensor
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

// remove a sensor with an invalid token
func TestRemoveSensorInvalidToken(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupRouter(db)

	// remove the sensor
	req, _ := http.NewRequest("DELETE", "/user/sensors/" + sensorID, nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid token")
}

// remove a sensor as an unauthorized user
func TestRemoveSensorUnauthorized(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupRouter(db)

	// remove the sensor
	req, _ := http.NewRequest("DELETE", "/user/sensors/" + sensorID, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Missing Authorization Header")
}

// remove a non-existing sensor
func TestRemoveSensorNonExistingSensor(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupRouter(db)

	token := SetupUser(t, db, r)

	// create the sensor payload
	sensor := models.UserSensor{
		Username: username,
		SensorID: sensorID,
	}
	payload, _ := json.Marshal(sensor)

	// remove the sensor
	req, _ := http.NewRequest("DELETE", "/user/sensors/" + sensorID, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Sensor not found")
}

// remove a sensor with an invalid payload
func TestRemoveSensorInvalidPayload(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupRouter(db)

	token := SetupUser(t, db, r)

	// send an invalid payload as a sensor
	sensor := struct{
		password string
	}{
		password,
	}
	payload, _ := json.Marshal(sensor)

	// remove the sensor
	req, _ := http.NewRequest("DELETE", "/user/sensors/" + sensorID, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid input")
}

// remove a sensor with a missing field
func TestRemoveSensorMissingField(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupRouter(db)

	token := SetupUser(t, db, r)

	// send a payload with a missing field
	sensor := models.UserSensor{
		Username: username,
	}
	payload, _ := json.Marshal(sensor)

	// remove the sensor
	req, _ := http.NewRequest("DELETE", "/user/sensors/" + sensorID, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid input")
}

// remove a sensor with an invalid token header format
func TestRemoveSensorInvalidTokenHeaderFormat(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupRouter(db)

	token := SetupUser(t, db, r)

	// create the sensor payload
	sensor := models.UserSensor{
		Username: username,
		SensorID: sensorID,
	}
	payload, _ := json.Marshal(sensor)

	// remove the sensor
	req, _ := http.NewRequest("DELETE", "/user/sensors/" + sensorID, bytes.NewBuffer(payload))
	req.Header.Set("Authorization", token) // missing Bearer

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid Authorization Header")
}