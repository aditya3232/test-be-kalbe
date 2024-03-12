package entity

import "time"

type Position struct {
	PositionId   int64      `gorm:"column:position_id;primaryKey"`
	DepartmentId int64      `gorm:"column:department_id"`
	PositionName string     `gorm:"column:position_name"`
	CreatedAt    time.Time  `gorm:"column:created_at;default:now()"`
	CreatedBy    string     `gorm:"column:created_by"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
	UpdatedBy    string     `gorm:"column:updated_by"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
}

func (p *Position) TableName() string {
	return "position"
}
