package entity

import "time"

type Attendance struct {
	AttendanceId int64      `gorm:"column:attendance_id;primaryKey"`
	EmployeeId   int64      `gorm:"column:employee_id"`
	LocationId   int64      `gorm:"column:location_id"`
	AbsentIn     time.Time  `gorm:"column:absent_in;default:null"`
	AbsentOut    time.Time  `gorm:"column:absent_out;default:null"`
	CreatedAt    *time.Time `gorm:"column:created_at;default:now()"`
	CreatedBy    string     `gorm:"column:created_by"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
	UpdatedBy    string     `gorm:"column:updated_by"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
}

func (a *Attendance) TableName() string {
	return "attendance"
}
