package main

import (
	"go-payroll-service/internal/config"
	"go-payroll-service/internal/db"
	controller2 "go-payroll-service/internal/payroll/controller"
	repository2 "go-payroll-service/internal/payroll/repository"
	service2 "go-payroll-service/internal/payroll/service"
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
	empRepo := repository2.NewEmployeeRepository(dbConn)
	payrollRepo := repository2.NewPayrollRepository(dbConn)

	empService := service2.NewEmployeeService(empRepo)
	payrollService := service2.NewPayrollService(empRepo, payrollRepo)

	empController := controller2.NewEmployeeController(empService)
	payrollController := controller2.NewPayrollController(payrollService)

	api := r.Group("/api/v1")
	empController.RegisterRoutes(api)
	payrollController.RegisterRoutes(api)

	addr := ":" + cfg.HTTPPort
	log.Println("Listening on " + addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
