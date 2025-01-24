// internal/dbutils/dbutils_test.go - dbutils.go unit tests
package dbutils

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/kajtekajtek/insight-naturae/pkg/models"
)

const testDBPath = "test.db"

func setupTestDB(t *testing.T) *sql.DB {
	// delete the test database if it exists
	if _, err := os.Stat(testDBPath); err == nil {
		err := os.Remove(testDBPath)
		assert.NoError(t, err)
	}

	// create the test database
	db, err := CreateDatabase(testDBPath)
	assert.NoError(t, err)
	assert.NotNil(t, db)
	return db
}

func tearDownTestDB(db *sql.DB) {
	db.Close()
	os.Remove(testDBPath)
}

// --- CREATE tables ---

func TestCreateSensorTable(t *testing.T) {
	db := setupTestDB(t)
	defer tearDownTestDB(db)

	err := CreateSensorTable(db)
	assert.NoError(t, err)
}

func TestCreateUserTable(t *testing.T) {
	db := setupTestDB(t)
	defer tearDownTestDB(db)

	err := CreateUserTable(db)
	assert.NoError(t, err)
}

func TestCreateUserSensorTable(t *testing.T) {
	db := setupTestDB(t)
	defer tearDownTestDB(db)

	err := CreateUserSensorTable(db)
	assert.NoError(t, err)
}

// --- INSERT operations ---

func TestInsertSensorData(t *testing.T) {
	db := setupTestDB(t)
	defer tearDownTestDB(db)

	data := models.SensorData{
		SensorID:  "test-sensor",
		Timestamp: "2021-01-01T00:00:00Z",
		Value:     42.0,
		Unit:      "test-unit",
	}

	err := InsertSensorData(db, data)
	assert.NoError(t, err)
}

func TestInsertUserData(t *testing.T) {
	db := setupTestDB(t)
	defer tearDownTestDB(db)

	user := models.User{
		Username: "testuser",
		Password: "securepassword",
	}

	err := InsertUserData(db, user)
	assert.NoError(t, err)
}

func TestInsertUserSensorData(t *testing.T) {
	db := setupTestDB(t)
	defer tearDownTestDB(db)

	userSensor := models.UserSensor{
		Username:   "testuser",
		SensorID: "testsensor",
	}

	err := InsertUserSensorData(db, userSensor)
	assert.NoError(t, err)
}

// --- SELECT operations ---

func TestDumpSensorData(t *testing.T) {
	db := setupTestDB(t)
	defer tearDownTestDB(db)

	data := models.SensorData{
		SensorID:  "test-sensor",
		Timestamp: "2021-01-01T00:00:00Z",
		Value:     42.0,
		Unit:      "test-unit",
	}

	err := InsertSensorData(db, data)
	assert.NoError(t, err)

	sensorData, err := DumpSensorData(db)
	assert.NoError(t, err)
	assert.Len(t, sensorData, 1)
	assert.Equal(t, data, sensorData[0])
}

func TestGetUsersByUsername(t *testing.T) {
	db := setupTestDB(t)
	defer tearDownTestDB(db)

	user := models.User{
		Username: "testuser",
		Password: "securepassword",
	}

	err := InsertUserData(db, user)
	assert.NoError(t, err)

	userData, err := GetUsersByUsername(db, "testuser")
	assert.NoError(t, err)
	assert.Len(t, userData, 1)
	assert.Equal(t, user.Username, userData[0].Username)
	assert.Equal(t, user.Password, userData[0].Password)
}

func TestGetUserSensors(t *testing.T) {
	db := setupTestDB(t)
	defer tearDownTestDB(db)

	userSensor := models.UserSensor{
		Username:   "testuser",
		SensorID: "testsensor",
	}

	err := InsertUserSensorData(db, userSensor)
	assert.NoError(t, err)

	userSensors, err := GetUserSensors(db, "testuser")
	assert.NoError(t, err)
	assert.Len(t, userSensors, 1)
	assert.Equal(t, userSensor.Username, userSensors[0].Username)
	assert.Equal(t, userSensor.SensorID, userSensors[0].SensorID)
}

// --- DELETE operations ---
func TestRemoveUserSensor(t *testing.T) {
	db := setupTestDB(t)
	defer tearDownTestDB(db)

	userSensor := models.UserSensor{
		Username:   "testuser",
		SensorID: "testsensor",
	}

	err := InsertUserSensorData(db, userSensor)
	assert.NoError(t, err)

	err = RemoveUserSensor(db, userSensor)
	assert.NoError(t, err)

	userSensors, err := GetUserSensors(db, "testuser")
	assert.NoError(t, err)
	assert.Len(t, userSensors, 0)
}