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

}
