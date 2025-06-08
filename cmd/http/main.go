package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/solrac97gr/telegram-followers-checker/app"
	handlers "github.com/solrac97gr/telegram-followers-checker/cmd/http/handlers"
	"github.com/solrac97gr/telegram-followers-checker/cmd/http/middleware"
	"github.com/solrac97gr/telegram-followers-checker/config"
	"github.com/solrac97gr/telegram-followers-checker/database"
	"github.com/solrac97gr/telegram-followers-checker/extractors/instagram"
	"github.com/solrac97gr/telegram-followers-checker/extractors/rutube"
	"github.com/solrac97gr/telegram-followers-checker/extractors/telegram"
	"github.com/solrac97gr/telegram-followers-checker/extractors/vk"
	"github.com/solrac97gr/telegram-followers-checker/filemanager"
)

const (
	EchoVar = "ECHO_VAR"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Error creating config: %v", err)
	}
	log.Print("Environment variables loaded successfully echo var [Ping]:", os.Getenv(EchoVar))
	fiberApp := fiber.New()
	fiberApp.Use(logger.New())

	// Initialize MongoDB client

	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	if err := database.InitializeDatabase(mongoClient, config); err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	// Initialize components
	repo, err := database.NewMongoRepository(mongoClient, config)
	if err != nil {
		log.Fatalf("Error creating MongoDB repository: %v", err)
	}
	repo.DeleteExpiredAnalyses()

	fm := filemanager.NewFileManager()
	telegramExtractor := telegram.NewTelegramExtractor()
	rutubeExtractor := rutube.NewRutubeExtractor()
	vkExtractor := vk.NewVKExtractor()
	instagramExtractor := instagram.NewInstagramExtractor()

	influencersApp := app.NewInfluencerApp(repo, fm, telegramExtractor, rutubeExtractor, vkExtractor, instagramExtractor)

	userRepo, err := database.NewUserMongoRepository(mongoClient, config)
	if err != nil {
		log.Fatalf("Error creating user MongoDB repository: %v", err)
	}
	usersApp := app.NewUserApp(userRepo, config.JWTSecret)

	hdl, err := handlers.NewHandlers(influencersApp, usersApp)
	if err != nil {
		log.Fatalf("Error creating handlers: %v", err)
	}

	auth, err := middleware.NewAuthMiddleware(&middleware.JWTConfig{
		Secret: config.JWTSecret,
	}, userRepo)
	if err != nil {
		log.Fatalf("Error creating middlewares: %v", err)
	}

	// Serve static files from the root directory
	fiberApp.Static("/", "./public")

	// Add route for dashboard protection
	fiberApp.Get("/dashboard.html", func(c *fiber.Ctx) error {
		return c.SendFile("./public/dashboard.html")
	})

	// Add route for database protection
	fiberApp.Get("/database.html", func(c *fiber.Ctx) error {
		return c.SendFile("./public/database.html")
	})

	// Public routes (no JWT required)
	publicGroup := fiberApp.Group("/api/v1")
	publicUserHandlers := publicGroup.Group("/users")
	publicUserHandlers.Post("/register", hdl.RegisterUserHandler)
	publicUserHandlers.Post("/login", hdl.LoginUserHandler)

	// Protected routes (JWT required)
	apiv1Group := fiberApp.Group("/api/v1", auth.WithJWT())

	// User routes that require authentication
	userHandlers := apiv1Group.Group("/users")
	userHandlers.Post("/logout", hdl.LogoutUserHandler)

	// Influencer routes
	influencersHandlers := apiv1Group.Group("/influencers")
	influencersHandlers.Get("/health", hdl.HealthCheckHandler)
	influencersHandlers.Post("/upload", hdl.UploadHandler)
	influencersHandlers.Get("/download", hdl.DownloadHandler)
	influencersHandlers.Post("/estimate-time", hdl.EstimateTimeHandler)
	influencersHandlers.Get("/analyses", hdl.AnalysesHandler)

	errors := make(chan error, 3)
	go func() {
		log.Println("Starting ticker for deleting expired analyses...")
		TickerDeleteExpiredAnalyses(repo)
	}()
	go func() {
		log.Println("Starting Fiber server...")
		if err := fiberApp.Listen(":3000"); err != nil {
			errors <- err
		}
	}()
	go func() {
		log.Println("Starting ticker for deleting expired tokens...")
		TickerDeleteExpiredTokens(userRepo)
	}()

	for err := range errors {
		log.Printf("Error occurred: %v", err)
	}
	log.Println("Shutting down the server gracefully...")
	if err := mongoClient.Disconnect(context.Background()); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
	} else {
		log.Println("Disconnected from MongoDB successfully")
	}
	if err := fiberApp.Shutdown(); err != nil {
		log.Printf("Error shutting down Fiber app: %v", err)
	} else {
		log.Println("Fiber app shut down successfully")
	}
	log.Println("Server stopped gracefully")
}

func TickerDeleteExpiredAnalyses(repo database.InfluencerRepository) {
	ticker := time.NewTicker(24 * time.Hour) // Adjust the interval as needed
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Deleting expired analyses...")
		err := repo.DeleteExpiredAnalyses()
		if err != nil {
			log.Printf("Error deleting expired analyses: %v", err)
		} else {
			log.Println("Expired analyses deleted successfully")
		}
	}
}

func TickerDeleteExpiredTokens(repo database.UserRepository) {
	ticker := time.NewTicker(24 * time.Hour) // Adjust the interval as needed
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Deleting expired tokens...")
		err := repo.DeleteExpiredTokens()
		if err != nil {
			log.Printf("Error deleting expired tokens: %v", err)
		} else {
			log.Println("Expired tokens deleted successfully")
		}
	}
}
