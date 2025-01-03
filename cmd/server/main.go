package main

import (
	"log"
	"fmt"
	"github.com/kajtekajtek/insight-naturae/pkg/dbutils"
	"github.com/kajtekajtek/insight-naturae/internal/database"
	"github.com/kajtekajtek/insight-naturae/pkg/models"
)

func main() {
	dbPath := "./insight-naturae.db"
	db, err := dbutils.Init(dbPath)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	if err := database.CreateSensorTable(db); err != nil {
		log.Fatalf("Error creating sensor table: %v", err)
	}

	data := models.SensorData{
		SensorID: "sensor1",
		Timestamp: "2021-01-01T00:00:00Z",
		Value: 25.0,
		Unit: "C",
	}

	if err := database.InsertSensorData(db, data); err != nil {
		log.Fatalf("Error inserting sensor data: %v", err)
	}

	dataDump, err := database.DumpSensorData(db)
	if err != nil {
		log.Fatalf("Error dumping sensor data: %v", err)
	}

	fmt.Println("Sensor data:")
	for _, d := range dataDump {
		fmt.Printf("%+v\n", d)
	}
}