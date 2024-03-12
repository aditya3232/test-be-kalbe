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
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type employeeService struct {
	DB                   *gorm.DB
	Log                  *logrus.Logger
	Validate             *validator.Validate
	EmployeeRepository   repository.EmployeeRepository
	PositionRepository   repository.PositionRepository
	DepartmentRepository repository.DepartmentRepository
}

type EmployeeService interface {
	Create(ctx context.Context, request *model.EmployeeCreateRequest) (*model.EmployeeResponse, error)
	Update(ctx context.Context, request *model.EmployeeUpdateRequest) (*model.EmployeeResponse, error)
	SoftDelete(ctx context.Context, request *model.EmployeeDeleteRequest) error
	FindById(ctx context.Context, request *model.EmployeeGetByIdRequest) (*model.EmployeeResponse, error)
	Search(ctx context.Context, request *model.EmployeeSearchRequest) ([]model.EmployeeResponse, int64, error)
}

func NewEmployeeService(db *gorm.DB, log *logrus.Logger, validate *validator.Validate, employeeRepository repository.EmployeeRepository, positionRepository repository.PositionRepository, departmentRepository repository.DepartmentRepository) *employeeService {
	return &employeeService{
		DB:                   db,
		Log:                  log,
		Validate:             validate,
		EmployeeRepository:   employeeRepository,
		PositionRepository:   positionRepository,
		DepartmentRepository: departmentRepository,
	}
}

func (s *employeeService) Create(ctx context.Context, request *model.EmployeeCreateRequest) (*model.EmployeeResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	// check department id
	department := new(entity.Department)
	if err := s.DepartmentRepository.FindById(tx, department, request.DepartmentId); err != nil {
		s.Log.WithError(err).Error("error finding department")
		return nil, fiber.ErrNotFound
	}

	// check position id
	position := new(entity.Position)
	if err := s.PositionRepository.FindById(tx, position, request.PositionId); err != nil {
		s.Log.WithError(err).Error("error finding position")
		return nil, fiber.ErrNotFound
	}

	currentTime := time.Now()
	departmentId, err := strconv.Atoi(request.DepartmentId)
	if err != nil {
		s.Log.WithError(err).Error("error parsing department ID")
		return nil, fiber.ErrInternalServerError
	}
	positionId, err := strconv.Atoi(request.PositionId)
	if err != nil {
		s.Log.WithError(err).Error("error parsing position ID")
		return nil, fiber.ErrInternalServerError
	}
	superior, err := strconv.Atoi(request.Superior)
	if err != nil {
		s.Log.WithError(err).Error("error parsing superior ID")
		return nil, fiber.ErrInternalServerError
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.MinCost)
	if err != nil {
		s.Log.WithError(err).Error("error hashing password")
		return nil, fiber.ErrInternalServerError
	}
	year := time.Now().Format("06")
	month := time.Now().Format("01")
	uuid := uuid.New().String()
	uniqueCode := year + month + "-" + uuid

	employee := &entity.Employee{
		EmployeeCode: uniqueCode,
		EmployeeName: request.EmployeeName,
		Password:     string(hashedPassword),
		DepartmentId: int64(departmentId),
		PositionId:   int64(positionId),
		Superior:     int64(superior),
		CreatedBy:    "system",
		CreatedAt:    &currentTime,
	}

	if err := s.EmployeeRepository.Create(tx, employee); err != nil {
		s.Log.WithError(err).Error("error creating employee")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error creating employee")
		return nil, fiber.ErrInternalServerError
	}

	return converter.EmployeeToResponse(employee), nil

}

func (s *employeeService) Update(ctx context.Context, request *model.EmployeeUpdateRequest) (*model.EmployeeResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	employee := new(entity.Employee)
	if err := s.EmployeeRepository.FindById(tx, employee, request.EmployeeId); err != nil {
		s.Log.WithError(err).Error("error finding employee")
		return nil, fiber.ErrNotFound
	}

	// check department id
	department := new(entity.Department)
	if err := s.DepartmentRepository.FindById(tx, department, request.DepartmentId); err != nil {
		s.Log.WithError(err).Error("error finding department")
		return nil, fiber.ErrNotFound
	}

	// check position id
	position := new(entity.Position)
	if err := s.PositionRepository.FindById(tx, position, request.PositionId); err != nil {
		s.Log.WithError(err).Error("error finding position")
		return nil, fiber.ErrNotFound
	}

	currentTime := time.Now()
	departmentId, err := strconv.Atoi(request.DepartmentId)
	if err != nil {
		s.Log.WithError(err).Error("error parsing department ID")
		return nil, fiber.ErrInternalServerError
	}
	positionId, err := strconv.Atoi(request.PositionId)
	if err != nil {
		s.Log.WithError(err).Error("error parsing position ID")
		return nil, fiber.ErrInternalServerError
	}
	superior, err := strconv.Atoi(request.Superior)
	if err != nil {
		s.Log.WithError(err).Error("error parsing superior ID")
		return nil, fiber.ErrInternalServerError
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.MinCost)
	if err != nil {
		s.Log.WithError(err).Error("error hashing password")
		return nil, fiber.ErrInternalServerError
	}

	employee.EmployeeName = request.EmployeeName
	employee.Password = string(hashedPassword)
	employee.DepartmentId = int64(departmentId)
	employee.PositionId = int64(positionId)
	employee.Superior = int64(superior)
	employee.UpdatedBy = "system"
	employee.UpdatedAt = &currentTime

	if err := s.EmployeeRepository.Update(tx, employee); err != nil {
		s.Log.WithError(err).Error("error updating employee")
		return nil, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error updating employee")
		return nil, fiber.ErrInternalServerError
	}

	return converter.EmployeeToResponse(employee), nil
}

func (s *employeeService) SoftDelete(ctx context.Context, request *model.EmployeeDeleteRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return fiber.ErrBadRequest
	}

	employee := new(entity.Employee)
	if err := s.EmployeeRepository.FindById(tx, employee, request.EmployeeId); err != nil {
		s.Log.WithError(err).Error("error deleting employee")
		return fiber.ErrNotFound
	}

	currentTime := time.Now()
	employee.DeletedAt = &currentTime

	if err := s.EmployeeRepository.Update(tx, employee); err != nil {
		s.Log.WithError(err).Error("error deleting employee")
		return fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error deleting employee")
		return fiber.ErrInternalServerError
	}

	return nil
}

func (s *employeeService) FindById(ctx context.Context, request *model.EmployeeGetByIdRequest) (*model.EmployeeResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, fiber.ErrBadRequest
	}

	employee := new(entity.Employee)
	if err := s.EmployeeRepository.FindById(tx, employee, request.EmployeeId); err != nil {
		s.Log.WithError(err).Error("error finding employee")
		return nil, fiber.ErrNotFound
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error getting employee")
		return nil, fiber.ErrInternalServerError
	}

	return converter.EmployeeToResponse(employee), nil
}

func (s *employeeService) Search(ctx context.Context, request *model.EmployeeSearchRequest) ([]model.EmployeeResponse, int64, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		s.Log.WithError(err).Error("error validating request body")
		return nil, 0, fiber.ErrBadRequest
	}

	employees, total, err := s.EmployeeRepository.Search(tx, request)
	if err != nil {
		s.Log.WithError(err).Error("error searching employee")
		return nil, 0, fiber.ErrInternalServerError
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.WithError(err).Error("error searching employee")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.EmployeeResponse, len(employees))
	for i, employee := range employees {
		responses[i] = *converter.EmployeeToResponse(&employee)
	}

	return responses, total, nil
}
