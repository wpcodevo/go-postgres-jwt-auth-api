package main

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/wpcodevo/go-postgres-jwt-auth-api/initializers"
	"github.com/wpcodevo/go-postgres-jwt-auth-api/models"
)

func DeserializeUser(c *fiber.Ctx) error {
	var tokenString string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("token") != "" {
		tokenString = c.Cookies("token")
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
	}

	config, _ := initializers.LoadEnv(".")

	tokenByte, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}

		return []byte(config.JwtSecret), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("invalidate token: %v", err)})
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "invalid token claim"})

	}

	var user models.User
	initializers.DB.First(&user, "id = ?", fmt.Sprint(claims["sub"]))

	if user.ID.String() != claims["sub"] {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no logger exists"})
	}

	c.Locals("user", models.FilterUserRecord(&user))

	return c.Next()
}

func allowedRoles(allowedRoles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(models.UserResponse)

		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "fail",
				"message": "Access denied. User not authenticated.",
			})
		}

		roleAllowed := false
		for _, allowedRole := range allowedRoles {
			if user.Role == allowedRole {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "fail",
				"message": "Access denied. You are not allowed to perform this action.",
			})
		}

		return c.Next()
	}
}
