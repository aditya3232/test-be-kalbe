package application

import (
	"math"
	"test-be-kalbe/internal/domain/model"
	"test-be-kalbe/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type DepartmentApplication struct {
	Service service.DepartmentService
	Log     *logrus.Logger
}

func NewDepartmentApplication(service *service.DepartmentService, log *logrus.Logger) *DepartmentApplication {
	return &DepartmentApplication{
		Service: *service,
		Log:     log,
	}
}

func (a *DepartmentApplication) Create(ctx *fiber.Ctx) error {
	request := new(model.DepartmentCreateRequest)
	if err := ctx.BodyParser(request); err != nil {
		a.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	response, err := a.Service.Create(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error creating customer")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.DepartmentResponse]{Data: response})
}

func (a *DepartmentApplication) List(ctx *fiber.Ctx) error {

	request := &model.DepartmentSearchRequest{
		DepartmentName: ctx.Query("department_name", ""),
		Page:           ctx.QueryInt("page", 1),
		Size:           ctx.QueryInt("size", 10),
	}

	responses, total, err := a.Service.Search(ctx.UserContext(), request)
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
		Data:   responses,
		Paging: paging,
	})
}

func (a *DepartmentApplication) Get(ctx *fiber.Ctx) error {
	request := &model.DepartmentGetByIdRequest{
		DepartmentId: ctx.Params("department_id"),
	}

	response, err := a.Service.FindById(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error getting department")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.DepartmentResponse]{Data: response})
}

func (a *DepartmentApplication) Update(ctx *fiber.Ctx) error {
	request := new(model.DepartmentUpdateRequest)
	if err := ctx.BodyParser(request); err != nil {
		a.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	request.DepartmentId = ctx.Params("department_id")

	response, err := a.Service.Update(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error updating department")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.DepartmentResponse]{Data: response})
}

func (a *DepartmentApplication) Delete(ctx *fiber.Ctx) error {
	request := &model.DepartmentDeleteRequest{
		DepartmentId: ctx.Params("department_id"),
	}

	if err := a.Service.Delete(ctx.UserContext(), request); err != nil {
		a.Log.WithError(err).Error("error deleting customer")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}
