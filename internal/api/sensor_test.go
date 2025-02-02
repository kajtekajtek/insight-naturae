/* internal/api/sensor_test.go - integration tests for the sensor 
	configuration endpoints */
package api

import (
	"testing"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"github.com/kajtekajtek/insight-naturae/pkg/models"
	"github.com/kajtekajtek/insight-naturae/internal/dbutils"
)

// subscribe the sensor
func TestSubscribeSensor(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	token := SetupUser(t, db, r)

	// subscribe the sensor
	req, _ := http.NewRequest("POST", "/user/sensors/" + sensorID, nil)
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Sensor subscribed successfully")
}

// subscribe the same sensor twice
func TestSubscribeSensorTwice(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	token := SetupUser(t, db, r)

	// subscribe the sensor
	req, _ := http.NewRequest("POST", "/user/sensors/" + sensorID, nil)
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// subscribe the sensor again
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "Sensor already subscribed")
}

// try to subscribe the sensor unauthorized
func TestSubscribeSensorUnauthorized(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	// subscribe the sensor
	req, _ := http.NewRequest("POST", "/user/sensors/" + sensorID, nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Missing Authorization Header")
}

// get sensors subscribed by the user
func TestGetSensors(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	token := SetupUser(t, db, r)

	req, _ := http.NewRequest("POST", "/user/sensors/" + sensorID, nil)
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

// 

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

// get sensor data
func TestGetSensorData(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)

	// insert sensor data to the database
	dbutils.InsertSensorData(db, models.SensorData{
		SensorID: sensorID,
		Timestamp: "2021-01-01T00:00:00Z",
		Value: 42.0,
		Unit: "test-unit",
	})
	
	r := SetupTestRouter(db)

	token := SetupUser(t, db, r)

	// subscribe the sensor
	req, _ := http.NewRequest("POST", "/user/sensors/" + sensorID, nil)
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)	

	// get the sensor data
	req, _ = http.NewRequest("GET", "/user/sensors/" + sensorID + "/data", nil)
	req.Header.Set("Authorization", "Bearer " + token)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), sensorID)
}

// get sensor data as an unauthorized user
func TestGetSensorDataUnauthorized(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	// get the sensor data
	req, _ := http.NewRequest("GET", "/user/sensors/" + sensorID + "/data", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "Missing Authorization Header")
}

// get sensor data from a sensor that does not exist
func TestGetSensorDataNonExistentSensor(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	token := SetupUser(t, db, r)

	// get the sensor data
	req, _ := http.NewRequest("GET", "/user/sensors/" + sensorID + "/data", nil)
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "No data found")
}

// unsubscribe the sensor from the user's list
func TestUnsubscribeSensor(t *testing.T) {
	db := SetupTestDB(t)
	defer TearDownTestDB(db)
	
	r := SetupTestRouter(db)

	token := SetupUser(t, db, r)

	// subscribe the sensor
	req, _ := http.NewRequest("POST", "/user/sensors/" + sensorID, nil)
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// unsubscribe the sensor
	req, _ = http.NewRequest("DELETE", "/user/sensors/" + sensorID, nil)
	req.Header.Set("Authorization", "Bearer " + token)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Sensor unsubscribed successfully")

	// get the sensors
	req, _ = http.NewRequest("GET", "/user/sensors", nil)
	req.Header.Set("Authorization", "Bearer " + token)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, w.Body.String(), "null")
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

	// unsubscribe the sensor
	req, _ := http.NewRequest("DELETE", "/user/sensors/" + sensorID, nil)
	req.Header.Set("Authorization", "Bearer " + token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Sensor not found")
}