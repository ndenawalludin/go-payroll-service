package repository

import (
	"context"
	"database/sql"
	"errors"
	"go-payroll-service/cmd/internal/payroll/model/domain"
	"go-payroll-service/cmd/internal/payroll/util"
	"time"
)

type PayrollRepository interface {
	GetOrCreatePeriod(ctx context.Context, code string, start, end time.Time) (domain.PayrollPeriod, error)
	CreatePayslip(ctx context.Context, p domain.Payslip) (domain.Payslip, error)
	ListPayslipByPeriodCode(ctx context.Context, periodCode string) ([]domain.PayslipWithEmployee, error)
}

type payrollRepository struct {
	db *sql.DB
}

func (r payrollRepository) GetOrCreatePeriod(ctx context.Context, code string, start, end time.Time) (domain.PayrollPeriod, error) {
	var p domain.PayrollPeriod
	err := r.db.QueryRowContext(ctx, `
		SELECT id, code, start_date, end_date, closed, created_at, updated_at
		FROM payroll_periods
		WHERE code = $1`, code,
	).Scan(&p.ID, &p.Code, &p.StartDate, &p.EndDate, &p.Closed, &p.CreatedAt, &p.UpdatedAt)

	if err == nil {
		return p, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return domain.PayrollPeriod{}, err
	}

	now := time.Now()
	p.Code = code
	p.StartDate = start.Format("2006-01-02")
	p.EndDate = end.Format("2006-01-02")
	p.Closed = false
	p.CreatedAt = now
	p.UpdatedAt = now

	err = r.db.QueryRowContext(ctx, `
		INSERT INTO payroll_periods(code, start_date, end_date, closed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		p.Code, p.StartDate, p.EndDate, p.Closed, p.CreatedAt, p.UpdatedAt).Scan(&p.ID)
	if err != nil {
		return domain.PayrollPeriod{}, err
	}
	return p, nil
}

func (r payrollRepository) CreatePayslip(ctx context.Context, p domain.Payslip) (domain.Payslip, error) {
	err := r.db.QueryRowContext(ctx, `
			INSERT INTO payslips(employee_id, payroll_period_id, base_salary, allowance, deduction, net_salary)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`,
		p.EmployeeID, p.PayrollPeriodID, p.BaseSalary, p.Allowance, p.Deduction, p.NetSalary,
	).Scan(&p.ID)
	if err != nil {
		return domain.Payslip{}, err
	}
	return p, nil
}

func (r payrollRepository) ListPayslipByPeriodCode(ctx context.Context, periodCode string) ([]domain.PayslipWithEmployee, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT ps.id,
		       ps.employee_id,
		       ps.payroll_period_id,
		       ps.base_salary,
		       ps.allowance,
		       ps.deduction,
		       ps.net_salary,
		       e.full_name as employee_name,
		       pp.code as period_code
		FROM payslips ps
		JOIN employees e ON e.id = ps.employee_id
		JOIN payroll_periods pp ON pp.id = ps.payroll_period_id
		WHERE pp.code = $1
		ORDER BY e.full_name`, periodCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.PayslipWithEmployee
	for rows.Next() {
		var p domain.PayslipWithEmployee
		if err := rows.Scan(
			&p.ID,
			&p.EmployeeID,
			&p.PayrollPeriodID,
			&p.BaseSalary,
			&p.Allowance,
			&p.Deduction,
			&p.NetSalary,
			&p.EmployeeName,
			&p.PeriodCode,
		); err != nil {
			return nil, err
		}
		result = append(result, p)
	}

	if len(result) == 0 {
		return nil, util.ErrNotFound
	}
	return result, rows.Err()
}

func NewPayrollRepository(db *sql.DB) PayrollRepository {
	return &payrollRepository{
		db: db,
	}
}
