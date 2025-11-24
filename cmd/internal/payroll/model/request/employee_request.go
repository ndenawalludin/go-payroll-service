package request

import "time"

type CreateEmployeeRequest struct {
	Code       string    `json:"code" binding:"required"`
	FullName   string    `json:"full_name" binding:"required"`
	Email      string    `json:"email" binding:"required,email"`
	BaseSalary int64     `json:"base_salary" binding:"required"`
	Allowance  int64     `json:"allowance"`
	HireDate   time.Time `json:"hire_date"`
}

type UpdateEmployeeRequest struct {
	FullName   *string    `json:"full_name"`
	Email      *string    `json:"email"`
	BaseSalary *int64     `json:"base_salary"`
	Allowance  *int64     `json:"allowance"`
	HireDate   *time.Time `json:"hire_date"`
	IsActive   *bool      `json:"is_active"`
}
