package application

import (
	"math"
	"test-be-kalbe/internal/service"

	"test-be-kalbe/internal/domain/model"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AttendanceApplication struct {
	AttendanceService service.AttendanceService
	Log               *logrus.Logger
}

func NewAttendanceApplication(attendanceService service.AttendanceService, log *logrus.Logger) *AttendanceApplication {
	return &AttendanceApplication{
		AttendanceService: attendanceService,
		Log:               log,
	}
}

func (a *AttendanceApplication) Create(ctx *fiber.Ctx) error {
	request := new(model.AttendanceCreateRequest)
	if err := ctx.BodyParser(request); err != nil {
		a.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	request.CreatedBy = ctx.Locals("employee_name").(string)

	response, err := a.AttendanceService.Create(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error creating attendance")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AttendanceResponse]{
		Meta: model.Meta{Code: 200, Status: "success create attendance"},
		Data: response,
	})
}

func (a *AttendanceApplication) List(ctx *fiber.Ctx) error {
	request := &model.AttendanceSearchRequest{
		EmployeeId: ctx.Query("employee_id", ""),
		LocationId: ctx.Query("location_id", ""),
		Page:       ctx.QueryInt("page", 1),
		Size:       ctx.QueryInt("size", 10),
	}

	responses, total, err := a.AttendanceService.Search(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error searching attendance")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.AttendanceResponse]{
		Meta:   model.Meta{Code: 200, Status: "OK"},
		Data:   responses,
		Paging: paging,
	})
}

func (a *AttendanceApplication) Get(ctx *fiber.Ctx) error {
	request := &model.AttendanceGetByIdRequest{
		AttendanceId: ctx.Params("attendanceId"),
	}

	response, err := a.AttendanceService.FindById(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error finding attendance")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AttendanceResponse]{
		Meta: model.Meta{Code: 200, Status: "OK"},
		Data: response,
	})
}

func (a *AttendanceApplication) Update(ctx *fiber.Ctx) error {
	request := new(model.AttendanceUpdateRequest)
	if err := ctx.BodyParser(request); err != nil {
		a.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	request.AttendanceId = ctx.Params("attendanceId")
	request.UpdatedBy = ctx.Locals("employee_name").(string)

	response, err := a.AttendanceService.Update(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error updating attendance")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AttendanceResponse]{
		Meta: model.Meta{Code: 200, Status: "success update attendance"},
		Data: response,
	})
}

func (a *AttendanceApplication) SoftDelete(ctx *fiber.Ctx) error {
	request := &model.AttendanceDeleteRequest{
		AttendanceId: ctx.Params("attendanceId"),
	}

	if err := a.AttendanceService.SoftDelete(ctx.UserContext(), request); err != nil {
		a.Log.WithError(err).Error("error deleting attendance")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{
		Meta: model.Meta{Code: 200, Status: "success delete attendance"},
		Data: true,
	})
}
