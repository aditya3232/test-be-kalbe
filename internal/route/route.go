package route

import (
	"fmt"
	"log"
	"test-be-kalbe/internal/application"
	"test-be-kalbe/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                   *fiber.App
	DepartmentApplication *application.DepartmentApplication
	PositionApplication   *application.PositionApplication
	EmployeeApplication   *application.EmployeeApplication
	AuthApplication       *application.AuthApplication
	MiddlewareApplication *middleware.JwtApplication
	LocationApplication   *application.LocationApplication
	AttendanceApplication *application.AttendanceApplication
}

func (r *RouteConfig) Setup() {
	r.App.Use(recoverPanic)
	r.SetupAuthenticationRoute()
	r.SetupMiddlewareRoute()
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

func (r *RouteConfig) SetupAuthenticationRoute() {
	authenticationGroup := r.App.Group("/api")

	authenticationGroup.Post("/login", r.AuthApplication.Login)
	authenticationGroup.Post("/logout", r.AuthApplication.Logout)

}

// SetupMiddlewareRoute
func (r *RouteConfig) SetupMiddlewareRoute() {
	middlewareGroup := r.App.Group("/api")
	middlewareGroup.Use(r.MiddlewareApplication.JWTMiddleware)

	middlewareGroup.Get("/departments", r.DepartmentApplication.List)
	middlewareGroup.Get("/departments", r.DepartmentApplication.List)
	middlewareGroup.Post("/department", r.DepartmentApplication.Create)
	middlewareGroup.Put("/department/:departmentId", r.DepartmentApplication.Update)
	middlewareGroup.Get("/department/:departmentId", r.DepartmentApplication.Get)
	middlewareGroup.Delete("/department/:departmentId", r.DepartmentApplication.SoftDelete)

	middlewareGroup.Get("/positions", r.PositionApplication.List)
	middlewareGroup.Post("/position", r.PositionApplication.Create)
	middlewareGroup.Put("/position/:positionId", r.PositionApplication.Update)
	middlewareGroup.Get("/position/:positionId", r.PositionApplication.Get)
	middlewareGroup.Delete("/position/:positionId", r.PositionApplication.SoftDelete)

	middlewareGroup.Get("/employees", r.EmployeeApplication.List)
	middlewareGroup.Post("/employee", r.EmployeeApplication.Create)
	middlewareGroup.Put("/employee/:employeeId", r.EmployeeApplication.Update)
	middlewareGroup.Get("/employee/:employeeId", r.EmployeeApplication.Get)
	middlewareGroup.Delete("/employee/:employeeId", r.EmployeeApplication.SoftDelete)

	middlewareGroup.Get("/locations", r.LocationApplication.List)
	middlewareGroup.Post("/location", r.LocationApplication.Create)
	middlewareGroup.Put("/location/:locationId", r.LocationApplication.Update)
	middlewareGroup.Get("/location/:locationId", r.LocationApplication.Get)
	middlewareGroup.Delete("/location/:locationId", r.LocationApplication.SoftDelete)

	middlewareGroup.Get("/attendances", r.AttendanceApplication.List)
	middlewareGroup.Post("/attendance", r.AttendanceApplication.Create)
	middlewareGroup.Put("/attendance/:attendanceId", r.AttendanceApplication.Update)
	middlewareGroup.Get("/attendance/:attendanceId", r.AttendanceApplication.Get)
	middlewareGroup.Delete("/attendance/:attendanceId", r.AttendanceApplication.SoftDelete)

}
