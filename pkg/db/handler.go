package db

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func calculateTaxHandler(c echo.Context) error {

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

func setPersonalDeductionHandler(c echo.Context) error {
	type personalDeductionRequest struct {
		Amount float64 `json:"amount"`
	}
	req := new(personalDeductionRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}

	res := map[string]interface{}{
		"personalDeduction": req.Amount,
	}

	return c.JSON(http.StatusOK, res)
}
