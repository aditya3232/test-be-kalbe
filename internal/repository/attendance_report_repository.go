package repository

import (
	"test-be-kalbe/internal/domain/model"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type attendanceReportRepository struct {
	DB  *gorm.DB
	Log *logrus.Logger
}

type AttendanceReportRepository interface {
	Search(db *gorm.DB, request *model.AttendanceReportSearchRequest) ([]model.AttendanceReportResponse, int64, error)
}

func NewAttendanceReportRepository(db *gorm.DB, log *logrus.Logger) *attendanceReportRepository {
	return &attendanceReportRepository{
		DB:  db,
		Log: log,
	}
}

func (r *attendanceReportRepository) Search(db *gorm.DB, request *model.AttendanceReportSearchRequest) ([]model.AttendanceReportResponse, int64, error) {
	var attendances []model.AttendanceReportResponse

	db = db.Table("attendance").
		Select("attendance.created_at as date, employee.employee_code as employee_code, employee.employee_name as employee_name, department.department_name as department_name, position.position_name as position_name, location.location_name as location_name, attendance.absent_in, attendance.absent_out").
		Joins("JOIN employee ON employee.employee_id = attendance.employee_id").
		Joins("JOIN department ON department.department_id = employee.department_id").
		Joins("JOIN location ON location.location_id = attendance.location_id").
		Joins("JOIN position ON position.position_id = employee.position_id")

	if err := db.Scopes(r.Filter(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&attendances).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := db.Scopes(r.Filter(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return attendances, total, nil
}

func (r *attendanceReportRepository) Filter(request *model.AttendanceReportSearchRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if request.TimeInterval != "" {
			// Initialize variables for start and end time
			var startTime, endTime time.Time

			// Get current time
			now := time.Now()

			// Calculate start and end time based on TimeInterval
			switch request.TimeInterval {
			case "day":
				startTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
				endTime = startTime.AddDate(0, 0, 1).Add(-time.Second)
			case "week":
				weekday := int(now.Weekday())
				startTime = now.AddDate(0, 0, -weekday)
				startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, startTime.Location())
				endTime = startTime.AddDate(0, 0, 7).Add(-time.Second)
			case "month":
				startTime = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
				endTime = startTime.AddDate(0, 1, 0).Add(-time.Second)
			case "year":
				startTime = time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())
				endTime = startTime.AddDate(1, 0, 0).Add(-time.Second)
			}

			// Apply time interval condition to the query
			tx = tx.Where("attendance.created_at >= ? AND attendance.created_at <= ?", startTime, endTime)
		} else {
			// If no TimeInterval specified, fetch all data
			tx = tx.Where("attendance.created_at IS NOT NULL")
		}

		tx = tx.Where("attendance.deleted_at IS NULL")

		return tx
	}
}
