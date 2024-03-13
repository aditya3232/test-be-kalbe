package model

type AttendanceResponse struct {
	AttendanceId int64  `json:"attendance_id"`
	EmployeeId   int64  `json:"employee_id"`
	LocationId   int64  `json:"location_id"`
	AbsentIn     string `json:"absent_in"`
	AbsentOut    string `json:"absent_out"`
	CreatedAt    string `json:"created_at"`
	CreatedBy    string `json:"created_by"`
	UpdatedAt    string `json:"updated_at"`
	UpdatedBy    string `json:"updated_by"`
}

type AttendanceCreateRequest struct {
	EmployeeId int64  `json:"employee_id" validate:"required"`
	LocationId int64  `json:"location_id" validate:"required"`
	AbsentIn   string `json:"absent_in" validate:"required"`
	AbsentOut  string `json:"absent_out" validate:"required"`
	CreatedBy  string `json:"created_by"`
}

type AttendanceUpdateRequest struct {
	AttendanceId int64  `json:"attendance_id" validate:"required"`
	EmployeeId   int64  `json:"employee_id" validate:"required"`
	LocationId   int64  `json:"location_id" validate:"required"`
	AbsentIn     string `json:"absent_in" validate:"required"`
	AbsentOut    string `json:"absent_out" validate:"required"`
	UpdatedBy    string `json:"updated_by"`
}

type AttendanceGetByIdRequest struct {
	AttendanceId int64 `json:"attendance_id" validate:"required"`
}

type AttendanceDeleteRequest struct {
	AttendanceId int64 `json:"attendance_id" validate:"required"`
}

type AttendanceSearchRequest struct {
	EmployeeId int64 `json:"employee_id"`
	LocationId int64 `json:"location_id"`
	Page       int   `json:"page" validate:"min=1"`
	Size       int   `json:"size" validate:"min=1,max=100"`
}
