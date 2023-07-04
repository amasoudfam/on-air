#!/bin/bash

# Build Go project
go build -o on-air .

# Check if the build was successful
if [ $? -ne 0 ]; then
  echo "Build failed. Exiting."
  exit 1
fi

# Run migrations
./on-air migrate --state=up

# Check if migration was successful
if [ $? -ne 0 ]; then
  echo "Migration failed. Exiting."
  exit 1
fi

# Run seed
./on-air seed

# Check if seed was successful
if [ $? -ne 0 ]; then
  echo "Seed failed. Exiting."
  exit 1
fi

# Run server
./on-air serve

# Check if server failed
if [ $? -ne 0 ]; then
  echo "Server failed. Exiting."
  exit 1
fi

# Cleanup - Uncomment the line below if you want to remove the built executable
# rm on-air
