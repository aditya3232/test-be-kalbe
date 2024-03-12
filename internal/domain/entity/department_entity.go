package entity

import "time"

type Department struct {
	DepartmentId   int64      `gorm:"column:department_id;primaryKey"`
	DepartmentName string     `gorm:"column:department_name"`
	CreatedAt      *time.Time `gorm:"column:created_at;default:now()"`
	CreatedBy      string     `gorm:"column:created_by"`
	UpdatedAt      *time.Time `gorm:"column:updated_at"`
	UpdatedBy      string     `gorm:"column:updated_by"`
	DeletedAt      *time.Time `gorm:"column:deleted_at"`
}

func (d *Department) TableName() string {
	return "department"
}
