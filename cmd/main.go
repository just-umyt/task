package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/just-umyt/task/database"
	"github.com/just-umyt/task/internal/configs"
	"github.com/just-umyt/task/internal/router"
)

func main() {
	//Load env file
	configs.LoadEnv()

	//Connect to database
	database.ConnectDB()

	//Create new fiber App
	app := fiber.New(fiber.Config{
		AppName: os.Getenv("APP_NAME"),
	})

	//Add deafult logger
	app.Use(logger.New())

	//Setup routes
	router.SetupRouter(app)

	fmt.Println("gee")

	//Listen and serve
	app.Listen(os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT"))

}
