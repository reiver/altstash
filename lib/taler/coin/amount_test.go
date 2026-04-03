package libcoin

import (
	"testing"
)

func TestAmountAdd(t *testing.T) {
	tests := []struct {
		Name             string
		A, B             Amount
		ExpectedValue    int64
		ExpectedFraction int64
		ExpectedCurrency string
	}{
		{
			Name:             "no overflow",
			A:                Amount{Currency: "EUR", Value: 1, Fraction: 20000000},
			B:                Amount{Currency: "EUR", Value: 2, Fraction: 30000000},
			ExpectedValue:    3,
			ExpectedFraction: 50000000,
			ExpectedCurrency: "EUR",
		},
		{
			Name:             "with overflow",
			A:                Amount{Currency: "KUDOS", Value: 0, Fraction: 80000000},
			B:                Amount{Currency: "KUDOS", Value: 0, Fraction: 30000000},
			ExpectedValue:    1,
			ExpectedFraction: 10000000,
			ExpectedCurrency: "KUDOS",
		},
		{
			Name:             "whole numbers",
			A:                Amount{Currency: "EUR", Value: 3, Fraction: 0},
			B:                Amount{Currency: "EUR", Value: 5, Fraction: 0},
			ExpectedValue:    8,
			ExpectedFraction: 0,
			ExpectedCurrency: "EUR",
		},
		{
			Name:             "double overflow",
			A:                Amount{Currency: "EUR", Value: 0, Fraction: 99999999},
			B:                Amount{Currency: "EUR", Value: 0, Fraction: 99999999},
			ExpectedValue:    1,
			ExpectedFraction: 99999998,
			ExpectedCurrency: "EUR",
		},

		{
			Name:             "0 CAD + 0 CAD",
			A:                Amount{Currency: "CAD", Value: 0, Fraction: 00000000},
			B:                Amount{Currency: "CAD", Value: 0, Fraction: 00000000},
			ExpectedValue:    0,
			ExpectedFraction: 0,
			ExpectedCurrency: "CAD",
		},

		{
			Name:             "0.00000001 CAD + 0 CAD",
			A:                Amount{Currency: "CAD", Value: 0, Fraction: 00000001},
			B:                Amount{Currency: "CAD", Value: 0, Fraction: 00000000},
			ExpectedValue:    0,
			ExpectedFraction: 00000001,
			ExpectedCurrency: "CAD",
		},
		{
			Name:             "0 CAD + 0.00000001 CAD",
			A:                Amount{Currency: "CAD", Value: 0, Fraction: 00000000},
			B:                Amount{Currency: "CAD", Value: 0, Fraction: 00000001},
			ExpectedValue:    0,
			ExpectedFraction: 00000001,
			ExpectedCurrency: "CAD",
		},
		{
			Name:             "0.00000007 CAD + 0 CAD",
			A:                Amount{Currency: "CAD", Value: 0, Fraction: 00000007},
			B:                Amount{Currency: "CAD", Value: 0, Fraction: 00000000},
			ExpectedValue:    0,
			ExpectedFraction: 00000007,
			ExpectedCurrency: "CAD",
		},
		{
			Name:             "0 CAD + 0.00000007 CAD",
			A:                Amount{Currency: "CAD", Value: 0, Fraction: 00000000},
			B:                Amount{Currency: "CAD", Value: 0, Fraction: 00000007},
			ExpectedValue:    0,
			ExpectedFraction: 00000007,
			ExpectedCurrency: "CAD",
		},
		{
			Name:             "0.00000123 CAD + 0 CAD",
			A:                Amount{Currency: "CAD", Value: 0, Fraction: 00000123},
			B:                Amount{Currency: "CAD", Value: 0, Fraction: 00000000},
			ExpectedValue:    0,
			ExpectedFraction: 00000123,
			ExpectedCurrency: "CAD",
		},
		{
			Name:             "0 CAD + 0.00000123 CAD",
			A:                Amount{Currency: "CAD", Value: 0, Fraction: 00000000},
			B:                Amount{Currency: "CAD", Value: 0, Fraction: 00000123},
			ExpectedValue:    0,
			ExpectedFraction: 00000123,
			ExpectedCurrency: "CAD",
		},
		{
			Name:             "0.00000123 CAD + 0 CAD",
			A:                Amount{Currency: "CAD", Value: 3, Fraction: 14159265},
			B:                Amount{Currency: "CAD", Value: 0, Fraction: 00000000},
			ExpectedValue:    3,
			ExpectedFraction: 14159265,
			ExpectedCurrency: "CAD",
		},
		{
			Name:             "0 CAD + 0.00000123 CAD",
			A:                Amount{Currency: "CAD", Value: 0, Fraction: 00000000},
			B:                Amount{Currency: "CAD", Value: 3, Fraction: 14159265},
			ExpectedValue:    3,
			ExpectedFraction: 14159265,
			ExpectedCurrency: "CAD",
		},
	}

	for testNumber, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			actual, err := test.A.Add(test.B)
			if nil != err {
				t.Fatalf("For test #%d, unexpected error: %v", testNumber, err)
			}

			if test.ExpectedValue != actual.Value {
				t.Errorf("For test #%d, the actual 'value' is not as expected.", testNumber)
				t.Logf("EXPECTED: %d", test.ExpectedValue)
				t.Logf("ACTUAL:   %d", actual.Value)
			}
			if test.ExpectedFraction != actual.Fraction {
				t.Errorf("For test #%d, the actual 'fraction' is not as expected.", testNumber)
				t.Logf("EXPECTED: %d", test.ExpectedFraction)
				t.Logf("ACTUAL:   %d", actual.Fraction)
			}
			if test.ExpectedCurrency != actual.Currency {
				t.Errorf("For test #%d, the actual 'currency' is not as expected.", testNumber)
				t.Logf("EXPECTED: %s", test.ExpectedCurrency)
				t.Logf("ACTUAL:   %s", actual.Currency)
			}
		})
	}
}

