package converter

import (
	"test-be-kalbe/internal/domain/entity"
	"test-be-kalbe/internal/domain/model"
)

func AttendanceToResponse(attendance *entity.Attendance) *model.AttendanceResponse {
	return &model.AttendanceResponse{
		AttendanceId: attendance.AttendanceId,
		EmployeeId:   attendance.EmployeeId,
		LocationId:   attendance.LocationId,
		AbsentIn:     attendance.AbsentIn.Format("2006-01-02 15:04:05"),
		AbsentOut:    attendance.AbsentOut.Format("2006-01-02 15:04:05"),
		CreatedAt:    attendance.CreatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:    attendance.CreatedBy,
		UpdatedAt:    attendance.UpdatedAt.Format("2006-01-02 15:04:05"),
		UpdatedBy:    attendance.UpdatedBy,
	}
}
