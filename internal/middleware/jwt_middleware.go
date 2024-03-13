package middleware

import (
	"strconv"
	"test-be-kalbe/internal/domain/model"
	"test-be-kalbe/internal/helper"
	"test-be-kalbe/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type JwtApplication struct {
	EmployeeService service.EmployeeService
	Log             *logrus.Logger
}

func NewJwtApplication(employeeService service.EmployeeService, log *logrus.Logger) *JwtApplication {
	return &JwtApplication{
		EmployeeService: employeeService,
		Log:             log,
	}
}

// JWTMiddleware is a middleware to protect routes
func (a *JwtApplication) JWTMiddleware(ctx *fiber.Ctx) error {
	// Get the JWT token from the request header
	authHeader := ctx.Get("Authorization")

	// Parse the JWT token
	_, err := helper.ValidateToken(authHeader)
	if err != nil {
		a.Log.WithError(err).Error("error validating token")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Extract the user ID from the token
	employeeId, err := helper.GetEmployeeIDFromToken(authHeader)
	if err != nil {
		a.Log.WithError(err).Error("error getting employee ID from token")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// get employee by id
	formatEmployeeId := strconv.Itoa(employeeId)
	employee, err := a.EmployeeService.FindById(ctx.UserContext(), &model.EmployeeGetByIdRequest{EmployeeId: formatEmployeeId})
	if err != nil {
		a.Log.WithError(err).Error("error finding employee by id")
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	if employee.Token != authHeader {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
			// "database":    employee.Token,
			// "request":     authHeader,
			// "allEmployee": employee,
		})
	}

	// Set the employee ID in the context
	ctx.Locals("employee_id", employeeId)

	// Continue with the next handler
	return ctx.Next()
}
