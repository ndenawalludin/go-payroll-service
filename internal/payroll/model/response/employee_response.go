package response

import "time"

type EmployeeResponse struct {
	ID         int64     `json:"id"`
	Code       string    `json:"code"`
	FullName   string    `json:"full_name"`
	Email      string    `json:"email"`
	BaseSalary int64     `json:"base_salary"`
	Allowance  int64     `json:"allowance"`
	IsActive   bool      `json:"is_active"`
	HireDate   time.Time `json:"hire_date"`
	CreateAt   time.Time `json:"create_at"`
	UpdateAt   time.Time `json:"update_at"`
}

type EmployeeListResponse []EmployeeResponse
