#!/bin/sh
set -e

echo "Running migrations..."
goose -dir migrations up

echo "Starting app..."
exec ./server
