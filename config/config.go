package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MongoURI          string `envconfig:"MONGO_URI"`
	JWTSecret         string `envconfig:"JWT_SECRET"`
	InfluencersDBName string `envconfig:"INFLUENCERS_DB"`
	UsersDBName       string `envconfig:"USERS_DB"`
	InstagramUsername string `envconfig:"INSTAGRAM_USERNAME"`
	InstagramPassword string `envconfig:"INSTAGRAM_PASSWORD"`
}

func NewConfig() (*Config, error) {
	var cfg Config
	// Try to load .env file, but don't fail if it doesn't exist (for Docker)
	_ = godotenv.Load()
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("error processing environment variables: %w", err)
	}
	return &cfg, nil
}
