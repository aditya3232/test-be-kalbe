package application

import (
	"math"
	"test-be-kalbe/internal/domain/model"
	"test-be-kalbe/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type DepartmentApplication struct {
	DepartmentService service.DepartmentService
	Log               *logrus.Logger
}

func NewDepartmentApplication(departmentService service.DepartmentService, log *logrus.Logger) *DepartmentApplication {
	return &DepartmentApplication{
		DepartmentService: departmentService,
		Log:               log,
	}
}

func (a *DepartmentApplication) Create(ctx *fiber.Ctx) error {
	request := new(model.DepartmentCreateRequest)
	if err := ctx.BodyParser(request); err != nil {
		a.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	request.CreatedBy = ctx.Locals("employee_name").(string)

	response, err := a.DepartmentService.Create(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error creating department")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.DepartmentResponse]{
		Meta: model.Meta{Code: 200, Status: "success create department"},
		Data: response,
	})
}

func (a *DepartmentApplication) List(ctx *fiber.Ctx) error {
	request := &model.DepartmentSearchRequest{
		DepartmentName: ctx.Query("department_name", ""),
		Page:           ctx.QueryInt("page", 1),
		Size:           ctx.QueryInt("size", 10),
	}

	responses, total, err := a.DepartmentService.Search(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error searching department")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.DepartmentResponse]{
		Meta:   model.Meta{Code: 200, Status: "OK"},
		Data:   responses,
		Paging: paging,
	})
}

func (a *DepartmentApplication) Get(ctx *fiber.Ctx) error {
	request := &model.DepartmentGetByIdRequest{
		DepartmentId: ctx.Params("departmentId"),
	}

	response, err := a.DepartmentService.FindById(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error getting department")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.DepartmentResponse]{
		Meta: model.Meta{Code: 200, Status: "OK"},
		Data: response,
	})
}

func (a *DepartmentApplication) Update(ctx *fiber.Ctx) error {
	request := new(model.DepartmentUpdateRequest)
	if err := ctx.BodyParser(request); err != nil {
		a.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	request.DepartmentId = ctx.Params("departmentId")

	response, err := a.DepartmentService.Update(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error updating department")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.DepartmentResponse]{
		Meta: model.Meta{Code: 200, Status: "success update department"},
		Data: response,
	})
}

func (a *DepartmentApplication) SoftDelete(ctx *fiber.Ctx) error {
	request := &model.DepartmentDeleteRequest{
		DepartmentId: ctx.Params("departmentId"),
	}

	if err := a.DepartmentService.SoftDelete(ctx.UserContext(), request); err != nil {
		a.Log.WithError(err).Error("error deleting department")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{
		Meta: model.Meta{Code: 200, Status: "success delete department"},
		Data: true,
	})
}
