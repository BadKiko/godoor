#!/bin/bash

# Simple script to build and run GoDoor with Docker

echo "Building Docker image..."
docker build -t godoor .

echo "Starting container..."
docker run -d \
  --name godoor \
  -p 8080:8080 \
  -v $(pwd)/godoor.db:/app/godoor.db \
  godoor

echo "Container started. Check logs with: docker logs godoor"
echo "Stop with: docker stop godoor && docker rm godoor"