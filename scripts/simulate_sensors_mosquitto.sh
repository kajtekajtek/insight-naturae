# scripts/simulate_sensors_mosquitto.sh - script for running the sensor 
# simulation with mosquitto mqtt broker
#!/bin/sh

topic=$TOPIC
if [ -z "$topic" ]; then
    topic="insight-naturae/sensors"
fi

trap 'kill $pid1 $pid2; exit' INT TERM

# start the mosquitto mqtt broker
mosquitto &
pid1=$!
sleep 1
# start the sensor simulator
go run ./cmd/sensorsim/main.go &
pid2=$!
# subscribe to the topic
mosquitto_sub -t $topic -v &

while true; do
    sleep 1
done

exit 0