package converter

import (
	"test-be-kalbe/internal/domain/entity"
	"test-be-kalbe/internal/domain/model"
)

func PositionToResponse(position *entity.Position) *model.PositionResponse {
	return &model.PositionResponse{
		PositionId:   position.PositionId,
		DepartmentId: position.DepartmentId,
		PositionName: position.PositionName,
		CreatedAt:    position.CreatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:    position.CreatedBy,
		UpdatedAt:    position.UpdatedAt.Format("2006-01-02 15:04:05"),
		UpdatedBy:    position.UpdatedBy,
	}
}
