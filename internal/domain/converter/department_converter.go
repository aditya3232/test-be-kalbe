package converter

import (
	"test-be-kalbe/internal/domain/entity"
	"test-be-kalbe/internal/domain/model"
)

func DepartmentToResponse(department *entity.Department) *model.DepartmentResponse {
	return &model.DepartmentResponse{
		DepartmentId:   department.DepartmentId,
		DepartmentName: department.DepartmentName,
		CreatedAt:      department.CreatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:      department.CreatedBy,
		UpdatedAt:      department.UpdatedAt.Format("2006-01-02 15:04:05"),
		UpdatedBy:      department.UpdatedBy,
	}
}
