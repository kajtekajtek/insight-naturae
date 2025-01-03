# scripts/simulate_sensors_mosquitto.sh - script for running the server 
# with the sensor simulator and mosquitto mqtt broker
#!/bin/sh

# start the mosquitto mqtt broker
mosquitto &
sleep 1
# start the sensor simulator
go run ./cmd/sensorsim/main.go &
# run the server
go run ./cmd/server/main.go

exit 0