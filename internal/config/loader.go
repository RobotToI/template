package config

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

var config *Common

// Initialize configuration by .env files
func Initialize(ctx context.Context) error {
	if err := initialize(ctx); err != nil {
		return err
	}
	return nil
}

func initialize(ctx context.Context) error {
	if err := loadEnvFiles(ctx); err != nil {
		return err
	}

	config = &Common{}

	return envconfig.Process("APP", config)
}

// Get returns common app configuration
func Get() Common {
	if config == nil {
		panic("Configuration is not initialized")
	}

	return *config
}

var loadEnvFiles = func(ctx context.Context) (err error) {
	logger := log.Ctx(ctx)
	env := os.Getenv(appEnv)
	if env == "" {
		env = devEnv
		_ = os.Setenv(appEnv, env)
		logger.Info().Msgf("Environment variable %q is not defined. Using %q environment\n", appEnv, env)
	} else {
		logger.Info().Msgf("Using %v environment", env)
	}

	if err = loadEnvFileIfExists(".env." + env + ".local"); err != nil {
		return err
	}

	if env != testEnv {
		if err = loadEnvFileIfExists(".env.local"); err != nil {
			return err
		}
	}

	if err = loadEnvFileIfExists(".env." + env); err != nil {
		return err
	}

	return loadEnvFileIfExists(".env")
}

func loadEnvFileIfExists(file string) error {
	_, err := os.Stat(file)

	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return err
	}

	if err = godotenv.Load(file); err != nil {
		return fmt.Errorf("failed to load %s: %s", file, err)
	}

	return err
}
