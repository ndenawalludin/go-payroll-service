package service

import (
	"context"
	"go-payroll-service/cmd/internal/payroll/model/domain"
	"go-payroll-service/cmd/internal/payroll/model/request"
	"go-payroll-service/cmd/internal/payroll/repository"
	"time"
)

type PayrollService interface {
	GeneratePayroll(ctx context.Context, req request.GeneratePayrollRequest) (int, error)
	ListPayslips(ctx context.Context, periodCode string) ([]domain.PayslipWithEmployee, error)
}

type payrollService struct {
	employeeRepository repository.EmployeeRepository
	payrollRepository  repository.PayrollRepository
}

func (s payrollService) GeneratePayroll(ctx context.Context, req request.GeneratePayrollRequest) (int, error) {
	periodCode := req.PeriodCode

	//example the period is 30 days
	start := time.Now().AddDate(0, 0, -30)
	end := time.Now()

	period, err := s.payrollRepository.GetOrCreatePeriod(ctx, periodCode, start, end)
	if err != nil {
		return 0, err
	}

	employees, err := s.employeeRepository.List(ctx)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, e := range employees {
		if !e.IsActive {
			continue
		}
		base := e.BaseSalary
		allow := e.Allowance
		deduction := int64(0) // this can be chance with pajak
		net := base + allow - deduction

		p := domain.Payslip{
			EmployeeID:      e.ID,
			PayrollPeriodID: period.ID,
			BaseSalary:      base,
			Allowance:       allow,
			Deduction:       deduction,
			NetSalary:       net,
		}
		if _, err := s.payrollRepository.CreatePayslip(ctx, p); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (s payrollService) ListPayslips(ctx context.Context, periodCode string) ([]domain.PayslipWithEmployee, error) {
	return s.payrollRepository.ListPayslipByPeriodCode(ctx, periodCode)
}

func NewPayrollService(employeeRepository repository.EmployeeRepository, payrollRepository repository.PayrollRepository) PayrollService {
	return &payrollService{
		employeeRepository: employeeRepository,
		payrollRepository:  payrollRepository,
	}
}
