#!/bin/bash

# Quick start script for social-scraper Docker setup
echo "🐋 Starting Social Scraper Docker Setup..."

# Check if .env.docker exists
if [ ! -f ".env.docker" ]; then
    echo "⚠️  .env.docker not found! Creating from template..."
    cp .env.docker.example .env.docker
    echo "📝 Please edit .env.docker with your actual values before proceeding"
    echo "   Required: MONGO_USERNAME, MONGO_PASSWORD, JWT_SECRET"
    exit 1
fi

# Build the Docker image
echo "🔨 Building Docker image..."
docker build -t social-scraper:latest .

if [ $? -ne 0 ]; then
    echo "❌ Docker build failed!"
    exit 1
fi

# Start the services
echo "🚀 Starting services with docker-compose..."
docker-compose up -d

if [ $? -eq 0 ]; then
    echo "✅ Services started successfully!"
    echo ""
    echo "📋 Service URLs:"
    echo "   - Application: http://localhost:8080"
    echo "   - MongoDB: localhost:27017"
    echo ""
    echo "📊 Check status: docker-compose ps"
    echo "📝 View logs: docker-compose logs -f"
    echo "🛑 Stop services: docker-compose down"
else
    echo "❌ Failed to start services!"
    exit 1
fi
