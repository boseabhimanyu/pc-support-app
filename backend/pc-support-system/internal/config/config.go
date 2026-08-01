package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoUri       string
	MongoDB        string
	ServerPort     string
	JWTSecret      string
	JWTExpiryHours int
	GinMode        string
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

	jwtSecret, err := extractEnv("JWT_SECRET")
	if err != nil {
		return Config{}, err
	}

	jwtExpiryHoursStr, err := extractEnv("JWT_EXPIRY_HOURS")
	if err != nil {
		return Config{}, err
	}

	jwtExpiryHours, err := strconv.Atoi(jwtExpiryHoursStr)
	if err != nil {
		return Config{}, fmt.Errorf("invalid JWT_EXPIRY_HOURS")
	}

	gin_mode, err := extractEnv("GIN_MODE")
	if err != nil {
		return Config{}, err
	}
	return Config{
		MongoUri:       mongoURI,
		MongoDB:        mongoDB,
		ServerPort:     port,
		JWTSecret:      jwtSecret,
		JWTExpiryHours: jwtExpiryHours,
		GinMode:        gin_mode,
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

	if c.JWTSecret == "" {
		return errors.New("jwt secret missing")
	}

	if strconv.Itoa(c.JWTExpiryHours) == "" {
		return errors.New("jwt expiry hours missing")
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
