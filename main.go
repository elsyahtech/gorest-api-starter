package main

import (
	"fmt"
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
	if err := godotenv.Load(); err != nil {
		failConfiguration(err)
	}

	// =====================================================================
	// 1. Set app configuration
	// =====================================================================
	appConfig, err := config.AppConfig()
	if err != nil {
		failConfiguration(err)
	}

	loggerConfig, err := config.LoggerConfig()
	if err != nil {
		failConfiguration(err)
	}

	serverConfig, err := config.ServerConfig()
	if err != nil {
		failConfiguration(err)
	}

	databaseConfig, err := config.DatabaseConfig()
	if err != nil {
		failConfiguration(err)
	}

	redisConfig, err := config.RedisConfig()
	if err != nil {
		failConfiguration(err)
	}

	filterConfig, err := config.FilterConfig()
	if err != nil {
		failConfiguration(err)
	}

	monitoringConfig, err := config.MonitoringConfig()
	if err != nil {
		failConfiguration(err)
	}

	// =====================================================================
	// 2. Initialize application
	// =====================================================================
	app := gorest.New(appConfig)

	// =====================================================================
	// 3. Register application services
	// =====================================================================
	app.Use(
		log.Run(loggerConfig),
		server.Run(serverConfig),
		database.Run(databaseConfig),
		redis.Run(redisConfig),
		filter.Run(filterConfig),
		monitoring.Run(monitoringConfig),
	)

	// =====================================================================
	// 4. Register the application's route modules.
	// =====================================================================
	app.RegisterServices([]gorest.RouterRegistrar{
		router.UserRouter,
	})

	// =====================================================================
	// 5. Start the application.
	// =====================================================================
	app.Start()
}

func failConfiguration(err error) {
	fmt.Fprintf(os.Stderr, "invalid application configuration: %v\n", err)

	os.Exit(1)
}
