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

type locationService struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	LocationRepository repository.LocationRepository
}

type LocationService interface {
	Create(ctx context.Context, request *model.LocationCreateRequest) (*model.LocationResponse, error)
	Update(ctx context.Context, request *model.LocationUpdateRequest) (*model.LocationResponse, error)
	SoftDelete(ctx context.Context, request *model.LocationDeleteRequest) error
	FindById(ctx context.Context, request *model.LocationGetByIdRequest) (*model.LocationResponse, error)
	Search(ctx context.Context, request *model.LocationSearchRequest) ([]model.LocationResponse, int64, error)
}

func NewLocationService(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, locationRepository repository.LocationRepository) *locationService {
	return &locationService{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		LocationRepository: locationRepository,
	}
}

func (s *locationService) Create(ctx context.Context, request *model.LocationCreateRequest) (*model.LocationResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	currentTime := time.Now()
	location := &entity.Location{
		LocationName: request.LocationName,
		CreatedBy:    request.CreatedBy,
		CreatedAt:    &currentTime,
	}

	if err := s.LocationRepository.Create(tx, location); err != nil {
		s.Log.WithError(err).Error("error creating location")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error creating department")
		return nil, fiber.ErrInternalServerError
	}

	return converter.LocationToResponse(location), nil
}

func (s *locationService) Update(ctx context.Context, request *model.LocationUpdateRequest) (*model.LocationResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	location := new(entity.Location)
	if err := s.LocationRepository.FindById(tx, location, request.LocationId); err != nil {
		s.Log.WithError(err).Error("error finding location")
		return nil, fiber.ErrNotFound
	}

	currentTime := time.Now()
	location.LocationName = request.LocationName
	location.UpdatedBy = request.UpdatedBy
	location.UpdatedAt = &currentTime

	if err := s.LocationRepository.Update(tx, location); err != nil {
		s.Log.WithError(err).Error("error updating location")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error updating location")
		return nil, fiber.ErrInternalServerError
	}

	return converter.LocationToResponse(location), nil
}

func (s *locationService) SoftDelete(ctx context.Context, request *model.LocationDeleteRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return fiber.ErrBadRequest
	}

	location := new(entity.Location)
	if err := s.LocationRepository.FindById(tx, location, request.LocationId); err != nil {
		s.Log.WithError(err).Error("error deleting location")
		return fiber.ErrNotFound
	}

	currentTime := time.Now()
	location.DeletedAt = &currentTime

	if err := s.LocationRepository.Update(tx, location); err != nil {
		s.Log.WithError(err).Error("error deleting location")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error deleting location")
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *locationService) FindById(ctx context.Context, request *model.LocationGetByIdRequest) (*model.LocationResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	location := new(entity.Location)
	if err := s.LocationRepository.FindById(tx, location, request.LocationId); err != nil {
		s.Log.WithError(err).Error("error finding location")
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error getting department")
		return nil, fiber.ErrInternalServerError
	}

	return converter.LocationToResponse(location), nil
}

func (s *locationService) Search(ctx context.Context, request *model.LocationSearchRequest) ([]model.LocationResponse, int64, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	locations, total, err := s.LocationRepository.Search(tx, request)
	if err != nil {
		s.Log.WithError(err).Error("error searching location")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error searching location")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.LocationResponse, len(locations))
	for i, location := range locations {
		responses[i] = *converter.LocationToResponse(&location)
	}

	return responses, total, nil
}
