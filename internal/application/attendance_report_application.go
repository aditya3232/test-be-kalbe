package application

import (
	"math"
	"test-be-kalbe/internal/service"

	"test-be-kalbe/internal/domain/model"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AttendanceReportApplication struct {
	AttendanceReportService service.AttendanceReportService
	Log                     *logrus.Logger
}

func NewAttendanceReportApplication(attendanceReportService service.AttendanceReportService, log *logrus.Logger) *AttendanceReportApplication {
	return &AttendanceReportApplication{
		AttendanceReportService: attendanceReportService,
		Log:                     log,
	}
}

func (a *AttendanceReportApplication) List(ctx *fiber.Ctx) error {
	request := &model.AttendanceReportSearchRequest{
		TimeInterval: ctx.Query("time_interval", ""),
		Page:         ctx.QueryInt("page", 1),
		Size:         ctx.QueryInt("size", 10),
	}

	responses, total, err := a.AttendanceReportService.Search(ctx.UserContext(), request)
	if err != nil {
		a.Log.WithError(err).Error("error searching attendance report")
		return err
	}

	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(request.Size))),
	}

	return ctx.JSON(model.WebResponse[[]model.AttendanceReportResponse]{
		Meta:   model.Meta{Code: 200, Status: "OK"},
		Data:   responses,
		Paging: paging,
	})
}
