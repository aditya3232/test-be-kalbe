package application

import (
	"math"
	"test-be-kalbe/internal/service"

	"test-be-kalbe/internal/domain/model"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type EmployeeApplication struct {
	EmployeeService service.EmployeeService
	Log             *logrus.Logger
}

func NewEmployeeApplication(employeeService service.EmployeeService, log *logrus.Logger) *EmployeeApplication {
	return &EmployeeApplication{
		EmployeeService: employeeService,
		Log:             log,
	}
}

func (a *EmployeeApplication) Create(ctx *fiber.Ctx) error {
	request := new(model.EmployeeCreateRequest)
	if err := ctx.BodyParser(request); err != nil {
		a.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	response, err := a.EmployeeService.Create(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error creating employee")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.EmployeeResponse]{
		Meta: model.Meta{Code: 200, Status: "success create employee"},
		Data: response,
	})
}

func (a *EmployeeApplication) List(ctx *fiber.Ctx) error {
	request := &model.EmployeeSearchRequest{
		EmployeeName: ctx.Query("employee_name", ""),
		Page:         ctx.QueryInt("page", 1),
		Size:         ctx.QueryInt("size", 10),
	}

	responses, total, err := a.EmployeeService.Search(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error searching employee")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.EmployeeResponse]{
		Meta:   model.Meta{Code: 200, Status: "OK"},
		Data:   responses,
		Paging: paging,
	})
}

func (a *EmployeeApplication) Get(ctx *fiber.Ctx) error {
	request := &model.EmployeeGetByIdRequest{
		EmployeeId: ctx.Params("employeeId"),
	}

	response, err := a.EmployeeService.FindById(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error getting employee")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.EmployeeResponse]{
		Meta: model.Meta{Code: 200, Status: "OK"},
		Data: response,
	})
}

func (a *EmployeeApplication) Update(ctx *fiber.Ctx) error {
	request := new(model.EmployeeUpdateRequest)
	if err := ctx.BodyParser(request); err != nil {
		a.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	request.EmployeeId = ctx.Params("employeeId")

	response, err := a.EmployeeService.Update(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error updating employee")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.EmployeeResponse]{
		Meta: model.Meta{Code: 200, Status: "success update employee"},
		Data: response,
	})
}

func (a *EmployeeApplication) SoftDelete(ctx *fiber.Ctx) error {
	request := &model.EmployeeDeleteRequest{
		EmployeeId: ctx.Params("employeeId"),
	}

	if err := a.EmployeeService.SoftDelete(ctx.UserContext(), request); err != nil {
		a.Log.WithError(err).Error("error deleting employee")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{
		Meta: model.Meta{Code: 200, Status: "success delete employee"},
		Data: true,
	})
}
