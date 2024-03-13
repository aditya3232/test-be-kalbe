package service

import (
	"context"
	"test-be-kalbe/internal/domain/converter"
	"test-be-kalbe/internal/domain/model"
	"test-be-kalbe/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type attendanceReportService struct {
	DB                         *gorm.DB
	Log                        *logrus.Logger
	Validate                   *validator.Validate
	AttendanceReportRepository repository.AttendanceReportRepository
}

type AttendanceReportService interface {
	Search(ctx context.Context, request *model.AttendanceReportSearchRequest) ([]model.AttendanceReportResponse, int64, error)
}

func NewAttendanceReportService(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, attendanceReportRepository repository.AttendanceReportRepository) *attendanceReportService {
	return &attendanceReportService{
		DB:                         db,
		Log:                        log,
		Validate:                   validate,
		AttendanceReportRepository: attendanceReportRepository,
	}
}

func (s *attendanceReportService) Search(ctx context.Context, request *model.AttendanceReportSearchRequest) ([]model.AttendanceReportResponse, int64, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	attendanceReports, total, err := s.AttendanceReportRepository.Search(tx, request)
	if err != nil {
		s.Log.WithError(err).Error("error searching attendance")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error searching attendance")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.AttendanceReportResponse, len(attendanceReports))
	for i, attendanceReport := range attendanceReports {
		responses[i] = *converter.AttendanceReportToResponse(attendanceReport)
	}

	return responses, total, nil
}
