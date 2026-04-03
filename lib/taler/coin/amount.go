package libcoin

import (
	"fmt"

	"codeberg.org/reiver/go-erorr"
)

const fractionBase int64 = 100_000_000

// Amount represents a Taler-compatible monetary amount.
//
// The Fraction field is in units of 1/100,000,000 (10^8).
// So, the following:
//
//	libcoin.Amount{
//		Currency: "CAD",
//		Value:    12,
//	        Fraction: 34000000,
//	}
//
// Represents:
//
//	12.34 CAD
//
// And, the following:
//
//	libcoin.Amount{
//		Currency: "IRR",
//		Value:    3,
//	        Fraction: 14159265,
//	}
//
// Represents:
//
//	3.14159265 IRR
//
// And, the following:
//
//	libcoin.Amount{
//		Currency: "KRW",
//		Value:    0,
//	        Fraction: 00000001,
//	}
//
// Represents:
//
//	0.00000001 KRW
type Amount struct {
	Currency string `json:"currency"`
	Value    int64  `json:"value"`
	Fraction int64  `json:"fraction"`
}

// Add returns the sum of two [Amount].
//
// Both [Amount] must have the same currency, else Add returns an error (if the currencies do not match).
//
// Add handles fraction overflow: if the two fractions sum to >= 100,000,000, then it carries that over to value.
func (receiver Amount) Add(other Amount) (Amount, error) {
	if receiver.Currency != other.Currency {
		return Amount{}, erorr.Errorf("cannot add amounts with different currencies: %s and %s", receiver.Currency, other.Currency)
	}

	fraction := receiver.Fraction + other.Fraction
	value := receiver.Value + other.Value

	if fraction >= fractionBase {
		value += fraction / fractionBase
		fraction = fraction % fractionBase
	}

	return Amount{
		Currency: receiver.Currency,
		Value:    value,
		Fraction: fraction,
	}, nil
}

// String formats the amount with currency for standalone display (e.g., "5.23 KUDOS").
//
// String makes [Amount] fit the [fmt.Stringer] interface.
func (receiver Amount) String() string {
	return fmt.Sprintf("%s %s", receiver.FormatValue(), receiver.Currency)
}

// Sub returns the difference of two [Amount].
//
// Both [Amount] must have the same currency, else Sub returns an error (if the currencies do not match).
//
// Returns an error if the currencies do not match or the result would be negative.
// Handles fraction underflow: if receiver.Fraction < other.Fraction, borrows from value.
func (receiver Amount) Sub(other Amount) (Amount, error) {
	if receiver.Currency != other.Currency {
		return Amount{}, erorr.Errorf("cannot subtract amounts with different currencies: %s and %s", receiver.Currency, other.Currency)
	}

	value := receiver.Value - other.Value
	fraction := receiver.Fraction - other.Fraction

	if fraction < 0 {
		value--
		fraction += fractionBase
	}

	if value < 0 {
		return Amount{}, erorr.Errorf("insufficient funds: result would be negative")
	}

	return Amount{
		Currency: receiver.Currency,
		Value:    value,
		Fraction: fraction,
	}, nil
}

// GreaterThanOrEqual returns true if receiver >= other.
// Both must have the same currency; returns false if currencies differ.
func (receiver Amount) GreaterThanOrEqual(other Amount) bool {
	if receiver.Currency != other.Currency {
		return false
	}

	if receiver.Value != other.Value {
		return receiver.Value > other.Value
	}

	return receiver.Fraction >= other.Fraction
}

// FormatValue formats the value only, without currency (e.g., "5.23").
// Used in UI rows where currency is shown separately as the row title.
func (receiver Amount) FormatValue() string {
	if 0 == receiver.Fraction {
		return fmt.Sprintf("%d.00", receiver.Value)
	}

	fractStr := fmt.Sprintf("%08d", receiver.Fraction)

	// trim trailing zeros
	last := len(fractStr) - 1
	for last > 1 && '0' == fractStr[last] {
		last--
	}
	fractStr = fractStr[:last+1]

	return fmt.Sprintf("%d.%s", receiver.Value, fractStr)
}
