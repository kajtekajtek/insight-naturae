// internal/database/dbutils.go - project specific database operations
package dbutils

import  (
	"database/sql"
	"fmt"

	"github.com/kajtekajtek/insight-naturae/pkg/models"
	"github.com/kajtekajtek/insight-naturae/pkg/database"
)

func CreateDatabase(DBPath string) (*sql.DB, error) {
	db, err := database.Init(DBPath)
	if err != nil {
		return nil, err
	}

	// create the sensor table
	if err := CreateSensorTable(db); err != nil {
		return nil, err
	}

	// create the user table
	if err := CreateUserTable(db); err != nil {
		return nil, err
	}

	// create the user sensor table
	if err := CreateSensorSubscriptionTable(db); err != nil {
		return nil, err
	}

	return db, nil
}

// --- CREATE tables --- 

func CreateSensorTable(db *sql.DB) error {
	sensorTableSQL := `
		CREATE TABLE IF NOT EXISTS SensorData (
			data_id INTEGER PRIMARY KEY AUTOINCREMENT,
			sensor_id TEXT NOT NULL,
			timestamp TEXT NOT NULL,
			value REAL NOT NULL,
			unit TEXT NOT NULL
		);`

	_, err := db.Exec(sensorTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create sensor table: %v", err)
	}

	return nil
}

func CreateUserTable(db *sql.DB) error {
	userTableSQL := `
		CREATE TABLE IF NOT EXISTS Users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);`

	_, err := db.Exec(userTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create user table: %v", err)
	}

	return nil;
}

func CreateSensorSubscriptionTable(db *sql.DB) error {
	sensorSubscriptionTableSQL := `
		CREATE TABLE IF NOT EXISTS SensorSubscription (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL,
			sensor_id TEXT NOT NULL
		);`

	_, err := db.Exec(sensorSubscriptionTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create user sensor table: %v", err)
	}

	return nil
}

// --- INSERT data ---

func InsertSensorData(db *sql.DB, data models.SensorData) error {
	insertSQL := `
		INSERT INTO SensorData (sensor_id, timestamp, value, unit)
		VALUES (?, ?, ?, ?);`
	
	_, err := db.Exec(insertSQL, data.SensorID, data.Timestamp, data.Value, data.Unit)
	if err != nil {
		return fmt.Errorf("failed to insert sensor data: %v", err)
	}

	return nil
}

func InsertUserData(db *sql.DB, user models.User) error {
	insertSQL := `
		INSERT INTO Users (username, password)
		VALUES (?, ?);`

	_, err := db.Exec(insertSQL, user.Username, user.Password)
	if err != nil {
		return fmt.Errorf("failed to insert user data: %v", err)
	}

	return nil
}

func InsertSensorSubscriptionData(db *sql.DB, sensorSubscription models.SensorSubscription) error {
	insertSQL := `
		INSERT INTO SensorSubscription (username, sensor_id)
		VALUES (?, ?);`

	_, err := db.Exec(insertSQL, sensorSubscription.Username, sensorSubscription.SensorID)
	if err != nil {
		return fmt.Errorf("failed to insert user sensor data: %v", err)
	}

	return nil
}

// --- SELECT data ---

func DumpSensorData(db *sql.DB) ([]models.SensorData, error) {
	querySQL := `
		SELECT sensor_id, timestamp, value, unit FROM SensorData;`
	
	rows, err := db.Query(querySQL)
	if err != nil {
		return nil, fmt.Errorf("failed to query sensor data: %v", err)
	}
	defer rows.Close()

	// iterate over the query results
	var data []models.SensorData
	for rows.Next() {
		var sd models.SensorData
		// copy each column into a field in the struct
		err := rows.Scan(&sd.SensorID, &sd.Timestamp, &sd.Value, &sd.Unit)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sensor data: %v", err) 
		}
		data = append(data, sd)
	}

	return data, nil
}

func GetUsersByUsername(db *sql.DB, username string) ([]models.User, error) {
	querySQL := `
		SELECT username, password FROM Users WHERE username = ?;`
	
	rows, err := db.Query(querySQL, username)
	if err != nil {
		return nil, fmt.Errorf("failed to query user data: %v", err)
	}
	defer rows.Close()

	// iterate over the query results
	var user []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.Username, &u.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user data: %v", err)
		}
		user = append(user, u)
	}

	return user, nil
}

func GetUserSubscriptions(db *sql.DB, username string) ([]models.SensorSubscription, error) {
	querySQL := `
		SELECT username, sensor_id FROM SensorSubscription WHERE username = ?;`
	
	rows, err := db.Query(querySQL, username)
	if err != nil {
		return nil, fmt.Errorf("failed to query user sensor data: %v", err)
	}
	defer rows.Close()

	// iterate over the query results
	var sensorSubscription []models.SensorSubscription
	for rows.Next() {
		var us models.SensorSubscription
		err := rows.Scan(&us.Username, &us.SensorID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user sensor data: %v", err)
		}
		sensorSubscription = append(sensorSubscription, us)
	}

	return sensorSubscription, nil
}

// get sensor data by sensor_id
func GetSensorData(db *sql.DB, sensorID string) ([]models.SensorData, error) {
	querySQL := `
		SELECT sensor_id, timestamp, value, unit FROM SensorData WHERE sensor_id = ?;`
	
	rows, err := db.Query(querySQL, sensorID)
	if err != nil {
		return nil, fmt.Errorf("failed to query sensor data: %v", err)
	}
	defer rows.Close()

	// iterate over the query results
	var data []models.SensorData
	for rows.Next() {
		var sd models.SensorData
		err := rows.Scan(&sd.SensorID, &sd.Timestamp, &sd.Value, &sd.Unit)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sensor data: %v", err)
		}
		data = append(data, sd)
	}

	return data, nil
}


// --- DELETE data ---
func RemoveSensorSubscription(db *sql.DB, sensorSubscription models.SensorSubscription) error {
	deleteSQL := `
		DELETE FROM SensorSubscription WHERE username = ? AND sensor_id = ?;`

	_, err := db.Exec(deleteSQL, sensorSubscription.Username, sensorSubscription.SensorID)
	if err != nil {
		return fmt.Errorf("failed to delete user sensor data: %v", err)
	}

	return nil
}