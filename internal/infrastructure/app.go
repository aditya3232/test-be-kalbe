package infrastructure

import (
	"test-be-kalbe/internal/application"
	"test-be-kalbe/internal/middleware"
	"test-be-kalbe/internal/repository"
	"test-be-kalbe/internal/route"
	"test-be-kalbe/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	// setup repository
	departmentRepository := repository.NewDepartmentRepository(config.DB, config.Log)
	positionRepository := repository.NewPositionRepository(config.DB, config.Log)
	employeeRepository := repository.NewEmployeeRepository(config.DB, config.Log)
	locationRepository := repository.NewLocationRepository(config.DB, config.Log)

	// setup service
	departmentService := service.NewDepartmentService(config.DB, config.Log, config.Validate, departmentRepository)
	positionService := service.NewPositionService(config.DB, config.Log, config.Validate, positionRepository, departmentRepository)
	employeeService := service.NewEmployeeService(config.DB, config.Log, config.Validate, employeeRepository, positionRepository, departmentRepository)
	authService := service.NewAuthService(config.DB, config.Log, config.Validate, employeeRepository)
	locationService := service.NewLocationService(config.DB, config.Log, config.Validate, locationRepository)

	// setup application
	departmentApplication := application.NewDepartmentApplication(departmentService, config.Log)
	positionApplication := application.NewPositionApplication(positionService, config.Log)
	employeeApplication := application.NewEmployeeApplication(employeeService, config.Log)
	authApplication := application.NewAuthApplication(authService, config.Log)
	jwtMiddleware := middleware.NewJwtApplication(employeeService, config.Log)
	locationApplication := application.NewLocationApplication(locationService, config.Log)

	// setup route
	routeConfig := route.RouteConfig{
		App:                   config.App,
		DepartmentApplication: departmentApplication,
		PositionApplication:   positionApplication,
		EmployeeApplication:   employeeApplication,
		AuthApplication:       authApplication,
		MiddlewareApplication: jwtMiddleware,
		LocationApplication:   locationApplication,
	}
	routeConfig.Setup()
}