func TestAmountSub(t *testing.T) {
	tests := []struct {
		Name             string
		A, B             Amount
		ExpectedValue    int64
		ExpectedFraction int64
	}{
		{
			Name:             "normal subtraction",
			A:                Amount{Currency: "EUR", Value: 5, Fraction: 50000000},
			B:                Amount{Currency: "EUR", Value: 2, Fraction: 30000000},
			ExpectedValue:    3,
			ExpectedFraction: 20000000,
		},
		{
			Name:             "fraction borrow",
			A:                Amount{Currency: "EUR", Value: 5, Fraction: 10000000},
			B:                Amount{Currency: "EUR", Value: 2, Fraction: 50000000},
			ExpectedValue:    2,
			ExpectedFraction: 60000000,
		},
		{
			Name:             "exact zero result",
			A:                Amount{Currency: "EUR", Value: 3, Fraction: 50000000},
			B:                Amount{Currency: "EUR", Value: 3, Fraction: 50000000},
			ExpectedValue:    0,
			ExpectedFraction: 0,
		},
		{
			Name:             "whole number subtraction",
			A:                Amount{Currency: "EUR", Value: 10, Fraction: 0},
			B:                Amount{Currency: "EUR", Value: 3, Fraction: 0},
			ExpectedValue:    7,
			ExpectedFraction: 0,
		},
		{
			Name:             "fraction only",
			A:                Amount{Currency: "EUR", Value: 0, Fraction: 80000000},
			B:                Amount{Currency: "EUR", Value: 0, Fraction: 30000000},
			ExpectedValue:    0,
			ExpectedFraction: 50000000,
		},
	}

	for testNumber, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			actual, err := test.A.Sub(test.B)
			if nil != err {
				t.Fatalf("For test #%d, unexpected error: %v", testNumber, err)
			}

			if test.ExpectedValue != actual.Value {
				t.Errorf("For test #%d, value mismatch.", testNumber)
				t.Logf("EXPECTED: %d", test.ExpectedValue)
				t.Logf("ACTUAL:   %d", actual.Value)
			}
			if test.ExpectedFraction != actual.Fraction {
				t.Errorf("For test #%d, fraction mismatch.", testNumber)
				t.Logf("EXPECTED: %d", test.ExpectedFraction)
				t.Logf("ACTUAL:   %d", actual.Fraction)
			}
		})
	}
}

