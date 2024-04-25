package db

import (
	"net/http"
	"encoding/csv"
    "strconv"
    "strings"

	"github.com/labstack/echo/v4"
)

func CalculateTaxHandler(c echo.Context) error {

	req := new(taxRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}

	tax, taxLevels := calculateTax(req.TotalIncome, req.WHT, req.Allowances)

	var convertedTaxLevels []map[string]interface{}
	for _, level := range taxLevels {
		convertedTaxLevel := map[string]interface{}{
			"level": level["level"],
			"tax":   level["tax"],
		}
		convertedTaxLevels = append(convertedTaxLevels, convertedTaxLevel)
	}

	var taxRefund float64
	taxRefund = calculateTaxRefund(tax, req.WHT)

	res := map[string]interface{}{
		"tax":      tax,
		"taxLevel": convertedTaxLevels,
	}
	if taxRefund > 0 {
		res["taxRefund"] = taxRefund
	}
	return c.JSON(http.StatusOK, res)
}

func SetPersonalDeductionHandler(c echo.Context) error {
	type personalDeductionRequest struct {
		Amount float64 `json:"amount"`
	}
	req := new(personalDeductionRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}

	setPersonalDeduction(req.Amount)

	res := map[string]interface{}{
		"personalDeduction": req.Amount,
	}

	return c.JSON(http.StatusOK, res)
}

func CalculateTaxFromCSVHandler(c echo.Context) error {
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