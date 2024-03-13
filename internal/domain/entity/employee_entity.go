package entity

import "time"

type Employee struct {
	EmployeeId   int64      `gorm:"column:employee_id;primaryKey"`
	EmployeeCode string     `gorm:"column:employee_code"`
	EmployeeName string     `gorm:"column:employee_name"`
	Password     string     `gorm:"column:password"`
	DepartmentId int64      `gorm:"column:department_id"`
	PositionId   int64      `gorm:"column:position_id"`
	Superior     int64      `gorm:"column:superior"`
	CreatedAt    *time.Time `gorm:"column:created_at;default:now()"`
	CreatedBy    string     `gorm:"column:created_by"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
	UpdatedBy    string     `gorm:"column:updated_by"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
	Token        string     `gorm:"column:token"`
}

func (e *Employee) TableName() string {
	return "employee"
}
