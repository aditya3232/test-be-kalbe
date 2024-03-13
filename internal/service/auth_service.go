package service

import (
	"context"
	"test-be-kalbe/internal/domain/converter"
	"test-be-kalbe/internal/domain/entity"
	"test-be-kalbe/internal/domain/model"
	"test-be-kalbe/internal/helper"
	"test-be-kalbe/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authService struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	EmployeeRepository repository.EmployeeRepository
}

type AuthService interface {
	Login(ctx context.Context, request *model.LoginRequest) (*model.LoginResponse, error)
	Logout(ctx context.Context, request *model.LogoutRequest) (*model.LogoutResponse, error)
}

func NewAuthService(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, employeeRepository repository.EmployeeRepository) *authService {
	return &authService{
		DB:                 db,
		Log:                log,
		Validate:           validate,
		EmployeeRepository: employeeRepository,
	}
}

func (s *authService) Login(ctx context.Context, request *model.LoginRequest) (*model.LoginResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	employee := new(entity.Employee)
	if err := s.EmployeeRepository.FindByEmployeeName(tx, employee, request.EmployeeName); err != nil {
		s.Log.WithError(err).Error("error finding employee name")
		return nil, fiber.ErrNotFound
	}

	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(request.Password)); err != nil {
		s.Log.WithError(err).Error("error comparing password")
		return nil, fiber.ErrUnauthorized
	}

	// generate token
	token, expirationTime, err := helper.GenerateToken(int(employee.EmployeeId))
	if err != nil {
		s.Log.WithError(err).Error("error generating token")
		return nil, fiber.ErrInternalServerError
	}

	// commit transaction
	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error commit transaction")
		return nil, fiber.ErrInternalServerError
	}

	return converter.AuthToResponse(token, expirationTime), nil

}

func (s *authService) Logout(ctx context.Context, request *model.LogoutRequest) (*model.LogoutResponse, error) {
	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	// validate token
	_, err := helper.ValidateToken(request.Token)
	if err != nil {
		s.Log.WithError(err).Error("error validating token")
		return nil, fiber.ErrUnauthorized
	}

	// validate employee id from claims
	employeeId, err := helper.GetEmployeeIDFromToken(request.Token)
	if err != nil {
		s.Log.WithError(err).Error("error getting employee id from token")
		return nil, fiber.ErrUnauthorized
	}

	// invalidate token
	expired, err := helper.InvalidateToken(request.Token)
	if err != nil {
		s.Log.WithError(err).Error("error invalidating token")
		return nil, fiber.ErrInternalServerError
	}

	return converter.LogoutToResponse(employeeId, expired), nil
}
