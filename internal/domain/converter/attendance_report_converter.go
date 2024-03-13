package converter

import (
	"test-be-kalbe/internal/domain/model"
)

func AttendanceReportToResponse(attendanceReport model.AttendanceReportResponse) *model.AttendanceReportResponse {
	return &model.AttendanceReportResponse{
		Date:           attendanceReport.Date,
		EmployeeCode:   attendanceReport.EmployeeCode,
		EmployeeName:   attendanceReport.EmployeeName,
		DepartmentName: attendanceReport.DepartmentName,
		PositionName:   attendanceReport.PositionName,
		LocationName:   attendanceReport.LocationName,
		AbsentIn:       attendanceReport.AbsentIn,
		AbsentOut:      attendanceReport.AbsentOut,
	}
}
