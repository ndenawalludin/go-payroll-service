package service

import (
	"context"
	"go-payroll-service/cmd/internal/payroll/model/domain"
	"go-payroll-service/cmd/internal/payroll/model/request"
	"go-payroll-service/cmd/internal/payroll/repository"
)

type EmployeeService interface {
	List(ctx context.Context) ([]domain.Employee, error)
	Create(ctx context.Context, req request.CreateEmployeeRequest) (domain.Employee, error)
	GetByID(ctx context.Context, id int64) (domain.Employee, error)
	Update(ctx context.Context, id int64, req request.UpdateEmployeeRequest) (domain.Employee, error)
	Delete(ctx context.Context, id int64) error
}

type employeeService struct {
	repository repository.EmployeeRepository
}

func (s employeeService) List(ctx context.Context) ([]domain.Employee, error) {
	return s.repository.List(ctx)
}

func (s employeeService) Create(ctx context.Context, req request.CreateEmployeeRequest) (domain.Employee, error) {
	e := domain.Employee{
		Code:       req.Code,
		FullName:   req.FullName,
		Email:      req.Email,
		BaseSalary: req.BaseSalary,
		Allowance:  req.Allowance,
		HireDate:   req.HireDate,
	}

	return s.repository.Create(ctx, e)
}

func (s employeeService) GetByID(ctx context.Context, id int64) (domain.Employee, error) {
	return s.repository.GetByID(ctx, id)
}

func (s employeeService) Update(ctx context.Context, id int64, req request.UpdateEmployeeRequest) (domain.Employee, error) {
	current, err := s.repository.GetByID(ctx, id)

	if err != nil {
		return domain.Employee{}, err
	}
	if req.FullName != nil {
		current.FullName = *req.FullName
	}
	if req.Email != nil {
		current.Email = *req.Email
	}
	if req.BaseSalary != nil {
		current.BaseSalary = *req.BaseSalary
	}
	if req.Allowance != nil {
		current.Allowance = *req.Allowance
	}
	if req.HireDate != nil {
		current.HireDate = *req.HireDate
	}
	if req.IsActive != nil {
		current.IsActive = *req.IsActive
	}

	return s.repository.Update(ctx, current)
}

func (s employeeService) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}

func NewEmployeeService(repository repository.EmployeeRepository) EmployeeService {
	return &employeeService{repository: repository}
}
