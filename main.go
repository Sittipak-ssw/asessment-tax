package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/Sittipak-ssw/assessment-tax/pkg/db" 
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        username, password, ok := c.Request().BasicAuth()
        if !ok {
            return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
        }
        if username != "adminTax" || password != "admin!"{
            return echo.NewHTTPError(http.StatusUnauthorized, "Username/Password incorrect.")
        }
        fmt.Println("Authorized passed")
        return next(c)
    }
}


func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/tax/calculations", db.CalculateTaxHandler)
	e.POST("/tax/calculations/upload-csv", db.CalculateTaxFromCSVHandler)

	e.POST("/admin/deductions/personal", AuthMiddleware(db.SetPersonalDeductionHandler))



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