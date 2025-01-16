// internal/database/dbutils.go - project specific database operations
package dbutils

import  (
	"database/sql"
	"fmt"

	"github.com/kajtekajtek/insight-naturae/pkg/models"
	"github.com/kajtekajtek/insight-naturae/pkg/utils"
	"github.com/kajtekajtek/insight-naturae/pkg/database"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func CreateDatabase() (*sql.DB, error) {
	// load the environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load .env file; continuing with the default values...")
	}

	// initialize the database
	dbPath := utils.Getenv("DB_FILE", "./insight-naturae.db")
	db, err := database.Init(dbPath)
	if err != nil {
		return nil, err
	}

	// create the sensor table
	if err := CreateSensorTable(db); err != nil {
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
	hashedPsswrd, err := bcrypt.GenerateFromPassword([]byte(user.Password),
		bcrypt.DefaultCost)	
	if err != nil {
		return fmt.Errorf("failed to insert user data: %v", err)
	}

	insertSQL := `
		INSERT INTO Users (username, password)
		VALUES (?, ?);`

	_, err := db.Exec(insertSQL, user.Username, string(hashedPsswrd))
	if err != nil {
		return fmt.Errorf("failed to insert user data: %v", err)
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

func GetUserByUsername(db *sql.DB, username string) (models.User, error) {
	querySQL = `
		SELECT username, password FROM Users WHERE username = ?;`
	
	row, err := db.QueryRow(querySQL, username)
	if err != nil {
		return nil, fmt.Errorf("failed to query user data: %v", err)
	}
	defer row.Close()

	var user models.User
	// copy each column into a field in the struct
	err := row.Scan(&user.Username, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to scan user data: %v", err)
	}

	return user, nil
}