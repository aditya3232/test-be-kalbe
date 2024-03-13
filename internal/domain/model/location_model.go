package model

type LocationResponse struct {
	LocationId   int64  `json:"location_id"`
	LocationName string `json:"location_name"`
	CreatedAt    string `json:"created_at"`
	CreatedBy    string `json:"created_by"`
	UpdatedAt    string `json:"updated_at"`
	UpdatedBy    string `json:"updated_by"`
}

type LocationCreateRequest struct {
	LocationName string `json:"location_name" validate:"required"`
	CreatedBy    string `json:"created_by"`
}

type LocationUpdateRequest struct {
	LocationId   string `json:"location_id" validate:"required"`
	LocationName string `json:"location_name" validate:"required"`
	UpdatedBy    string `json:"updated_by"`
}

type LocationGetByIdRequest struct {
	LocationId string `json:"-" validate:"required"`
}

type LocationDeleteRequest struct {
	LocationId string `json:"-" validate:"required"`
}

type LocationSearchRequest struct {
	LocationName string `json:"location_name"`
	Page         int    `json:"page" validate:"min=1"`
	Size         int    `json:"size" validate:"min=1,max=100"`
}
