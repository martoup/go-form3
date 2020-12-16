#!/bin/sh
echo Waiting for accountapi
while ! nc -z accountapi 8080; do sleep 2; done
go run test/integration.go