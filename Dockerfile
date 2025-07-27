# Multi-stage build for social-scraper with Puppeteer support

# Stage 1: Build the Go application
FROM golang:1.23-alpine AS go-builder

WORKDIR /app

# Install git (required for go mod download)
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the HTTP server binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o http-server ./cmd/http

# Stage 2: Node.js image with Puppeteer and Chrome
FROM node:20-bullseye-slim

# Install necessary system packages and Chromium (works on both amd64 and arm64)
RUN apt-get update && apt-get install -y \
    wget \
    gnupg \
    ca-certificates \
    chromium \
    fonts-liberation \
    libappindicator3-1 \
    libasound2 \
    libatk-bridge2.0-0 \
    libdrm2 \
    libgtk-3-0 \
    libnspr4 \
    libnss3 \
    libxcomposite1 \
    libxrandr2 \
    libxss1 \
    libxtst6 \
    xdg-utils \
    && rm -rf /var/lib/apt/lists/*

# Create app directory
WORKDIR /app

# Set Puppeteer environment variables
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true \
    PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium

# Create a non-root user
RUN groupadd -r pptruser && useradd -r -g pptruser -G audio,video pptruser \
    && mkdir -p /home/pptruser/Downloads \
    && chown -R pptruser:pptruser /home/pptruser \
    && chown -R pptruser:pptruser /app

# Copy the Go binary from builder stage
COPY --from=go-builder /app/http-server ./

# Copy Node.js scripts and package.json
COPY package.json ./
COPY scripts/ ./scripts/

# Copy public files (HTML, CSS, JS)
COPY public/ ./public/

# Install Node.js dependencies
RUN npm install

# Change ownership to pptruser
RUN chown -R pptruser:pptruser /app

# Switch to non-root user
USER pptruser

# Expose the port the app runs on
EXPOSE 8080

# Run the HTTP server
CMD ["./http-server"]
