package db

import (
	"net/http"
	
	
	"github.com/labstack/echo/v4"
)

func CalculateTaxHandler(c echo.Context) error {

	req := new(taxRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}

	tax, _, taxRefund := calculateTax(req.TotalIncome, req.WHT, req.Allowances)

	res := map[string]interface{}{
		"tax":      tax,
	}
	if taxRefund > 0 {
		res["taxRefund"] = taxRefund
	}
	return c.JSON(http.StatusOK, res)
}
