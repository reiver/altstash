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