func TestAmountSubInsufficientFunds(t *testing.T) {
	a := Amount{Currency: "EUR", Value: 2, Fraction: 0}
	b := Amount{Currency: "EUR", Value: 5, Fraction: 0}

	_, err := a.Sub(b)
	if nil == err {
		t.Error("expected error for insufficient funds, got nil")
	}
}

func TestAmountSubCurrencyMismatch(t *testing.T) {
	a := Amount{Currency: "EUR", Value: 5, Fraction: 0}
	b := Amount{Currency: "KUDOS", Value: 2, Fraction: 0}

	_, err := a.Sub(b)
	if nil == err {
		t.Error("expected error for currency mismatch, got nil")
	}
}

func TestAmountGreaterThanOrEqual(t *testing.T) {
	tests := []struct {
		Name     string
		A, B     Amount
		Expected bool
	}{
		{
			Name:     "equal amounts",
			A:        Amount{Currency: "EUR", Value: 5, Fraction: 50000000},
			B:        Amount{Currency: "EUR", Value: 5, Fraction: 50000000},
			Expected: true,
		},
		{
			Name:     "greater value",
			A:        Amount{Currency: "EUR", Value: 10, Fraction: 0},
			B:        Amount{Currency: "EUR", Value: 5, Fraction: 0},
			Expected: true,
		},
		{
			Name:     "lesser value",
			A:        Amount{Currency: "EUR", Value: 3, Fraction: 0},
			B:        Amount{Currency: "EUR", Value: 5, Fraction: 0},
			Expected: false,
		},
		{
			Name:     "equal value greater fraction",
			A:        Amount{Currency: "EUR", Value: 5, Fraction: 80000000},
			B:        Amount{Currency: "EUR", Value: 5, Fraction: 50000000},
			Expected: true,
		},
		{
			Name:     "equal value lesser fraction",
			A:        Amount{Currency: "EUR", Value: 5, Fraction: 30000000},
			B:        Amount{Currency: "EUR", Value: 5, Fraction: 50000000},
			Expected: false,
		},
		{
			Name:     "currency mismatch",
			A:        Amount{Currency: "EUR", Value: 100, Fraction: 0},
			B:        Amount{Currency: "KUDOS", Value: 1, Fraction: 0},
			Expected: false,
		},
	}

	for testNumber, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			actual := test.A.GreaterThanOrEqual(test.B)
			if test.Expected != actual {
				t.Errorf("For test #%d, expected %v, got %v", testNumber, test.Expected, actual)
			}
		})
	}
}

func TestAmountAddCurrencyMismatch(t *testing.T) {
	a := Amount{Currency: "EUR", Value: 1, Fraction: 0}
	b := Amount{Currency: "KUDOS", Value: 2, Fraction: 0}

	_, err := a.Add(b)
	if nil == err {
		t.Error("expected error when adding amounts with different currencies, got nil")
	}
}

