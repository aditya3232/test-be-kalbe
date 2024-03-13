package repository

import (
	"test-be-kalbe/internal/domain/entity"
	"test-be-kalbe/internal/domain/model"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type attendanceRepository struct {
	DB  *gorm.DB
	Log *logrus.Logger
}

type AttendanceRepository interface {
	Create(db *gorm.DB, entity *entity.Attendance) error
	Update(db *gorm.DB, entity *entity.Attendance) error
	Delete(db *gorm.DB, entity *entity.Attendance) error
	CountById(db *gorm.DB, id any) (int64, error)
	FindById(db *gorm.DB, entity *entity.Attendance, id any) error
	Search(db *gorm.DB, request *model.AttendanceSearchRequest) ([]entity.Attendance, int64, error)
	Filter(request *model.AttendanceSearchRequest) func(tx *gorm.DB) *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB, log *logrus.Logger) *attendanceRepository {
	return &attendanceRepository{
		DB:  db,
		Log: log,
	}
}

func (r *attendanceRepository) Create(db *gorm.DB, entity *entity.Attendance) error {
	return db.Create(entity).Error
}

func (r *attendanceRepository) Update(db *gorm.DB, entity *entity.Attendance) error {
	return db.Save(entity).Error
}

func (r *attendanceRepository) Delete(db *gorm.DB, entity *entity.Attendance) error {
	return db.Delete(entity).Error
}

func (r *attendanceRepository) CountById(db *gorm.DB, id any) (int64, error) {
	var total int64
	err := db.Model(&entity.Attendance{}).Where("attendance_id = ? AND deleted_at IS NULL", id).Count(&total).Error
	return total, err
}

func (r *attendanceRepository) FindById(db *gorm.DB, entity *entity.Attendance, id any) error {
	return db.Where("attendance_id = ? AND deleted_at IS NULL", id).Take(entity).Error
}

func (r *attendanceRepository) Search(db *gorm.DB, request *model.AttendanceSearchRequest) ([]entity.Attendance, int64, error) {
	var attendances []entity.Attendance
	if err := db.Scopes(r.Filter(request)).Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&attendances).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.Attendance{}).Scopes(r.Filter(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return attendances, total, nil
}

func (r *attendanceRepository) Filter(request *model.AttendanceSearchRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if request.EmployeeId != "" {
			tx = tx.Where("employee_id = ?", request.EmployeeId)
		}
		if request.LocationId != "" {
			tx = tx.Where("location_id = ?", request.LocationId)
		}
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
			tx = tx.Where("created_at >= ? AND created_at <= ?", startTime, endTime)
		} else {
			// If no TimeInterval specified, fetch all data
			tx = tx.Where("created_at IS NOT NULL")
		}

		tx = tx.Where("deleted_at IS NULL")

		return tx
	}
}
