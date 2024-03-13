package model

type AttendanceReportResponse struct {
	Date           string `json:"date"`
	EmployeeCode   string `json:"employee_code"`
	EmployeeName   string `json:"employee_name"`
	DepartmentName string `json:"department_name"`
	PositionName   string `json:"position_name"`
	LocationName   string `json:"location_name"`
	AbsentIn       string `json:"absent_in"`
	AbsentOut      string `json:"absent_out"`
}

type AttendanceReportSearchRequest struct {
	TimeInterval string `json:"time_interval"` // day, week, month
	Page         int    `json:"page" validate:"min=1"`
	Size         int    `json:"size" validate:"min=1,max=100"`
}
