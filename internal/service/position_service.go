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

type positionService struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	PositionRepository repository.PositionRepository
}

type PositionService interface {
	Create(ctx context.Context, request *model.PositionCreateRequest) (*model.PositionResponse, error)
	Update(ctx context.Context, request *model.PositionUpdateRequest) (*model.PositionResponse, error)
	SoftDelete(ctx context.Context, request *model.PositionDeleteRequest) error
	FindById(ctx context.Context, request *model.PositionGetByIdRequest) (*model.PositionResponse, error)
	Search(ctx context.Context, request *model.PositionSearchRequest) ([]model.PositionResponse, int64, error)
}

func NewPositionService(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, positionRepository repository.PositionRepository) *positionService {
	return &positionService{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		PositionRepository: positionRepository,
	}
}

func (s *positionService) Create(ctx context.Context, request *model.PositionCreateRequest) (*model.PositionResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	currentTime := time.Now()
	position := &entity.Position{
		DepartmentId: request.DepartmentId,
		PositionName: request.PositionName,
		CreatedBy:    "system",
		CreatedAt:    &currentTime,
	}

	if err := s.PositionRepository.Create(tx, position); err != nil {
		s.Log.WithError(err).Error("error creating position")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error creating position")
		return nil, fiber.ErrInternalServerError
	}

	return converter.PositionToResponse(position), nil
}

func (s *positionService) Update(ctx context.Context, request *model.PositionUpdateRequest) (*model.PositionResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	position := new(entity.Position)
	if err := s.PositionRepository.FindById(tx, position, request.PositionId); err != nil {
		s.Log.WithError(err).Error("error finding position")
		return nil, fiber.ErrNotFound
	}

	currentTime := time.Now()
	position.PositionName = request.PositionName
	position.UpdatedBy = "system"
	position.UpdatedAt = &currentTime

	if err := s.PositionRepository.Update(tx, position); err != nil {
		s.Log.WithError(err).Error("error updating position")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error updating position")
		return nil, fiber.ErrInternalServerError
	}

	return converter.PositionToResponse(position), nil
}

func (s *positionService) SoftDelete(ctx context.Context, request *model.PositionDeleteRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return fiber.ErrBadRequest
	}

	position := new(entity.Position)
	if err := s.PositionRepository.FindById(tx, position, request.PositionId); err != nil {
		s.Log.WithError(err).Error("error deleting position")
		return fiber.ErrNotFound
	}

	currentTime := time.Now()
	position.DeletedAt = &currentTime

	if err := s.PositionRepository.Update(tx, position); err != nil {
		s.Log.WithError(err).Error("error deleting position")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error deleting position")
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *positionService) FindById(ctx context.Context, request *model.PositionGetByIdRequest) (*model.PositionResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	position := new(entity.Position)
	if err := s.PositionRepository.FindById(tx, position, request.PositionId); err != nil {
		s.Log.WithError(err).Error("error getting position")
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error getting position")
		return nil, fiber.ErrInternalServerError
	}

	return converter.PositionToResponse(position), nil
}

func (s *positionService) Search(ctx context.Context, request *model.PositionSearchRequest) ([]model.PositionResponse, int64, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	positions, total, err := s.PositionRepository.Search(tx, request)
	if err != nil {
		s.Log.WithError(err).Error("error searching position")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.PositionResponse, len(positions))
	for i, position := range positions {
		responses[i] = *converter.PositionToResponse(&position)
	}

	return responses, total, nil
}
