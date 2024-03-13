package converter

import (
	"test-be-kalbe/internal/domain/model"
	"time"
)

func AttendanceReportToResponse(attendanceReport model.AttendanceReportResponse) *model.AttendanceReportResponse {
	parsedDate, _ := time.Parse(time.RFC3339, attendanceReport.Date)
	parsedAbsentIn, _ := time.Parse(time.RFC3339, attendanceReport.AbsentIn)
	parsedAbsentOut, _ := time.Parse(time.RFC3339, attendanceReport.AbsentOut)

	return &model.AttendanceReportResponse{
		Date:           parsedDate.Format("2006-01-02 15:04:05"),
		EmployeeCode:   attendanceReport.EmployeeCode,
		EmployeeName:   attendanceReport.EmployeeName,
		DepartmentName: attendanceReport.DepartmentName,
		PositionName:   attendanceReport.PositionName,
		LocationName:   attendanceReport.LocationName,
		AbsentIn:       parsedAbsentIn.Format("2006-01-02 15:04:05"),
		AbsentOut:      parsedAbsentOut.Format("2006-01-02 15:04:05"),
	}
}
