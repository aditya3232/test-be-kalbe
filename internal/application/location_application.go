package application

import (
	"math"
	"test-be-kalbe/internal/domain/model"
	"test-be-kalbe/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type LocationApplication struct {
	LocationService service.LocationService
	Log             *logrus.Logger
}

func NewLocationApplication(locationService service.LocationService, log *logrus.Logger) *LocationApplication {
	return &LocationApplication{
		LocationService: locationService,
		Log:             log,
	}
}

func (a *LocationApplication) Create(ctx *fiber.Ctx) error {
	request := new(model.LocationCreateRequest)
	if err := ctx.BodyParser(request); err != nil {
		a.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	request.CreatedBy = ctx.Locals("employee_name").(string)

	response, err := a.LocationService.Create(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error creating location")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.LocationResponse]{
		Meta: model.Meta{Code: 200, Status: "success create location"},
		Data: response,
	})
}

func (a *LocationApplication) List(ctx *fiber.Ctx) error {
	request := &model.LocationSearchRequest{
		LocationName: ctx.Query("location_name", ""),
		Page:         ctx.QueryInt("page", 1),
		Size:         ctx.QueryInt("size", 10),
	}

	responses, total, err := a.LocationService.Search(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error searching location")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.LocationResponse]{
		Meta:   model.Meta{Code: 200, Status: "OK"},
		Data:   responses,
		Paging: paging,
	})
}

func (a *LocationApplication) Get(ctx *fiber.Ctx) error {
	request := &model.LocationGetByIdRequest{
		LocationId: ctx.Params("locationId"),
	}

	response, err := a.LocationService.FindById(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error getting location")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.LocationResponse]{
		Meta: model.Meta{Code: 200, Status: "OK"},
		Data: response,
	})
}

func (a *LocationApplication) Update(ctx *fiber.Ctx) error {
	request := new(model.LocationUpdateRequest)
	if err := ctx.BodyParser(request); err != nil {
		a.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	request.LocationId = ctx.Params("locationId")
	request.UpdatedBy = ctx.Locals("employee_name").(string)

	response, err := a.LocationService.Update(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error updating location")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.LocationResponse]{
		Meta: model.Meta{Code: 200, Status: "success update location"},
		Data: response,
	})
}

func (a *LocationApplication) Delete(ctx *fiber.Ctx) error {
	request := &model.LocationDeleteRequest{
		LocationId: ctx.Params("locationId"),
	}

	if err := a.LocationService.SoftDelete(ctx.UserContext(), request); err != nil {
		a.Log.WithError(err).Error("error deleting location")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{
		Meta: model.Meta{Code: 200, Status: "success delete location"},
		Data: true,
	})
}
