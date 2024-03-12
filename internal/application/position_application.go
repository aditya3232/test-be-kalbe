package application

import (
	"math"
	"test-be-kalbe/internal/domain/model"
	"test-be-kalbe/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type PositionApplication struct {
	PositionService service.PositionService
	Log             *logrus.Logger
}

func NewPositionApplication(positionService service.PositionService, log *logrus.Logger) *PositionApplication {
	return &PositionApplication{
		PositionService: positionService,
		Log:             log,
	}
}

func (a *PositionApplication) Create(ctx *fiber.Ctx) error {
	request := new(model.PositionCreateRequest)
	if err := ctx.BodyParser(request); err != nil {
		a.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	response, err := a.PositionService.Create(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error creating position")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.PositionResponse]{
		Meta: model.Meta{Code: 200, Status: "success create position"},
		Data: response,
	})
}

func (a *PositionApplication) List(ctx *fiber.Ctx) error {
	request := &model.PositionSearchRequest{
		PositionName: ctx.Query("position_name", ""),
		Page:         ctx.QueryInt("page", 1),
		Size:         ctx.QueryInt("size", 10),
	}

	responses, total, err := a.PositionService.Search(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error searching position")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.PositionResponse]{
		Meta:   model.Meta{Code: 200, Status: "OK"},
		Data:   responses,
		Paging: paging,
	})
}

func (a *PositionApplication) Get(ctx *fiber.Ctx) error {
	request := &model.PositionGetByIdRequest{
		PositionId: ctx.Params("positionId"),
	}

	response, err := a.PositionService.FindById(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error getting position")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.PositionResponse]{
		Meta: model.Meta{Code: 200, Status: "OK"},
		Data: response,
	})
}

func (a *PositionApplication) Update(ctx *fiber.Ctx) error {
	request := new(model.PositionUpdateRequest)
	if err := ctx.BodyParser(request); err != nil {
		a.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	request.PositionId = ctx.Params("positionId")

	response, err := a.PositionService.Update(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error updating position")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.PositionResponse]{
		Meta: model.Meta{Code: 200, Status: "success update position"},
		Data: response,
	})
}

func (a *PositionApplication) SoftDelete(ctx *fiber.Ctx) error {
	request := &model.PositionDeleteRequest{
		PositionId: ctx.Params("positionId"),
	}

	if err := a.PositionService.SoftDelete(ctx.UserContext(), request); err != nil {
		a.Log.WithError(err).Error("error deleting position")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{
		Meta: model.Meta{Code: 200, Status: "success delete position"},
		Data: true,
	})
}
