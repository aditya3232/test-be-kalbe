package model

type EmployeeResponse struct {
	EmployeeId   int64  `json:"employee_id"`
	EmployeeCode string `json:"employee_code"`
	EmployeeName string `json:"employee_name"`
	Password     string `json:"password"`
	DepartmentId int64  `json:"department_id"`
	PositionId   int64  `json:"position_id"`
	Superior     int64  `json:"superior"`
	CreatedAt    string `json:"created_at"`
	CreatedBy    string `json:"created_by"`
	UpdatedAt    string `json:"updated_at"`
	UpdatedBy    string `json:"updated_by"`
	Token        string `json:"token"`
}

type EmployeeCreateRequest struct {
	EmployeeName string `json:"employee_name" validate:"required"`
	Password     string `json:"password" validate:"required"`
	DepartmentId string `json:"department_id" validate:"required"`
	PositionId   string `json:"position_id" validate:"required"`
	Superior     string `json:"superior" validate:"required"`
	CreatedBy    string `json:"created_by"`
}

type EmployeeUpdateRequest struct {
	EmployeeId   string `json:"employee_id" validate:"required"`
	EmployeeName string `json:"employee_name" validate:"required"`
	Password     string `json:"password" validate:"required"`
	DepartmentId string `json:"department_id" validate:"required"`
	PositionId   string `json:"position_id" validate:"required"`
	Superior     string `json:"superior" validate:"required"`
	Token        string `json:"token"`
	UpdatedBy    string `json:"updated_by"`
}

type EmployeeGetByIdRequest struct {
	EmployeeId string `json:"-" validate:"required"`
}

type EmployeeDeleteRequest struct {
	EmployeeId string `json:"-" validate:"required"`
}

type EmployeeSearchRequest struct {
	EmployeeName string `json:"employee_name"`
	Page         int    `json:"page" validate:"min=1"`
	Size         int    `json:"size" validate:"min=1,max=100"`
}
