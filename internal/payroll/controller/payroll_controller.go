package controller

import (
	"errors"
	"go-payroll-service/internal/payroll/model/request"
	"go-payroll-service/internal/payroll/model/response"
	"go-payroll-service/internal/payroll/service"
	"go-payroll-service/internal/payroll/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PayrollController struct {
	svc service.PayrollService
}

func NewPayrollController(svc service.PayrollService) *PayrollController {
	return &PayrollController{svc: svc}
}

func (h *PayrollController) RegisterRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/payroll")
	r.POST("/generate", h.Generate)
	r.GET("/payslips/:periodCode", h.ListPayslips)
}

func (h *PayrollController) Generate(c *gin.Context) {
	var req request.GeneratePayrollRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := h.svc.GeneratePayroll(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate payroll"})
		return
	}

	resp := response.GeneratePayrollResponse{
		PeriodCode:   req.PeriodCode,
		TotalPayslip: count,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *PayrollController) ListPayslips(c *gin.Context) {
	periodCode := c.Param("periodCode")

	list, err := h.svc.ListPayslips(c.Request.Context(), periodCode)

	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no payslips found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list payslips"})
		return
	}

	var resp response.PayslipListResponse
	for _, p := range list {
		resp = append(resp, response.PayslipResponse{
			ID:           p.ID,
			EmployeeID:   p.EmployeeID,
			EmployeeName: p.EmployeeName,
			PeriodCode:   p.PeriodCode,
			BaseSalary:   p.BaseSalary,
			Allowance:    p.Allowance,
			Deduction:    p.Deduction,
			NetSalary:    p.NetSalary,
		})
	}
	c.JSON(http.StatusOK, resp)
}
