package repository

import (
	"test-be-kalbe/internal/domain/entity"
	"test-be-kalbe/internal/domain/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type positionRepository struct {
	DB  *gorm.DB
	Log *logrus.Logger
}

type PositionRepository interface {
	Create(db *gorm.DB, entity *entity.Position) error
	Update(db *gorm.DB, entity *entity.Position) error
	Delete(db *gorm.DB, entity *entity.Position) error
	CountById(db *gorm.DB, id any) (int64, error)
	FindById(db *gorm.DB, entity *entity.Position, id any) error
	Search(db *gorm.DB, request *model.PositionSearchRequest) ([]entity.Position, int64, error)
	Filter(request *model.PositionSearchRequest) func(tx *gorm.DB) *gorm.DB
}

func NewPositionRepository(db *gorm.DB, log *logrus.Logger) *positionRepository {
	return &positionRepository{
		DB:  db,
		Log: log,
	}
}

func (r *positionRepository) Create(db *gorm.DB, entity *entity.Position) error {
	return db.Create(entity).Error
}

func (r *positionRepository) Update(db *gorm.DB, entity *entity.Position) error {
	return db.Save(entity).Error
}

func (r *positionRepository) Delete(db *gorm.DB, entity *entity.Position) error {
	return db.Delete(entity).Error
}

func (r *positionRepository) CountById(db *gorm.DB, id any) (int64, error) {
	var total int64
	err := db.Model(&entity.Position{}).Where("position_id = ?", id).Count(&total).Error
	return total, err
}

func (r *positionRepository) FindById(db *gorm.DB, entity *entity.Position, id any) error {
	return db.Where("position_id = ? AND deleted_at IS NULL", id).Take(entity).Error
}

func (r *positionRepository) Search(db *gorm.DB, request *model.PositionSearchRequest) ([]entity.Position, int64, error) {
	var positions []entity.Position
	if err := db.Scopes(r.Filter(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&positions).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.Position{}).Scopes(r.Filter(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return positions, total, nil
}

func (r *positionRepository) Filter(request *model.PositionSearchRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if request.PositionName != "" {
			tx = tx.Where("position_name LIKE ?", "%"+request.PositionName+"%")
		}

		tx = tx.Where("deleted_at IS NULL")

		return tx
	}
}
