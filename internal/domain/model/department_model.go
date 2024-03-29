package model

type DepartmentResponse struct {
	DepartmentId   int64  `json:"department_id"`
	DepartmentName string `json:"department_name"`
	CreatedAt      string `json:"created_at"`
	CreatedBy      string `json:"created_by"`
	UpdatedAt      string `json:"updated_at"`
	UpdatedBy      string `json:"updated_by"`
}

type DepartmentCreateRequest struct {
	DepartmentName string `json:"department_name" validate:"required"`
	CreatedBy      string `json:"created_by"`
}

type DepartmentUpdateRequest struct {
	DepartmentId   string `json:"department_id" validate:"required"`
	DepartmentName string `json:"department_name" validate:"required"`
	UpdatedBy      string `json:"updated_by"`
}

type DepartmentGetByIdRequest struct {
	DepartmentId string `json:"-" validate:"required"`
}

type DepartmentDeleteRequest struct {
	DepartmentId string `json:"-" validate:"required"`
}

type DepartmentSearchRequest struct {
	DepartmentName string `json:"department_name"`
	Page           int    `json:"page" validate:"min=1"`
	Size           int    `json:"size" validate:"min=1,max=100"`
}
