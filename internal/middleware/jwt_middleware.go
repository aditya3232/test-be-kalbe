package middleware

import (
	"strings"
	"test-be-kalbe/internal/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type JwtApplication struct {
	Log *logrus.Logger
}

func NewJwtApplication(log *logrus.Logger) *JwtApplication {
	return &JwtApplication{
		Log: log,
	}
}

// JWTMiddleware is a middleware to protect routes
func (a *JwtApplication) JWTMiddleware(ctx *fiber.Ctx) error {
	// Get the JWT token from the request header
	authHeader := ctx.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse the JWT token
	_, err := helper.ValidateToken(tokenString)
	if err != nil {
		a.Log.WithError(err).Error("error validating token")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Extract the user ID from the token
	employeeId, err := helper.GetEmployeeIDFromToken(tokenString)
	if err != nil {
		a.Log.WithError(err).Error("error getting employee ID from token")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Set the employee ID in the context
	ctx.Locals("employee_id", employeeId)

	// Continue with the next handler
	return ctx.Next()
}
