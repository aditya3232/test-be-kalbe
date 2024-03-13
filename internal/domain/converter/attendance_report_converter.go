package converter

import (
	"test-be-kalbe/internal/domain/entity"
	"test-be-kalbe/internal/domain/model"
)

func AttendanceReportToResponse(department *entity.Department, position *entity.Position, location *entity.Location, employee *entity.Employee, attendance *entity.Attendance) *model.AttendanceReportResponse {
	return &model.AttendanceReportResponse{
		Date:           attendance.CreatedAt.Format("2006-01-02 15:04:05"),
		EmployeeCode:   employee.EmployeeCode,
		EmployeeName:   employee.EmployeeName,
		DepartmentName: department.DepartmentName,
		PositionName:   position.PositionName,
		LocationName:   location.LocationName,
		AbsentIn:       attendance.AbsentIn.Format("2006-01-02 15:04:05"),
		AbsentOut:      attendance.AbsentOut.Format("2006-01-02 15:04:05"),
	}
}
