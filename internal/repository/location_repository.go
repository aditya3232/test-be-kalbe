package repository

import (
	"test-be-kalbe/internal/domain/entity"
	"test-be-kalbe/internal/domain/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type locationRepository struct {
	DB  *gorm.DB
	Log *logrus.Logger
}

type LocationRepository interface {
	Create(db *gorm.DB, entity *entity.Location) error
	Update(db *gorm.DB, entity *entity.Location) error
	Delete(db *gorm.DB, entity *entity.Location) error
	CountById(db *gorm.DB, id any) (int64, error)
	FindById(db *gorm.DB, entity *entity.Location, id any) error
	Search(db *gorm.DB, request *model.LocationSearchRequest) ([]entity.Location, int64, error)
	Filter(request *model.LocationSearchRequest) func(tx *gorm.DB) *gorm.DB
}

func NewLocationRepository(db *gorm.DB, log *logrus.Logger) *locationRepository {
	return &locationRepository{
		DB:  db,
		Log: log,
	}
}

func (r *locationRepository) Create(db *gorm.DB, entity *entity.Location) error {
	return db.Create(entity).Error
}

func (r *locationRepository) Update(db *gorm.DB, entity *entity.Location) error {
	return db.Save(entity).Error
}

func (r *locationRepository) Delete(db *gorm.DB, entity *entity.Location) error {
	return db.Delete(entity).Error
}

func (r *locationRepository) CountById(db *gorm.DB, id any) (int64, error) {
	var total int64
	err := db.Model(&entity.Location{}).Where("location_id = ? AND deleted_at IS NULL", id).Count(&total).Error
	return total, err
}

func (r *locationRepository) FindById(db *gorm.DB, entity *entity.Location, id any) error {
	return db.Where("location_id = ? AND deleted_at IS NULL", id).Take(entity).Error
}

func (r *locationRepository) Search(db *gorm.DB, request *model.LocationSearchRequest) ([]entity.Location, int64, error) {
	var locations []entity.Location
	if err := db.Scopes(r.Filter(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&locations).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.Location{}).Scopes(r.Filter(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return locations, total, nil
}

func (r *locationRepository) Filter(request *model.LocationSearchRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if request.LocationName != "" {
			tx = tx.Where("location_name LIKE ?", "%"+request.LocationName+"%")
		}

		tx = tx.Where("deleted_at IS NULL")

		return tx
	}
}
