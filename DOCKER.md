# Docker Setup for Social Scraper

This document explains how to run the Social Scraper application using Docker and Docker Compose.

## Prerequisites

- Docker
- Docker Compose
- Git

## Quick Start

### Option 1: Automated Setup (Recommended)
```bash
# Clone the repository
git clone <your-repo-url>
cd social-scraper

# Use the quick start script
./docker-start.sh
```

### Option 2: Manual Setup
1. **Clone the repository**:
   ```bash
   git clone <your-repo-url>
   cd social-scraper
   ```

2. **Set up environment variables**:
   ```bash
   cp .env.docker.example .env.docker
   # Edit .env.docker with your actual credentials
   ```

3. **Build and start the application**:
   ```bash
   # Build the Docker image
   ./docker-build.sh
   
   # Start services
   docker-compose up -d
   ```

4. **Access the application**:
   - Web Interface: http://localhost:8080
   - MongoDB: localhost:27017

## Environment Configuration

### Required Environment Variables

Create a `.env` file with the following variables:

```bash
# Instagram Credentials (required for Instagram scraping)
INSTAGRAM_USERNAME=your_instagram_username
INSTAGRAM_PASSWORD=your_instagram_password

# JWT Secret for authentication
JWT_SECRET=your_super_secure_jwt_secret_key

# MongoDB Passwords
MONGO_ROOT_PASSWORD=secure_root_password
MONGO_USER_PASSWORD=secure_user_password

# Optional
ECHO_VAR=success
```

### Security Notes

- **Change default passwords**: Always use secure, unique passwords for production
- **JWT Secret**: Use a long, random string for the JWT secret
- **Instagram Credentials**: Use a dedicated account for scraping

## Services

### Social Scraper Application
- **Port**: 8080
- **Health Check**: http://localhost:8080/health
- **Features**: 
  - Multi-platform social media scraping
  - File upload and processing
  - User authentication
  - Results export

### MongoDB Database
- **Port**: 27017
- **Database**: `influencer_db`
- **Collections**: `influencers`, `users`
- **Authentication**: Enabled with custom user

## Docker Commands

### Start services
```bash
docker-compose up -d
```

### View logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f social-scraper
docker-compose logs -f mongo
```

### Stop services
```bash
docker-compose down
```

### Rebuild application
```bash
docker-compose build social-scraper
docker-compose up -d social-scraper
```

### Clean up (removes volumes)
```bash
docker-compose down -v
```

## Development

### Local Development with Docker
```bash
# Start only MongoDB
docker-compose up -d mongo

# Run application locally
make run-http
```

### Debugging
```bash
# Access application container
docker exec -it social-scraper-app /bin/bash

# Access MongoDB container
docker exec -it social-scraper-mongo mongosh
```

## Data Persistence

- **MongoDB Data**: Stored in Docker volume `mongo_data`
- **Uploads**: Mounted to `./uploads` directory
- **Results**: Mounted to `./results` directory

## Architecture

```
┌─────────────────┐    ┌──────────────────┐
│   Load Balancer │    │   Social Scraper │
│   (Future)      │────│   Application    │
└─────────────────┘    │   (Port 8080)    │
                       └──────────────────┘
                                │
                                │
                       ┌──────────────────┐
                       │   MongoDB        │
                       │   (Port 27017)   │
                       └──────────────────┘
```

## Troubleshooting

### Common Issues

1. **Port conflicts**:
   ```bash
   # Check if ports are in use
   netstat -tlnp | grep :8080
   netstat -tlnp | grep :27017
   ```

2. **Permission errors**:
   ```bash
   # Fix ownership of mounted directories
   sudo chown -R $USER:$USER uploads/ results/
   ```

3. **MongoDB connection issues**:
   ```bash
   # Check MongoDB logs
   docker-compose logs mongo
   
   # Test connection
   docker exec -it social-scraper-mongo mongosh -u admin -p
   ```

4. **Puppeteer issues**:
   ```bash
   # Check if Chrome is properly installed
   docker exec -it social-scraper-app google-chrome-stable --version
   ```

### Health Checks

```bash
# Application health
curl http://localhost:8080/health

# MongoDB health
docker exec social-scraper-mongo mongosh --eval "db.adminCommand('ping')"
```

## Production Deployment

For production deployment, consider:

1. **Use secrets management** for sensitive environment variables
2. **Enable SSL/TLS** with reverse proxy (nginx, traefik)
3. **Set up monitoring** with health checks
4. **Configure backup strategy** for MongoDB
5. **Use multi-stage builds** for smaller images
6. **Implement log aggregation**

## Performance Tuning

- **MongoDB**: Adjust memory settings based on data size
- **Puppeteer**: Configure concurrency limits for scraping
- **Go Application**: Tune worker pools and timeouts

## Support

For issues and questions:
1. Check application logs: `docker-compose logs social-scraper`
2. Check MongoDB logs: `docker-compose logs mongo`
3. Verify environment configuration
4. Test individual components
