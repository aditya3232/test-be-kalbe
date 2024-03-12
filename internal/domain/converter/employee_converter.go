package converter

import (
	"test-be-kalbe/internal/domain/entity"
	"test-be-kalbe/internal/domain/model"
)

func EmployeeToResponse(employee *entity.Employee) *model.EmployeeResponse {
	return &model.EmployeeResponse{
		EmployeeId:   employee.EmployeeId,
		EmployeeCode: employee.EmployeeCode,
		EmployeeName: employee.EmployeeName,
		Password:     employee.Password,
		DepartmentId: employee.DepartmentId,
		PositionId:   employee.PositionId,
		Superior:     employee.Superior,
		CreatedAt:    employee.CreatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy:    employee.CreatedBy,
		UpdatedAt:    employee.UpdatedAt.Format("2006-01-02 15:04:05"),
		UpdatedBy:    employee.UpdatedBy,
	}
}
