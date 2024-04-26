package db

import (
	"fmt"
	"math"
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
var kReceipt float64 = 50000.0

func calculateTax(totalIncome float64, wht float64, allowances []Allowance) (float64, []map[string]interface{}, float64) {

	var totalDeductions float64
	for _, allowance := range allowances {
		if allowance.AllowanceType == "donation" {
			if allowance.Amount > 100000.0 {
				totalDeductions += 100000.0
			} else {
				totalDeductions += allowance.Amount
			}
		} 
		if allowance.AllowanceType == "k-receipt" {
			totalDeductions += kReceipt
		}
	}

	taxableIncome := totalIncome - totalDeductions - personalDeduction

	taxLevels := []map[string]interface{}{
		{"level": "0-150,000", "tax": 0.0},
		{"level": "150,001-500,000", "tax": 0.0},
		{"level": "500,001-1,000,000", "tax": 0.0},
		{"level": "1,000,001-2,000,000", "tax": 0.0},
		{"level": "2,000,001 ขึ้นไป", "tax": 0.0},
	}

	var tax float64
	var taxFinal float64
	var taxRefund float64

	switch {
	case taxableIncome <= 150000.0:
		tax = 0.0
		taxLevels[0]["tax"] = tax

	case taxableIncome <= 500000.0:
		taxFinal = (taxableIncome - 150000.0) * 0.1
		taxRefund = calculateTaxRefund(taxFinal, wht)

		var tax_ float64
		if wht > 0 {
			tax_ = taxFinal - wht
			if tax_ > 0 {
				taxFinal = tax_
			}
		}

		taxLevels[1]["tax"] = tax

	case taxableIncome <= 1000000.0:
		tax = (taxableIncome - 500000.0) * 0.15
		taxFinal = tax + 35000.0 
		taxRefund = calculateTaxRefund(taxFinal, wht)

		var tax_ float64
		if wht > 0 {
			tax_ = taxFinal - wht
			if tax_ > 0 {
				taxFinal = tax_
			}
		}

		taxLevels[1]["tax"] = 35000.0
		taxLevels[2]["tax"] = tax

	case taxableIncome <= 2000000.0:
		tax = (taxableIncome-1000000.0)*0.2
		taxFinal = tax + 35000.0 + 75000.0
		taxRefund = calculateTaxRefund(taxFinal, wht)

		var tax_ float64
		if wht > 0 {
			tax_ = taxFinal - wht
			if tax_ > 0 {
				taxFinal = tax_
			}
		}

		taxLevels[1]["tax"] = 35000.0
		taxLevels[2]["tax"] = 75000.0
		taxLevels[3]["tax"] = tax

	default:
		tax = (taxableIncome - 2000000)*0.35
		taxFinal = tax + 35000.0 + 75000.0 + 200000.0
		fmt.Println(taxFinal)
		taxRefund = calculateTaxRefund(taxFinal, wht)
		fmt.Println(taxRefund)
		var tax_ float64
		if wht > 0 {
			tax_ = taxFinal - wht
			if tax_ > 0 {
				taxFinal = tax_
			}
		}

		taxLevels[1]["tax"] = 35000.0
		taxLevels[2]["tax"] = 75000.0
		taxLevels[3]["tax"] = 200000.0
		taxLevels[4]["tax"] = tax
		
	}

	if taxFinal < 0.0 {
		taxFinal = 0.0
	}

	return math.Round(taxFinal), taxLevels, taxRefund
}

func calculateTaxRefund(tax float64, wht float64) float64 {

	taxRefund := wht - tax
	if taxRefund < 0.0 {
		taxRefund = 0.0
	}

	return math.Round(taxRefund)
}

func setPersonalDeduction(newPersonalDeduction float64) {
	personalDeduction = newPersonalDeduction
}

func setKReceipt(newKReceipt float64) {
	kReceipt = newKReceipt
}