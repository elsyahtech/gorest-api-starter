package main

import (
	"errors"
	"fmt"
	golog "log"
	"os"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/config"
	"github.com/elsyahtech/gorest-api-starter/app/router"
	"github.com/elsyahtech/gorest/database"
	"github.com/elsyahtech/gorest/filter"
	"github.com/elsyahtech/gorest/log"
	"github.com/elsyahtech/gorest/monitoring"
	"github.com/elsyahtech/gorest/redis"
	"github.com/elsyahtech/gorest/server"
	"github.com/joho/godotenv"
)

func main() {
	if err := loadDotEnv(); err != nil {
		failConfiguration(err)
	}

	// =====================================================================
	// 1. Set app configuration
	// =====================================================================
	appConfig, err := config.AppConfig()
	if err != nil {
		failConfiguration(err)
	}

	// Logger configuration
	loggerConfig, err := config.LoggerConfig()
	if err != nil {
		failConfiguration(err)
	}

	// Server configuration
	serverConfig, err := config.ServerConfig()
	if err != nil {
		failConfiguration(err)
	}

	// Database configuration
	databaseConfig, err := config.DatabaseConfig()
	if err != nil {
		failConfiguration(err)
	}

	// Redis configuration
	redisConfig, err := config.RedisConfig()
	if err != nil {
		failConfiguration(err)
	}

	// Filter (JWT) configuration
	filterConfig, err := config.FilterConfig()
	if err != nil {
		failConfiguration(err)
	}

	// Monitoring configuration
	monitoringConfig, err := config.MonitoringConfig()
	if err != nil {
		failConfiguration(err)
	}

	// =====================================================================
	// 2. Initialize application
	// =====================================================================
	app := gorest.New(*appConfig)

	// =====================================================================
	// 3. Register application services
	// =====================================================================
	app.Use(
		log.Run(*loggerConfig),
		server.Run(*serverConfig),
		database.Run(*databaseConfig),
		redis.Run(*redisConfig),
		filter.Run(*filterConfig),
		monitoring.Run(*monitoringConfig),
	)

	// =====================================================================
	// 4. Register the application's route modules.
	// =====================================================================
	app.RegisterServices([]gorest.RouterRegistrar{
		router.UserRouter,

		// Add your route modules here.
	})

	// =====================================================================
	// 5. Start the application.
	// =====================================================================
	app.Start()
}

func loadDotEnv() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load .env: %w", err)
	}

	return nil
}

func failConfiguration(err error) {
	golog.Fatalf("invalid application configuration: %v", err)
}
