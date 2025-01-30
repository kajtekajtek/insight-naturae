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

func TestCreateSensorSubscriptionTable(t *testing.T) {
	db := setupTestDB(t)
	defer tearDownTestDB(db)

	err := CreateSensorSubscriptionTable(db)
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

func TestInsertSensorSubscriptionData(t *testing.T) {
	db := setupTestDB(t)
	defer tearDownTestDB(db)

	sensorSubscription := models.SensorSubscription{
		Username:   "testuser",
		SensorID: "testsensor",
	}

	err := InsertSensorSubscriptionData(db, sensorSubscription)
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

func TestGetSensorSubscriptions(t *testing.T) {
	db := setupTestDB(t)
	defer tearDownTestDB(db)

	sensorSubscription := models.SensorSubscription{
		Username:   "testuser",
		SensorID: "testsensor",
	}

	err := InsertSensorSubscriptionData(db, sensorSubscription)
	assert.NoError(t, err)

	sensorSubscriptions, err := GetSensorSubscriptions(db, "testuser")
	assert.NoError(t, err)
	assert.Len(t, sensorSubscriptions, 1)
	assert.Equal(t, sensorSubscription.Username, sensorSubscriptions[0].Username)
	assert.Equal(t, sensorSubscription.SensorID, sensorSubscriptions[0].SensorID)
}

// --- DELETE operations ---
func TestRemoveSensorSubscription(t *testing.T) {
	db := setupTestDB(t)
	defer tearDownTestDB(db)

	sensorSubscription := models.SensorSubscription{
		Username:   "testuser",
		SensorID: "testsensor",
	}

	err := InsertSensorSubscriptionData(db, sensorSubscription)
	assert.NoError(t, err)

	err = RemoveSensorSubscription(db, sensorSubscription)
	assert.NoError(t, err)

	sensorSubscriptions, err := GetSensorSubscriptions(db, "testuser")
	assert.NoError(t, err)
	assert.Len(t, sensorSubscriptions, 0)
}