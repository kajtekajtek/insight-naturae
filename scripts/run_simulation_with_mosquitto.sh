# scripts/simulate_sensors_mosquitto.sh - script for running the server 
# with the sensor simulator and mosquitto mqtt broker
#!/bin/sh

cleanup() {
    pkill mosquitto

    pkill -f "go run ./cmd/sensorsim/main.go"

    pkill -f "go run ./cmd/server/main.go"
}

# Przypisanie funkcji cleanup do sygnałów zakończenia skryptu
trap cleanup EXIT

# start the mosquitto mqtt broker
mosquitto -d
sleep 1
# start the sensor simulator
go run ./cmd/sensorsim/main.go &
# run the server
go run ./cmd/server/main.go

exit 0