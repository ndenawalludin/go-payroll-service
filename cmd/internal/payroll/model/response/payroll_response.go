package response

type PayslipResponse struct {
	ID           int64  `json:"id"`
	EmployeeID   int64  `json:"employee_id"`
	EmployeeName string `json:"employee_name"`
	PeriodCode   string `json:"period_code"`
	BaseSalary   int64  `json:"base_salary"`
	Allowance    int64  `json:"allowance"`
	Deduction    int64  `json:"deduction"`
	NetSalary    int64  `json:"net_salary"`
}

type GeneratePayrollResponse struct {
	PeriodCode   string `json:"period_code"`
	TotalPayslip int    `json:"total_payslip"`
}

type PayslipListResponse []PayslipResponse
