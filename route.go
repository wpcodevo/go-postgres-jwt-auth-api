package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wpcodevo/go-postgres-jwt-auth-api/handlers"
)

func SetupRoutes(app fiber.Router) {
	authRoutes := app.Group("/auth")
	authRoutes.Post("/register", handlers.SignUpUser)
	authRoutes.Post("/login", handlers.SignInUser)
	authRoutes.Get("/logout", DeserializeUser, handlers.LogoutUser)

	app.Get("/users/me", DeserializeUser, handlers.GetMeHandler)
	app.Get("/users/", DeserializeUser, allowedRoles([]string{"admin", "moderator"}), handlers.GetUsersHandler)
}
