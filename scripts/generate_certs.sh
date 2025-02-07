#!/bin/sh

if [ ! -d "certs" ]; then
  mkdir certs
fi

# generate CA certificate
openssl genrsa -out certs/ca.key 2048
openssl req -x509 -new -nodes -key certs/ca.key -sha256 -days 365 -out certs/ca.crt -subj "/CN=MyCA"

# generate server certificate
openssl genrsa -out certs/server.key 2048
openssl req -new -key certs/server.key -out certs/server.csr -subj "/CN=localhost"
openssl x509 -req -in certs/server.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/server.crt -days 365 -sha256

# generate mqtt client certificate
openssl genrsa -out certs/mqtt_client.key 2048
openssl req -new -key certs/mqtt_client.key -out certs/mqtt_client.csr -subj "/CN=mqtt-client"
openssl x509 -req -in certs/mqtt_client.csr -CA certs/ca.crt -CAkey certs/ca.key -CAcreateserial -out certs/mqtt_client.crt -days 365 -sha256