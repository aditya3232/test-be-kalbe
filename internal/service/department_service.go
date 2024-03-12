package service

import (
	"context"
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

type departmentService struct {
	DB                   *gorm.DB
	Log                  *logrus.Logger
	Validate             *validator.Validate
	DepartmentRepository repository.DepartmentRepository
}

type DepartmentService interface {
	Create(ctx context.Context, request *model.DepartmentCreateRequest) (*model.DepartmentResponse, error)
	Update(ctx context.Context, request *model.DepartmentUpdateRequest) (*model.DepartmentResponse, error)
	SoftDelete(ctx context.Context, request *model.DepartmentDeleteRequest) error
	FindById(ctx context.Context, request *model.DepartmentGetByIdRequest) (*model.DepartmentResponse, error)
	Search(ctx context.Context, request *model.DepartmentSearchRequest) ([]model.DepartmentResponse, int64, error)
}

func NewDepartmentService(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, departmentRepository repository.DepartmentRepository) *departmentService {
	return &departmentService{
		DB:                   db,
		Log:                  log,
		Validate:             validate,
		DepartmentRepository: departmentRepository,
	}
}

func (s *departmentService) Create(ctx context.Context, request *model.DepartmentCreateRequest) (*model.DepartmentResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	currentTime := time.Now()
	department := &entity.Department{
		DepartmentName: request.DepartmentName,
		CreatedBy:      "system",
		CreatedAt:      &currentTime,
	}

	if err := s.DepartmentRepository.Create(tx, department); err != nil {
		s.Log.WithError(err).Error("error creating department")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error creating department")
		return nil, fiber.ErrInternalServerError
	}

	return converter.DepartmentToResponse(department), nil

}

func (s *departmentService) Update(ctx context.Context, request *model.DepartmentUpdateRequest) (*model.DepartmentResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	department := new(entity.Department)
	if err := s.DepartmentRepository.FindById(tx, department, request.DepartmentId); err != nil {
		s.Log.WithError(err).Error("error updating department")
		return nil, fiber.ErrNotFound
	}

	currentTime := time.Now()
	department.DepartmentName = request.DepartmentName
	department.UpdatedBy = "system"
	department.UpdatedAt = &currentTime

	if err := s.DepartmentRepository.Update(tx, department); err != nil {
		s.Log.WithError(err).Error("error updating department")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error updating department")
		return nil, fiber.ErrInternalServerError
	}

	return converter.DepartmentToResponse(department), nil
}

func (s *departmentService) SoftDelete(ctx context.Context, request *model.DepartmentDeleteRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return fiber.ErrBadRequest
	}

	department := new(entity.Department)
	if err := s.DepartmentRepository.FindById(tx, department, request.DepartmentId); err != nil {
		s.Log.WithError(err).Error("error deleting department")
		return fiber.ErrNotFound
	}

	currentTime := time.Now()
	department.DeletedAt = &currentTime

	if err := s.DepartmentRepository.Update(tx, department); err != nil {
		s.Log.WithError(err).Error("error deleting department")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error deleting department")
		return fiber.ErrInternalServerError
	}

	return nil

}

func (s *departmentService) FindById(ctx context.Context, request *model.DepartmentGetByIdRequest) (*model.DepartmentResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	department := new(entity.Department)
	if err := s.DepartmentRepository.FindById(tx, department, request.DepartmentId); err != nil {
		s.Log.WithError(err).Error("error getting department")
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error getting department")
		return nil, fiber.ErrInternalServerError
	}

	return converter.DepartmentToResponse(department), nil
}

func (s *departmentService) Search(ctx context.Context, request *model.DepartmentSearchRequest) ([]model.DepartmentResponse, int64, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	departments, total, err := s.DepartmentRepository.Search(tx, request)
	if err != nil {
		s.Log.WithError(err).Error("error searching department")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error searching department")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.DepartmentResponse, len(departments))
	for i, department := range departments {
		responses[i] = *converter.DepartmentToResponse(&department)
	}

	return responses, total, nil
}
