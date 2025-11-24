package main

import (
	"go-payroll-service/cmd/internal/config"
	"go-payroll-service/cmd/internal/db"
	"go-payroll-service/cmd/internal/payroll/controller"
	"go-payroll-service/cmd/internal/payroll/repository"
	"go-payroll-service/cmd/internal/payroll/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	dbConn, err := db.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	r := gin.Default()

	//dependency injection
	empRepo := repository.NewEmployeeRepository(dbConn)
	payrollRepo := repository.NewPayrollRepository(dbConn)

	empService := service.NewEmployeeService(empRepo)
	payrollService := service.NewPayrollService(empRepo, payrollRepo)

	empController := controller.NewEmployeeController(empService)
	payrollController := controller.NewPayrollController(payrollService)

	api := r.Group("/api/v1")
	empController.RegisterRoutes(api)
	payrollController.RegisterRoutes(api)

	addr := ":" + cfg.HTTPPort
	log.Println("Listening on " + addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
