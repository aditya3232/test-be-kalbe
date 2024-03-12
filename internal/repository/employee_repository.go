package repository

import (
	"test-be-kalbe/internal/domain/entity"
	"test-be-kalbe/internal/domain/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type employeeRepository struct {
	DB  *gorm.DB
	Log *logrus.Logger
}

type EmployeeRepository interface {
	Create(db *gorm.DB, entity *entity.Employee) error
	Update(db *gorm.DB, entity *entity.Employee) error
	Delete(db *gorm.DB, entity *entity.Employee) error
	CountById(db *gorm.DB, id any) (int64, error)
	FindById(db *gorm.DB, entity *entity.Employee, id any) error
	Search(db *gorm.DB, request *model.EmployeeSearchRequest) ([]entity.Employee, int64, error)
	Filter(request *model.EmployeeSearchRequest) func(tx *gorm.DB) *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB, log *logrus.Logger) *employeeRepository {
	return &employeeRepository{
		DB:  db,
		Log: log,
	}
}

func (r *employeeRepository) Create(db *gorm.DB, entity *entity.Employee) error {
	return db.Create(entity).Error
}

func (r *employeeRepository) Update(db *gorm.DB, entity *entity.Employee) error {
	return db.Save(entity).Error
}

func (r *employeeRepository) Delete(db *gorm.DB, entity *entity.Employee) error {
	return db.Delete(entity).Error
}

func (r *employeeRepository) CountById(db *gorm.DB, id any) (int64, error) {
	var total int64
	err := db.Model(&entity.Employee{}).Where("employee_id = ?", id).Count(&total).Error
	return total, err
}

func (r *employeeRepository) FindById(db *gorm.DB, entity *entity.Employee, id any) error {
	return db.Where("employee_id = ? AND deleted_at IS NULL", id).Take(entity).Error
}

func (r *employeeRepository) Search(db *gorm.DB, request *model.EmployeeSearchRequest) ([]entity.Employee, int64, error) {
	var employees []entity.Employee
	if err := db.Scopes(r.Filter(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&employees).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := db.Model(&entity.Employee{}).Scopes(r.Filter(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return employees, total, nil
}

func (r *employeeRepository) Filter(request *model.EmployeeSearchRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if request.EmployeeName != "" {
			tx = tx.Where("employee_name LIKE ?", "%"+request.EmployeeName+"%")
		}

		tx = tx.Where("deleted_at IS NULL")

		return tx
	}
}
