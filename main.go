package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/config"
	"github.com/elsyahtech/gorest-api-starter/app/router"
	"github.com/elsyahtech/gorest/database"
	"github.com/elsyahtech/gorest/log"
	"github.com/elsyahtech/gorest/redis"
	"github.com/elsyahtech/gorest/server"
)

func main() {
	databaseConfig, err := config.LoadDatabase()
	if err != nil {
		failConfiguration(err)
	}
	serverConfig, err := config.LoadServer()
	if err != nil {
		failConfiguration(err)
	}
	redisConfig, err := config.LoadRedis()
	if err != nil {
		failConfiguration(err)
	}
	if config.Environment() == "production" {
		if err := requireProductionRedis(redisConfig); err != nil {
			failConfiguration(err)
		}
	}

	// Initialize the application and its shared configuration.
	app := gorest.New(gorest.Config{
		AppName:     "gorest-api-starter",
		Version:     "1.0.0",
		Environment: config.Environment(),
		Timezone:    "Asia/Jakarta",
		Debug:       false,
	})

	// Register application services before registering routes.
	app.Use(
		log.Run(config.Logger()),
		redis.Run(redisConfig),
		database.Run(databaseConfig),
		server.Run(serverConfig),
	)

	// Register the application's route modules.
	app.RegisterServices([]gorest.RouterRegistrar{
		router.UserRouter,
		router.DepartmentRouter,
		router.OrderRouter,
	})

	// Start serves requests and manages the application shutdown lifecycle.
	app.Start()
}

func requireProductionRedis(cfg redis.Config) error {
	if !cfg.Enabled {
		return nil
	}
	for _, key := range []string{"GOREST_REDIS_HOST", "GOREST_REDIS_PASSWORD"} {
		if strings.TrimSpace(os.Getenv(key)) == "" {
			return fmt.Errorf("%s must be set when Redis is enabled in production", key)
		}
	}
	if !cfg.TLS {
		return fmt.Errorf("GOREST_REDIS_TLS must be true when Redis is enabled in production")
	}
	return nil
}

func failConfiguration(err error) {
	fmt.Fprintf(os.Stderr, "invalid application configuration: %v\n", err)
	os.Exit(1)
}
