package converter

import (
	"test-be-kalbe/internal/domain/entity"
	"test-be-kalbe/internal/domain/model"
)

func LocationToResponse(location *entity.Location) *model.LocationResponse {
	return &model.LocationResponse{
		LocationId:   location.LocationId,
		LocationName: location.LocationName,
		CreatedAt:    location.CreatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:    location.CreatedBy,
		UpdatedAt:    location.UpdatedAt.Format("2006-01-02 15:04:05"),
		UpdatedBy:    location.UpdatedBy,
	}
}
