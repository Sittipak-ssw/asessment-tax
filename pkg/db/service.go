package db

import (
	"encoding/csv"
	"math"
	"net/http"
	"strconv"
	"strings"
	"fmt"

	"github.com/labstack/echo/v4"
)

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type taxRequest struct {
	TotalIncome       float64     `json:"totalIncome"`
	WHT               float64     `json:"wht"`
	Allowances        []Allowance `json:"allowances"`
  
}

var personalDeduction float64 = 60000.0
var totalDeductions float64 

func calculateTax(totalIncome float64, wht float64, allowances []Allowance) (float64, []map[string]interface{}) {

	for _, allowance := range allowances {
		if allowance.AllowanceType == "donation" {
			if allowance.Amount > 100000.0 {
				totalDeductions += 100000.0
			} else {
				totalDeductions += allowance.Amount
			}
		} 
	}

	taxableIncome := totalIncome - personalDeduction - totalDeductions

	taxLevels := []map[string]interface{}{
		{"level": "0-150,000", "tax": 0.0},
		{"level": "150,001-500,000", "tax": 0.0},
		{"level": "500,001-1,000,000", "tax": 0.0},
		{"level": "1,000,001-2,000,000", "tax": 0.0},
		{"level": "2,000,001 ขึ้นไป", "tax": 0.0},
	}

	var tax float64

	switch {
	case taxableIncome <= 150000.0:
		tax = 0.0

		taxLevels[0]["tax"] = math.Round(tax)

	case taxableIncome <= 500000.0:
		tax = (taxableIncome - 150000.0) * 0.1
		if wht > 0 {
			tax -= wht
		}
		taxLevels[1]["tax"] = math.Round(tax)

	case taxableIncome <= 1000000.0:
		tax = 35000.0 + ((taxableIncome - 500000.0) * 0.15)
		if wht > 0 {
			tax -= wht
		}
		taxLevels[2]["tax"] = math.Round(tax)

	case taxableIncome <= 2000000.0:
		tax = 75000.0 + (taxableIncome-1000000.0)*0.2
		if wht > 0 {
			tax -= wht
		}
		taxLevels[3]["tax"] = math.Round(tax)

	default:
		tax = 175000 + (taxableIncome-2000000)*0.35
		if wht > 0 {
			tax -= wht
		}
		taxLevels[4]["tax"] = math.Round(tax)
	}

	if tax < 0.0 {
		tax = 0.0
	}

	return math.Round(tax), taxLevels
}

func calculateTaxRefund(tax float64, wht float64) float64 {

	taxRefund := wht - tax
	if taxRefund < 0 {
		taxRefund = 0
	}

	return math.Round(taxRefund)
}

func calculateTaxFromCSVHandler(c echo.Context) error {
	file, err := c.FormFile("taxFile")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to get CSV file")
	}

	src, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to open CSV file")
	}
	defer src.Close()

	reader := csv.NewReader(src)
	records, err := reader.ReadAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to read CSV file")
	}

	var taxes []map[string]interface{}
	for _, record := range records {

		totalIncome, err := strconv.ParseFloat(strings.TrimSpace(record[0]), 64)
		if err != nil {
			continue
		}

		wht, err := strconv.ParseFloat(strings.TrimSpace(record[1]), 64)
		if err != nil {
			continue
		}

		donation, err := strconv.ParseFloat(strings.TrimSpace(record[2]), 64)
		if err != nil {
			continue
		}

		tax, _ := calculateTax(totalIncome, wht, []Allowance{{AllowanceType: "donation", Amount: donation}})
		taxData := map[string]interface{}{
			"totalIncome": totalIncome,
			"tax":         tax,
		}
		taxes = append(taxes, taxData)
	}

	res := map[string]interface{}{
		"taxes": taxes,
	}

	return c.JSON(http.StatusOK, res)
}

func setPersonalDeduction(newPersonalDeduction float64)  {
	personalDeduction = newPersonalDeduction
	fmt.Println("Personal deduction set to: ", personalDeduction)
}
