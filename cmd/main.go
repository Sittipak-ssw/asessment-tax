package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/Sittipak-ssw/pkg/db"

)

// type Allowance struct {
// 	AllowanceType string  `json:"allowanceType"`
// 	Amount        float64 `json:"amount"`
// }

// type taxRequest struct {
// 	TotalIncome       float64     `json:"totalIncome"`
// 	WHT               float64     `json:"wht"`
// 	Allowances        []Allowance `json:"allowances"`
// 	PersonalDeduction float64     `json:"personalDeduction"`
// }

// func calculateTax(totalIncome float64, wht float64, allowances []Allowance) (float64, []map[string]interface{}) {

// 	var totalDeductions float64
// 	for _, allowance := range allowances {
// 		if allowance.AllowanceType == "donation" {
// 			if allowance.Amount > 100000.0 {
// 				totalDeductions += 100000.0
// 			} else {
// 				totalDeductions += allowance.Amount
// 			}
// 		} else if allowance.AllowanceType == "personal" {
// 			if allowance.Amount < 10000.0 {
// 				totalDeductions += 10000.0
// 			} else {
// 				totalDeductions += allowance.Amount
// 			}
// 		} else if allowance.AllowanceType == "k-receipt" {
// 			if allowance.Amount > 50000.0 {
// 				totalDeductions += 50000.0
// 			} else {
// 				totalDeductions += allowance.Amount
// 			}
// 		}
// 	}

// 	totalDeductions += 60000.0

// 	taxableIncome := totalIncome - totalDeductions

// 	// Calculate tax based on tax brackets
// 	taxLevels := []map[string]interface{}{
// 		{"level": "0-150,000", "tax": 0.0},
// 		{"level": "150,001-500,000", "tax": 0.0},
// 		{"level": "500,001-1,000,000", "tax": 0.0},
// 		{"level": "1,000,001-2,000,000", "tax": 0.0},
// 		{"level": "2,000,001 ขึ้นไป", "tax": 0.0},
// 	}

// 	var tax float64

// 	switch {
// 	case taxableIncome <= 150000.0:
// 		tax = 0.0

// 		taxLevels[0]["tax"] = math.Round(tax)

// 	case taxableIncome <= 500000.0:
// 		tax = (taxableIncome - 150000.0) * 0.1
// 		if wht > 0 {
// 			tax -= wht
// 		}
// 		taxLevels[1]["tax"] = math.Round(tax)

// 	case taxableIncome <= 1000000.0:
// 		tax = 35000.0 + ((taxableIncome - 500000.0) * 0.15)
// 		if wht > 0 {
// 			tax -= wht
// 		}
// 		taxLevels[2]["tax"] = math.Round(tax)

// 	case taxableIncome <= 2000000.0:
// 		tax = 75000.0 + (taxableIncome-1000000.0)*0.2
// 		if wht > 0 {
// 			tax -= wht
// 		}
// 		taxLevels[3]["tax"] = math.Round(tax)

// 	default:
// 		tax = 175000 + (taxableIncome-2000000)*0.35
// 		if wht > 0 {
// 			tax -= wht
// 		}
// 		taxLevels[4]["tax"] = math.Round(tax)
// 	}

// 	if tax < 0.0 {
// 		tax = 0.0
// 	}

// 	return math.Round(tax), taxLevels
// }

// func calculateTaxRefund(tax float64, wht float64) float64 {

// 	taxRefund := wht - tax
// 	if taxRefund < 0 {
// 		taxRefund = 0
// 	}

// 	return math.Round(taxRefund)
// }

// func calculateTaxFromCSVHandler(c echo.Context) error {
// 	file, err := c.FormFile("taxFile")
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, "failed to get CSV file")
// 	}

// 	src, err := file.Open()
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open CSV file")
// 	}
// 	defer src.Close()

// 	reader := csv.NewReader(src)
// 	records, err := reader.ReadAll()
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError, "failed to read CSV file")
// 	}

// 	var taxes []map[string]interface{}
// 	for _, record := range records {

// 		totalIncome, err := strconv.ParseFloat(strings.TrimSpace(record[0]), 64)
// 		if err != nil {
// 			continue
// 		}

// 		wht, err := strconv.ParseFloat(strings.TrimSpace(record[1]), 64)
// 		if err != nil {
// 			continue
// 		}

// 		donation, err := strconv.ParseFloat(strings.TrimSpace(record[2]), 64)
// 		if err != nil {
// 			continue
// 		}

// 		tax, _ := calculateTax(totalIncome, wht, []Allowance{{AllowanceType: "donation", Amount: donation}})
// 		taxData := map[string]interface{}{
// 			"totalIncome": totalIncome,
// 			"tax":         tax,
// 		}
// 		taxes = append(taxes, taxData)
// 	}

// 	res := map[string]interface{}{
// 		"taxes": taxes,
// 	}

// 	return c.JSON(http.StatusOK, res)
// }

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/tax/calculations", db.calculateTaxHandler())
	e.POST("/admin/deductions/personal", db.setPersonalDeductionHandler())
	e.POST("/tax/calculations/upload-csv", db.calculateTaxFromCSVHandler())

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

// func calculateTaxHandler(c echo.Context) error {

// 	req := new(taxRequest)
// 	if err := c.Bind(req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
// 	}

// 	tax, taxLevels := calculateTax(req.TotalIncome, req.WHT, req.Allowances)

// 	var convertedTaxLevels []map[string]interface{}
// 	for _, level := range taxLevels {
// 		convertedTaxLevel := map[string]interface{}{
// 			"level": level["level"],
// 			"tax":   level["tax"],
// 		}
// 		convertedTaxLevels = append(convertedTaxLevels, convertedTaxLevel)
// 	}

// 	var taxRefund float64
// 	taxRefund = calculateTaxRefund(tax, req.WHT)

// 	res := map[string]interface{}{
// 		"tax":      tax,
// 		"taxLevel": convertedTaxLevels,
// 	}
// 	if taxRefund > 0 {
// 		res["taxRefund"] = taxRefund
// 	}
// 	return c.JSON(http.StatusOK, res)
// }

// func setPersonalDeductionHandler(c echo.Context) error {
// 	type personalDeductionRequest struct {
// 		Amount float64 `json:"amount"`
// 	}
// 	req := new(personalDeductionRequest)
// 	if err := c.Bind(req); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
// 	}

// 	res := map[string]interface{}{
// 		"personalDeduction": req.Amount,
// 	}

// 	return c.JSON(http.StatusOK, res)
// }
