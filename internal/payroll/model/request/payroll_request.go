package request

type GeneratePayrollRequest struct {
	PeriodCode string `json:"period_code" binding:"required"`
}
