package route

import (
	"fmt"
	"log"
	"test-be-kalbe/internal/application"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                   *fiber.App
	DepartmentApplication *application.DepartmentApplication
	PositionApplication   *application.PositionApplication
	EmployeeApplication   *application.EmployeeApplication
}

func (c *RouteConfig) Setup() {
	c.App.Use(recoverPanic)
	c.SetupGuestRoute()
}

func recoverPanic(ctx *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("panic occurred: %v", r)
			log.Println(err)
			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
				"error":   err.Error(),
			})
		}
	}()

	return ctx.Next()
}

func (r *RouteConfig) SetupGuestRoute() {
	r.App.Get("/api/departments", r.DepartmentApplication.List)
	r.App.Post("/api/department", r.DepartmentApplication.Create)
	r.App.Put("/api/department/:departmentId", r.DepartmentApplication.Update)
	r.App.Get("/api/department/:departmentId", r.DepartmentApplication.Get)
	r.App.Delete("/api/department/:departmentId", r.DepartmentApplication.SoftDelete)

	r.App.Get("/api/positions", r.PositionApplication.List)
	r.App.Post("/api/position", r.PositionApplication.Create)
	r.App.Put("/api/position/:positionId", r.PositionApplication.Update)
	r.App.Get("/api/position/:positionId", r.PositionApplication.Get)
	r.App.Delete("/api/position/:positionId", r.PositionApplication.SoftDelete)

	r.App.Get("/api/employees", r.EmployeeApplication.List)
	r.App.Post("/api/employee", r.EmployeeApplication.Create)
	r.App.Put("/api/employee/:employeeId", r.EmployeeApplication.Update)
	r.App.Get("/api/employee/:employeeId", r.EmployeeApplication.Get)
	r.App.Delete("/api/employee/:employeeId", r.EmployeeApplication.SoftDelete)

}
