package middleware

import (
	"strings"
	"test-be-kalbe/internal/helper"

	"github.com/gofiber/fiber/v2"
)

// JWTMiddleware is a middleware to protect routes
func JWTMiddleware(ctx *fiber.Ctx) error {
	// Get the JWT token from the request header
	authHeader := ctx.Get("Authorization")
	tokenString := strings.Split(authHeader, " ")[1]

	// Parse the JWT token
	_, err := helper.ValidateToken(tokenString)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Extract the user ID from the token
	employeeId, err := helper.GetEmployeeIDFromToken(tokenString)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Set the employee ID in the context
	ctx.Locals("employee_id", employeeId)

	// Continue with the next handler
	return ctx.Next()
}
