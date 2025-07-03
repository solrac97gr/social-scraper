# Gemini Code Assistant Project Overview

This document provides a comprehensive overview of the "Social Scraper" project, designed to assist Gemini in understanding and interacting with the codebase.

## Project Description

"Social Scraper" is a Go-based application that extracts follower counts and channel names from Telegram, Rutube, VK, and Instagram. It takes an Excel file containing a list of channel URLs as input and outputs a new Excel file with the scraped data. The project includes both a command-line interface (CLI) for direct execution and an HTTP server for a web-based user interface.

## Key Technologies

- **Go**: The primary programming language for the backend logic.
- **Node.js**: Used for running the Puppeteer script for scraping.
- **MongoDB**: The project is configured to use MongoDB as a database, as indicated by the `infracli run mongo` command in the Makefile.
- **Puppeteer**: A Node.js library used for headless browser scraping of Instagram and VK.
- **Go Libraries**:
    - `github.com/PuerkitoBio/goquery`: For parsing HTML.
    - `github.com/xuri/excelize/v2`: For reading and writing Excel files.
    - `github.com/gofiber/fiber/v2`: For the HTTP server.
    - `github.com/golang-jwt/jwt/v5`: For JSON Web Token (JWT) authentication.
    - `go.mongodb.org/mongo-driver`: For interacting with MongoDB.

## Project Structure

The project is organized into several packages:

- `cmd/`: Contains the main entry points for the CLI (`cmd/cli/main.go`) and the HTTP server (`cmd/http/main.go`).
- `app/`: Defines the core application logic and data structures for influencers and users.
- `config/`: Handles application configuration.
- `database/`: Manages database connections and data persistence, with specific implementations for MongoDB.
- `extractors/`: Contains the scraping logic for different social media platforms (Telegram, Rutube, VK, Instagram).
- `filemanager/`: Handles file operations, such as reading and writing Excel files.
- `public/`: Contains the static assets for the web interface (HTML, CSS, JavaScript).
- `scripts/`: Includes the Node.js-based Puppeteer script for web scraping.

## How to Run the Project

The `Makefile` provides several commands for managing the project:

- `make deps`: Installs the necessary Go and Node.js dependencies.
- `make build`: Compiles the Go application.
- `make run`: Starts the MongoDB container and runs the HTTP server.
- `make cli`: Starts the MongoDB container and runs the CLI application.
- `make clean`: Removes generated files and directories.

### Running the CLI

1. Ensure MongoDB is running.
2. Execute the following command:
   ```sh
   go run cmd/cli/main.go /path/to/your_excel_file.xlsx
   ```

### Running the HTTP Server

1. Ensure MongoDB is running.
2. Execute the following command:
   ```sh
   go run cmd/http/main.go
   ```
3. Open a web browser and navigate to `http://localhost:3000`.

## Authentication

The HTTP server uses JWT for authentication, with endpoints for user registration and login. The `middleware/jwt.go` file contains the JWT authentication middleware.

## Caching and State

The application does not appear to have an explicit caching layer beyond in-memory data structures. State is primarily managed through the MongoDB database.
