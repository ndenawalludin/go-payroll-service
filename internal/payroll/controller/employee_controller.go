package controller

import (
	"errors"
	"go-payroll-service/internal/payroll/model/request"
	"go-payroll-service/internal/payroll/model/response"
	"go-payroll-service/internal/payroll/service"
	"go-payroll-service/internal/payroll/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const employeeNotFound = "employee not found"

type EmployeeController struct {
	svc service.EmployeeService
}

func NewEmployeeController(svc service.EmployeeService) *EmployeeController {
	return &EmployeeController{svc: svc}
}

func (h *EmployeeController) RegisterRoutes(rg *gin.RouterGroup) {
	r := rg.Group("/employees")
	r.GET("", h.List)
	r.POST("", h.Create)
	r.GET("/:id", h.GetById)
	r.PUT("/:id", h.UpdateById)
	r.DELETE("/:id", h.DeleteById)
}

func (h *EmployeeController) List(c *gin.Context) {
	emps, err := h.svc.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var resp response.EmployeeListResponse
	for _, e := range emps {
		resp = append(resp, response.EmployeeResponse{
			ID:         e.ID,
			Code:       e.Code,
			FullName:   e.FullName,
			BaseSalary: e.BaseSalary,
			Allowance:  e.Allowance,
			IsActive:   e.IsActive,
			HireDate:   e.HireDate,
			CreateAt:   e.CreatedAt,
			UpdateAt:   e.UpdatedAt,
		})
	}

	if len(resp) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"data": resp})
	}

	c.JSON(http.StatusOK, resp)
}

func (h *EmployeeController) Create(c *gin.Context) {
	var req request.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	e, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := response.EmployeeResponse{
		ID:         e.ID,
		Code:       e.Code,
		FullName:   e.FullName,
		Email:      e.Email,
		BaseSalary: e.BaseSalary,
		Allowance:  e.Allowance,
		IsActive:   e.IsActive,
		HireDate:   e.HireDate,
		CreateAt:   e.CreatedAt,
		UpdateAt:   e.UpdatedAt,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *EmployeeController) GetById(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	e, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": employeeNotFound})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to fetch employee"})
		return
	}

	resp := response.EmployeeResponse{
		ID:         e.ID,
		Code:       e.Code,
		FullName:   e.FullName,
		Email:      e.Email,
		BaseSalary: e.BaseSalary,
		Allowance:  e.Allowance,
		IsActive:   e.IsActive,
		HireDate:   e.HireDate,
		CreateAt:   e.CreatedAt,
		UpdateAt:   e.UpdatedAt,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *EmployeeController) UpdateById(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var req request.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	e, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": employeeNotFound})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to update employee"})
		return
	}

	resp := response.EmployeeResponse{
		ID:         e.ID,
		Code:       e.Code,
		FullName:   e.FullName,
		Email:      e.Email,
		BaseSalary: e.BaseSalary,
		Allowance:  e.Allowance,
		IsActive:   e.IsActive,
		HireDate:   e.HireDate,
		CreateAt:   e.CreatedAt,
		UpdateAt:   e.UpdatedAt,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *EmployeeController) DeleteById(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := h.svc.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": employeeNotFound})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to delete employee"})
		return
	}
	c.Status(http.StatusNoContent)
}
