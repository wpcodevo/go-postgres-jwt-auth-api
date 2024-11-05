package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wpcodevo/go-postgres-jwt-auth-api/initializers"
	"github.com/wpcodevo/go-postgres-jwt-auth-api/models"
)

func GetMeHandler(c *fiber.Ctx) error {
	user := c.Locals("user").(models.UserResponse)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "data": fiber.Map{"user": user}})
}

func GetUsersHandler(c *fiber.Ctx) error {
	var page = c.Query("page", "1")
	var limit = c.Query("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var users []models.User
	results := initializers.DB.Limit(intLimit).Offset(offset).Find(&users)
	if results.Error != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"status": "error", "message": results.Error})
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, models.UserResponse{
			ID:        *user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      *user.Role,
			Photo:     *user.Photo,
			Provider:  *user.Provider,
			CreatedAt: *user.CreatedAt,
			UpdatedAt: *user.UpdatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "results": len(userResponses), "users": userResponses})
}
