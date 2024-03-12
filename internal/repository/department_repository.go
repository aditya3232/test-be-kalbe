package repository

import (
	"test-be-kalbe/internal/domain/entity"
	"test-be-kalbe/internal/domain/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type departmentRepository struct {
	DB  *gorm.DB
	Log *logrus.Logger
}

type DepartmentRepository interface {
	Create(db *gorm.DB, entity *entity.Department) error
	Update(db *gorm.DB, entity *entity.Department) error
	Delete(db *gorm.DB, entity *entity.Department) error
	CountById(db *gorm.DB, id any) (int64, error)
	FindById(db *gorm.DB, entity *entity.Department, id any) error
	Search(db *gorm.DB, request *model.DepartmentSearchRequest) ([]entity.Department, int64, error)
	Filter(request *model.DepartmentSearchRequest) func(tx *gorm.DB) *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB, log *logrus.Logger) *departmentRepository {
	return &departmentRepository{
		DB:  db,
		Log: log,
	}
}

func (r *departmentRepository) Create(db *gorm.DB, entity *entity.Department) error {
	return db.Create(entity).Error
}

func (r *departmentRepository) Update(db *gorm.DB, entity *entity.Department) error {
	return db.Save(entity).Error
}

func (r *departmentRepository) Delete(db *gorm.DB, entity *entity.Department) error {
	return db.Delete(entity).Error
}

func (r *departmentRepository) CountById(db *gorm.DB, id any) (int64, error) {
	var total int64
	err := db.Model(&entity.Department{}).Where("id = ?", id).Count(&total).Error
	return total, err
}

func (r *departmentRepository) FindById(db *gorm.DB, entity *entity.Department, id any) error {
	return db.Where("id = ?", id).Take(entity).Error
}

func (r *departmentRepository) Search(db *gorm.DB, request *model.DepartmentSearchRequest) ([]entity.Department, int64, error) {
	var departements []entity.Department
	if err := db.Scopes(r.Filter(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&departements).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.Department{}).Scopes(r.Filter(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return departements, total, nil
}

func (r *departmentRepository) Filter(request *model.DepartmentSearchRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {

		if department_name := request.DepartmentName; department_name != "" {
			department_name = "%" + department_name + "%"
			tx = tx.Where("name LIKE ?", department_name)
		}

		return tx
	}
}
