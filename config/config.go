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
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("error processing environment variables: %w", err)
	}
	return &cfg, nil
}
