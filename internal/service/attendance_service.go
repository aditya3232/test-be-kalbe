package service

import (
	"context"
	"strconv"
	"test-be-kalbe/internal/domain/converter"
	"test-be-kalbe/internal/domain/entity"
	"test-be-kalbe/internal/domain/model"
	"test-be-kalbe/internal/repository"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type attendanceService struct {
	DB                   *gorm.DB
	Log                  *logrus.Logger
	Validate             *validator.Validate
	AttendanceRepository repository.AttendanceRepository
	EmployeeRepository   repository.EmployeeRepository
	LocationRepository   repository.LocationRepository
}

type AttendanceService interface {
	Create(ctx context.Context, request *model.AttendanceCreateRequest) (*model.AttendanceResponse, error)
	Update(ctx context.Context, request *model.AttendanceUpdateRequest) (*model.AttendanceResponse, error)
	SoftDelete(ctx context.Context, request *model.AttendanceDeleteRequest) error
	FindById(ctx context.Context, request *model.AttendanceGetByIdRequest) (*model.AttendanceResponse, error)
	Search(ctx context.Context, request *model.AttendanceSearchRequest) ([]model.AttendanceResponse, int64, error)
}

func NewAttendanceService(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, attendanceRepository repository.AttendanceRepository, employeeRepository repository.EmployeeRepository, locationRepository repository.LocationRepository) *attendanceService {
	return &attendanceService{
		DB:                   db,
		Log:                  log,
		Validate:             validate,
		AttendanceRepository: attendanceRepository,
		EmployeeRepository:   employeeRepository,
		LocationRepository:   locationRepository,
	}
}

func (s *attendanceService) Create(ctx context.Context, request *model.AttendanceCreateRequest) (*model.AttendanceResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	// check employee id
	employee := new(entity.Employee)
	if err := s.EmployeeRepository.FindById(tx, employee, request.EmployeeId); err != nil {
		s.Log.WithError(err).Error("error finding employee")
		return nil, fiber.ErrNotFound
	}

	// check location id
	location := new(entity.Location)
	if err := s.LocationRepository.FindById(tx, location, request.LocationId); err != nil {
		s.Log.WithError(err).Error("error finding location")
		return nil, fiber.ErrNotFound
	}

	currentTime := time.Now()

	employeeId, err := strconv.Atoi(request.EmployeeId)
	if err != nil {
		s.Log.WithError(err).Error("error parsing employee id")
		return nil, fiber.ErrInternalServerError
	}

	locationId, err := strconv.Atoi(request.LocationId)
	if err != nil {
		s.Log.WithError(err).Error("error parsing location id")
		return nil, fiber.ErrInternalServerError
	}

	absentIn := time.Time{}
	if request.AbsentIn != "" {
		var err error
		absentIn, err = time.Parse("2006-01-02 15:04:05", request.AbsentIn)
		if err != nil {
			s.Log.WithError(err).Error("error parsing AbsentIn")
			return nil, fiber.ErrInternalServerError
		}
	}

	absentOut := time.Time{}
	if request.AbsentOut != "" {
		absentOut, err = time.Parse("2006-01-02 15:04:05", request.AbsentOut)
		if err != nil {
			s.Log.WithError(err).Error("error parsing AbsentOut")
			return nil, fiber.ErrInternalServerError
		}
	}

	attendance := &entity.Attendance{
		EmployeeId: int64(employeeId),
		LocationId: int64(locationId),
		CreatedBy:  request.CreatedBy,
		CreatedAt:  &currentTime,
	}

	if request.AbsentIn != "" {
		attendance.AbsentIn = absentIn
	} else {
		attendance.AbsentIn = time.Time{}
	}

	if request.AbsentOut != "" {
		attendance.AbsentOut = absentOut
	} else {
		attendance.AbsentOut = time.Time{}
	}

	if err := s.AttendanceRepository.Create(tx, attendance); err != nil {
		s.Log.WithError(err).Error("error creating attendance")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error creating employee")
		return nil, fiber.ErrInternalServerError
	}

	return converter.AttendanceToResponse(attendance), nil
}

func (s *attendanceService) Update(ctx context.Context, request *model.AttendanceUpdateRequest) (*model.AttendanceResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	attendance := new(entity.Attendance)
	if err := s.AttendanceRepository.FindById(tx, attendance, request.AttendanceId); err != nil {
		s.Log.WithError(err).Error("error finding attendance")
		return nil, fiber.ErrNotFound
	}

	// check employee id
	employee := new(entity.Employee)
	if err := s.EmployeeRepository.FindById(tx, employee, request.EmployeeId); err != nil {
		s.Log.WithError(err).Error("error finding employee")
		return nil, fiber.ErrNotFound
	}

	// check location id
	location := new(entity.Location)
	if err := s.LocationRepository.FindById(tx, location, request.LocationId); err != nil {
		s.Log.WithError(err).Error("error finding location")
		return nil, fiber.ErrNotFound
	}

	currentTime := time.Now()
	employeeId, err := strconv.Atoi(request.EmployeeId)
	if err != nil {
		s.Log.WithError(err).Error("error parsing employee id")
		return nil, fiber.ErrInternalServerError
	}
	locationId, err := strconv.Atoi(request.LocationId)
	if err != nil {
		s.Log.WithError(err).Error("error parsing location id")
		return nil, fiber.ErrInternalServerError
	}
	absentIn, err := time.Parse("2006-01-02 15:04:05", request.AbsentIn)
	if err != nil {
		s.Log.WithError(err).Error("error parsing AbsentIn")
		return nil, fiber.ErrInternalServerError
	}
	absentOut, err := time.Parse("2006-01-02 15:04:05", request.AbsentOut)
	if err != nil {
		s.Log.WithError(err).Error("error parsing AbsentOut")
		return nil, fiber.ErrInternalServerError
	}

	attendance.EmployeeId = int64(employeeId)
	attendance.LocationId = int64(locationId)
	attendance.AbsentIn = absentIn
	attendance.AbsentOut = absentOut
	attendance.UpdatedBy = request.UpdatedBy
	attendance.UpdatedAt = &currentTime

	if err := s.AttendanceRepository.Update(tx, attendance); err != nil {
		s.Log.WithError(err).Error("error updating attendance")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error updating attendance")
		return nil, fiber.ErrInternalServerError
	}

	return converter.AttendanceToResponse(attendance), nil
}

func (s *attendanceService) SoftDelete(ctx context.Context, request *model.AttendanceDeleteRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return fiber.ErrBadRequest
	}

	attendance := new(entity.Attendance)
	if err := s.AttendanceRepository.FindById(tx, attendance, request.AttendanceId); err != nil {
		s.Log.WithError(err).Error("error deleting attendance")
		return fiber.ErrNotFound
	}

	currentTime := time.Now()
	attendance.DeletedAt = &currentTime

	if err := s.AttendanceRepository.Update(tx, attendance); err != nil {
		s.Log.WithError(err).Error("error deleting attendance")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error deleting attendance")
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *attendanceService) FindById(ctx context.Context, request *model.AttendanceGetByIdRequest) (*model.AttendanceResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	attendance := new(entity.Attendance)
	if err := s.AttendanceRepository.FindById(tx, attendance, request.AttendanceId); err != nil {
		s.Log.WithError(err).Error("error finding attendance")
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error getting employee")
		return nil, fiber.ErrInternalServerError
	}

	return converter.AttendanceToResponse(attendance), nil
}

func (s *attendanceService) Search(ctx context.Context, request *model.AttendanceSearchRequest) ([]model.AttendanceResponse, int64, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	attendances, total, err := s.AttendanceRepository.Search(tx, request)
	if err != nil {
		s.Log.WithError(err).Error("error searching attendance")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error searching attendance")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.AttendanceResponse, len(attendances))
	for i, attendance := range attendances {
		responses[i] = *converter.AttendanceToResponse(&attendance)
	}

	return responses, total, nil
}
