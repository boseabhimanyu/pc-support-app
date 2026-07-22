package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoUri   string
	MongoDB    string
	ServerPort string
}

func Load() (Config, error) {
	paths := []string{
		".env",
		"../.env",
		"../../.env",
	}

	loaded := false

	for _, p := range paths {
		if err := godotenv.Load(p); err == nil {
			loaded = true
			break
		}
	}

	if !loaded {
		return Config{}, fmt.Errorf(".env not found")
	}

	mongoURI, err := extractEnv("MONGO_URI")
	if err != nil {
		return Config{}, err
	}

	mongoDB, err := extractEnv("MONGO_DB_NAME")
	if err != nil {
		return Config{}, err
	}

	port, err := extractEnv("PORT")
	if err != nil {
		return Config{}, err
	}

	return Config{
		MongoUri:   mongoURI,
		MongoDB:    mongoDB,
		ServerPort: port,
	}, nil
}

// Add a Configuration Validation Function
func (c Config) Validate() error {

	if c.MongoUri == "" {
		return errors.New("mongo uri missing")
	}

	if c.MongoDB == "" {
		return errors.New("mongo database missing")
	}

	if c.ServerPort == "" {
		return errors.New("server port missing")
	}

	return nil
}

func extractEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("missing required environment variable: %s", key)
	}
	return val, nil
}
