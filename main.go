package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/Sittipak-ssw/assessment-tax/pkg/db" 


)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/tax/calculations", db.CalculateTaxHandler)
	e.POST("/admin/deductions/personal", db.SetPersonalDeductionHandler)
	e.POST("/tax/calculations/upload-csv", db.CalculateTaxFromCSVHandler)
	e.POST("/admin/deductions/k-receipt", db.SetKReceiptHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	go func() {
		if err := e.Start(fmt.Sprintf(":%s", port)); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := e.Shutdown(context.Background()); err != nil {
		e.Logger.Fatal(err)
	}
}

