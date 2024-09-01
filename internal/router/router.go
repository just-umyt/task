package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/just-umyt/task/database"
	"github.com/just-umyt/task/internal/handlers"
	"github.com/just-umyt/task/internal/repository"
)

func SetupRouter(app *fiber.App) {
	userRepo := repository.NewUserRepository(database.DB)
	authHandler := handlers.NewAuthHandler(userRepo)
	tokenHandler := handlers.NewTokenHandler(userRepo)

	app.Post("/signup", authHandler.UserSignUp)
	app.Post("/signin", authHandler.UserSignIn)
	app.Post("/refresh", tokenHandler.RefreshToken)
}
