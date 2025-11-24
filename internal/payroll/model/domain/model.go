package domain

import "time"

type Employee struct {
	ID         int64     `db:"id"`
	Code       string    `db:"code"`
	FullName   string    `db:"full_name"`
	Email      string    `db:"email"`
	BaseSalary int64     `db:"base_salary"`
	Allowance  int64     `db:"allowance"`
	IsActive   bool      `db:"is_active"`
	HireDate   time.Time `db:"hire_date"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type PayrollPeriod struct {
	ID        int64     `db:"id"`
	Code      string    `db:"code"`
	StartDate string    `db:"start_date"`
	EndDate   string    `db:"end_date"`
	Closed    bool      `db:"closed"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Payslip struct {
	ID              int64 `db:"id"`
	EmployeeID      int64 `db:"employee_id"`
	PayrollPeriodID int64 `db:"payroll_period_id"`
	BaseSalary      int64 `db:"base_salary"`
	Allowance       int64 `db:"allowance"`
	Deduction       int64 `db:"deduction"`
	NetSalary       int64 `db:"net_salary"`
}

type PayslipWithEmployee struct {
	Payslip
	EmployeeName string `db:"employee_name"`
	PeriodCode   string `db:"period_code"`
}
