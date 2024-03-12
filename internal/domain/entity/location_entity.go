package entity

import "time"

type Location struct {
	LocationId   int64      `gorm:"column:location_id;primaryKey"`
	LocationName string     `gorm:"column:location_name"`
	CreatedAt    time.Time  `gorm:"column:created_at;default:now()"`
	CreatedBy    string     `gorm:"column:created_by"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
	UpdatedBy    string     `gorm:"column:updated_by"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
}

func (l *Location) TableName() string {
	return "location"
}
