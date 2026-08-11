package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoUri       string
	MongoDB        string
	ServerPort     string
	JWTSecret      string
	JWTExpiryHours int
	GinMode        string
	AllowedOrigins []string
}

func Load() (Config, error) {
	// .env is optional.
	//
	// Local development:
	//   .env can be loaded by godotenv.
	//
	// Docker/cloud:
	//   Environment variables can be supplied externally.
	//
	// Therefore, failure to find .env should NOT be an error.
	paths := []string{
		".env",
		"../.env",
		"../../.env",
	}

	for _, p := range paths {
		if err := godotenv.Load(p); err == nil {
			break
		}
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
		return Config{}, fmt.Errorf("invalid JWT_EXPIRY_HOURS: %w", err)
	}

	ginMode, err := extractEnv("GIN_MODE")
	if err != nil {
		return Config{}, err
	}

	allowedOriginsStr, err := extractEnv("ALLOWED_ORIGINS")
	if err != nil {
		return Config{}, err
	}

	// Convert:
	//
	// http://localhost:3000,http://localhost:5173
	//
	// into:
	//
	// []string{
	//     "http://localhost:3000",
	//     "http://localhost:5173",
	// }
	allowedOrigins := strings.Split(allowedOriginsStr, ",")

	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}

	config := Config{
		MongoUri:       mongoURI,
		MongoDB:        mongoDB,
		ServerPort:     port,
		JWTSecret:      jwtSecret,
		JWTExpiryHours: jwtExpiryHours,
		GinMode:        ginMode,
		AllowedOrigins: allowedOrigins,
	}

	if err := config.Validate(); err != nil {
		return Config{}, err
	}

	return config, nil
}

// Validate validates the loaded configuration.
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

	if c.JWTExpiryHours <= 0 {
		return errors.New("JWT expiry hours must be greater than zero")
	}

	if c.GinMode == "" {
		return errors.New("GIN_MODE missing")
	}

	if len(c.AllowedOrigins) == 0 {
		return errors.New("ALLOWED_ORIGINS missing")
	}

	return nil
}

func extractEnv(key string) (string, error) {
	val := strings.TrimSpace(os.Getenv(key))

	if val == "" {
		return "", fmt.Errorf(
			"missing required environment variable: %s",
			key,
		)
	}

	return val, nil
}
