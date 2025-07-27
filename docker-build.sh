#!/bin/bash

# Docker build script for social-scraper
echo "Building social-scraper Docker image..."

# Build the image
docker build -t social-scraper:latest .

if [ $? -eq 0 ]; then
    echo "✅ Docker image built successfully!"
    echo "🚀 You can now run: docker-compose up"
else
    echo "❌ Docker build failed!"
    exit 1
fi
