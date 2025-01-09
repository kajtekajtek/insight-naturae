// internal/database/dbutils.go - project specific database operations
package dbutils

import  (
	"database/sql"
	"fmt"
	"github.com/kajtekajtek/insight-naturae/pkg/models"
	"github.com/kajtekajtek/insight-naturae/pkg/utils"
	"github.com/kajtekajtek/insight-naturae/pkg/database"
	"github.com/joho/godotenv"
)

func CreateDatabase() error {
	// load the environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load .env file; continuing with the default values...")
	}

	// initialize the database
	dbPath := utils.Getenv("DB_FILE", "./insight-naturae.db")
	db, err := database.Init(dbPath)
	if err != nil {
		return err
	}

	// create the sensor table
	if err := CreateSensorTable(db); err != nil {
		return err
	}

	return nil
}

// create table for sensor data
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

// insert sensor data into the database
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

// dump sensor data from the database
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