func TestAmountFormatValue(t *testing.T) {
	tests := []struct {
		Name     string
		Amount   Amount
		Expected string
	}{
		{
			Name:     "whole number",
			Amount:   Amount{Currency: "EUR", Value: 5, Fraction: 0},
			Expected: "5.00",
		},
		{
			Name:     "with fraction",
			Amount:   Amount{Currency: "KUDOS", Value: 12, Fraction: 50000000},
			Expected: "12.50",
		},
		{
			Name:     "small fraction",
			Amount:   Amount{Currency: "EUR", Value: 1, Fraction: 23000000},
			Expected: "1.23",
		},
		{
			Name:     "micropayment",
			Amount:   Amount{Currency: "EUR", Value: 0, Fraction: 1},
			Expected: "0.00000001",
		},
		{
			Name:     "bigger micropayment",
			Amount:   Amount{Currency: "EUR", Value: 0, Fraction: 10},
			Expected: "0.0000001",
		},
		{
			Name:     "even bigger micropayment",
			Amount:   Amount{Currency: "EUR", Value: 0, Fraction: 100},
			Expected: "0.000001",
		},
	}

	for testNumber, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			actual := test.Amount.FormatValue()

			if test.Expected != actual {
				t.Errorf("For test #%d, the actual 'format-value' is not as expected.", testNumber)
				t.Logf("EXPECTED: %s", test.Expected)
				t.Logf("ACTUAL:   %s", actual)
			}
		})
	}
}

func TestAmountString(t *testing.T) {
	tests := []struct {
		Name     string
		Amount   Amount
		Expected string
	}{
		{
			Name:     "with fraction",
			Amount:   Amount{Currency: "KUDOS", Value: 5, Fraction: 23000000},
			Expected: "5.23 KUDOS",
		},
		{
			Name:     "whole number",
			Amount:   Amount{Currency: "EUR", Value: 10, Fraction: 0},
			Expected: "10.00 EUR",
		},
		{
			Name:     "zero CAD",
			Amount:   Amount{Currency: "CAD", Value: 0, Fraction: 0},
			Expected: "0.00 CAD",
		},

		{
			Name:     "10^-8 CAD (smallest)",
			Amount:   Amount{Currency: "CAD", Value: 0, Fraction: 1},
			Expected: "0.00000001 CAD",
		},
		{
			Name:     "10^-7 CAD",
			Amount:   Amount{Currency: "CAD", Value: 0, Fraction: 10},
			Expected: "0.0000001 CAD",
		},
		{
			Name:     "10^-6 CAD",
			Amount:   Amount{Currency: "CAD", Value: 0, Fraction: 100},
			Expected: "0.000001 CAD",
		},
		{
			Name:     "10^-5 CAD",
			Amount:   Amount{Currency: "CAD", Value: 0, Fraction: 1000},
			Expected: "0.00001 CAD",
		},
		{
			Name:     "10^-4 CAD",
			Amount:   Amount{Currency: "CAD", Value: 0, Fraction: 10000},
			Expected: "0.0001 CAD",
		},
		{
			Name:     "10^-3 CAD",
			Amount:   Amount{Currency: "CAD", Value: 0, Fraction: 100000},
			Expected: "0.001 CAD",
		},
		{
			Name:     "10^-2 CAD",
			Amount:   Amount{Currency: "CAD", Value: 0, Fraction: 1000000},
			Expected: "0.01 CAD",
		},
		{
			Name:     "10^-1 CAD",
			Amount:   Amount{Currency: "CAD", Value: 0, Fraction: 10000000},
			Expected: "0.10 CAD",
		},

		{
			Name:     "(2 + 10^-1) CAD",
			Amount:   Amount{Currency: "CAD", Value: 2, Fraction: 10000000},
			Expected: "2.10 CAD",
		},
		{
			Name:     "(3 + (2*10^-1)) CAD",
			Amount:   Amount{Currency: "CAD", Value: 3, Fraction: 20000000},
			Expected: "3.20 CAD",
		},

		{
			Name:     "Pi CAD",
			Amount:   Amount{Currency: "CAD", Value: 3, Fraction: 14159265},
			Expected: "3.14159265 CAD",
		},

		{
			Name:     "5.4321 FEDI",
			Amount:   Amount{Currency: "FEDI", Value: 5, Fraction: 43210000},
			Expected: "5.4321 FEDI",
		},
	}

	for testNumber, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			actual := test.Amount.String()

			if test.Expected != actual {
				t.Errorf("For test #%d, the actual 'string' is not as expected.", testNumber)
				t.Logf("EXPECTED: %s", test.Expected)
				t.Logf("ACTUAL:   %s", actual)
			}
		})
	}
}
