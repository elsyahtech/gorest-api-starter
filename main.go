package main

import (
	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/config"
	"github.com/elsyahtech/gorest-api-starter/app/router"
	"github.com/elsyahtech/gorest/database"
	"github.com/elsyahtech/gorest/log"
	"github.com/elsyahtech/gorest/redis"
	"github.com/elsyahtech/gorest/server"
)

func main() {
	// Initialize the application and its shared configuration.
	app := gorest.New(gorest.Config{
		AppName:     "gorest-api-starter",
		Version:     "1.0.0",
		Environment: "Development",
		Timezone:    "Asia/Jakarta",
		Debug:       false,
	})

	// Register application services before registering routes.
	app.Use(
		log.Run(config.Logger()),
		redis.Run(config.Redis()),
		database.Run(config.Database()),
		server.Run(config.Server()),
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
