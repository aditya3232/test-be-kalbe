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

	// join data from attendance, employee, department, location and position
	db = db.Table("attendances").
		Select("attendances.created_at, employees.code as employee_code, employees.name as employee_name, departments.name as department_name, positions.name as position_name, locations.name as location_name, attendances.absent_in, attendances.absent_out").
		Joins("JOIN employees ON employees.id = attendances.employee_id").
		Joins("JOIN departments ON departments.id = employees.department_id").
		Joins("JOIN locations ON locations.id = attendances.location_id").
		Joins("JOIN positions ON positions.id = employees.position_id")

	if err := db.Scopes(r.Filter(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&attendances).Error; err != nil {
		return nil, 0, err
	}
	return attendances, 0, nil
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
			tx = tx.Where("attendances.created_at >= ? AND attendances.created_at <= ?", startTime, endTime)
		} else {
			// If no TimeInterval specified, fetch all data
			tx = tx.Where("attendances.created_at IS NOT NULL")
		}

		tx = tx.Where("attendances.deleted_at IS NULL")

		return tx
	}
}
