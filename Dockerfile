# Stage 1: Build the Go application and install Node.js dependencies
FROM golang:1.22-alpine AS builder

# Install build tools and Node.js
RUN apk add --no-cache git build-base nodejs npm

# Set up the working directory
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy Node.js package files and install dependencies
# We assume package.json is inside the scripts directory
COPY scripts/package.json scripts/package-lock.json* ./scripts/
RUN npm install --prefix ./scripts

# Copy the rest of the application source code
COPY . .

# Build the Go application
# This will create a static binary named 'social-scraper' in the /app directory
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /social-scraper cmd/http/main.go

# Stage 2: Create the final, lean production image
FROM alpine:latest

# Set up the working directory
WORKDIR /app

# Install Node.js and Chromium for Puppeteer
# These are the minimal dependencies required to run headless Chrome on Alpine
RUN apk add --no-cache \
    nodejs \
    npm \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont

# Set the environment variable for Puppeteer to find the installed Chromium
ENV PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium-browser

# Copy the built Go binary from the builder stage
COPY --from=builder /social-scraper .

# Copy the public assets and scripts (with node_modules)
COPY --from=builder /app/public ./public
COPY --from=builder /app/scripts ./scripts

# Expose the port the application runs on
EXPOSE 3000

# The command to run the application
CMD ["/app/social-scraper"]
