package db

import (
	"testing"
)

func TestCalculateTax(t *testing.T) {
	
	testCases := []struct {
		totalIncome  float64
		wht          float64
		allowances   []Allowance
		expectedTax  float64
		expectedTaxLevels	[]map[string]interface{}
		expectedRefund float64
	}{
		{
			totalIncome:  500000.0,
			wht:          0.0,
			allowances:   []Allowance{{AllowanceType: "donation", Amount: 200000}},
			expectedTax:  19000,
			expectedTaxLevels: []map[string]interface{}{
				{
					"level": "0-150,000",
					"tax": 0,
				},
				{
					"level": "150,001-500,000",
					"tax": 19000,
				},
				{
					"level": "500,001-1,000,000",
					"tax": 0,
				},
				{
					"level": "1,000,001-2,000,000",
					"tax": 0,
				},
				{
					"level": "2,000,001 ขึ้นไป",
					"tax": 0,
				},
			},
			
			expectedRefund: 0,
		},
		
	}


	for _, tc := range testCases {
		tax, taxLevels, refund := calculateTax(tc.totalIncome, tc.wht, tc.allowances)

		if tax != tc.expectedTax {
			t.Errorf("Expected tax %f, but got %f", tc.expectedTax, tax)
		}

		if taxLevels != tc.expectedTaxLevels {
			t.Errorf("Expected tax levels %v, but got %v", tc.expectedTaxLevels, taxLevels)
		}
				
		if refund != tc.expectedRefund {
			t.Errorf("Expected tax refund %f, but got %f", tc.expectedRefund, refund)
		}
	}
}

func TestCalculateTaxRefund(t *testing.T) {

	testCases := []struct {
		tax         float64
		wht         float64
		expectedRefund float64
	}{
		{
			tax:  5000,
			wht:  10000,
			expectedRefund: 5000,
		},
		{
			tax:  20000,
			wht:  10000,
			expectedRefund: 0,
		},

	}


	for _, tc := range testCases {
		refund := calculateTaxRefund(tc.tax, tc.wht)


		if refund != tc.expectedRefund {
			t.Errorf("Expected tax refund %f, but got %f", tc.expectedRefund, refund)
		}
	}
}
