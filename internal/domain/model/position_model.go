package model

type PositionResponse struct {
	PositionId   int64  `json:"position_id"`
	DepartmentId int64  `json:"department_id"`
	PositionName string `json:"position_name"`
	CreatedAt    string `json:"created_at"`
	CreatedBy    string `json:"created_by"`
	UpdatedAt    string `json:"updated_at"`
	UpdatedBy    string `json:"updated_by"`
}

type PositionCreateRequest struct {
	DepartmentId string `json:"department_id" validate:"required"`
	PositionName string `json:"position_name" validate:"required"`
}

type PositionUpdateRequest struct {
	PositionId   string `json:"position_id" validate:"required"`
	DepartmentId string `json:"department_id" validate:"required"`
	PositionName string `json:"position_name" validate:"required"`
}

type PositionGetByIdRequest struct {
	PositionId string `json:"-" validate:"required"`
}

type PositionDeleteRequest struct {
	PositionId string `json:"-" validate:"required"`
}

type PositionSearchRequest struct {
	PositionName string `json:"position_name"`
	Page         int    `json:"page" validate:"min=1"`
	Size         int    `json:"size" validate:"min=1,max=100"`
}
