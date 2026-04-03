package libcoin

import (
	"encoding/binary"
	"testing"
)

func TestMarshalNBO(t *testing.T) {
	tests := []struct {
		Name     string
		Amount   Amount
		Expected [24]byte
	}{
		{
			Name:   "10.50 KUDOS",
			Amount: Amount{Currency: "KUDOS", Value: 10, Fraction: 50000000},
			Expected: func() [24]byte {
				var b [24]byte
				binary.BigEndian.PutUint64(b[0:8], 10)
				binary.BigEndian.PutUint32(b[8:12], 50000000)
				copy(b[12:24], "KUDOS")
				return b
			}(),
		},
		{
			Name:   "zero EUR",
			Amount: Amount{Currency: "EUR", Value: 0, Fraction: 0},
			Expected: func() [24]byte {
				var b [24]byte
				copy(b[12:24], "EUR")
				return b
			}(),
		},
		{
			Name:   "micropayment",
			Amount: Amount{Currency: "KUDOS", Value: 0, Fraction: 1},
			Expected: func() [24]byte {
				var b [24]byte
				binary.BigEndian.PutUint32(b[8:12], 1)
				copy(b[12:24], "KUDOS")
				return b
			}(),
		},
	}

	for testNumber, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			actual := test.Amount.MarshalNBO()
			if actual != test.Expected {
				t.Errorf("For test #%d, NBO mismatch.", testNumber)
				t.Logf("EXPECTED: %x", test.Expected)
				t.Logf("ACTUAL:   %x", actual)
			}
		})
	}
}

func TestMarshalNBOCurrencyZeroPadded(t *testing.T) {
	amt := Amount{Currency: "EUR", Value: 1, Fraction: 0}
	nbo := amt.MarshalNBO()

	// Currency field is bytes 12-23. "EUR" is 3 bytes, remaining 9 should be zero.
	for i := 15; i < 24; i++ {
		if nbo[i] != 0 {
			t.Errorf("expected zero at byte %d, got %d", i, nbo[i])
		}
	}
}

func TestMarshalNBORoundTripWithParseWireAmount(t *testing.T) {
	amt, err := ParseWireAmount("KUDOS:10.50")
	if nil != err {
		t.Fatalf("ParseWireAmount error: %v", err)
	}

	nbo := amt.MarshalNBO()

	// Verify fields from NBO
	value := binary.BigEndian.Uint64(nbo[0:8])
	fraction := binary.BigEndian.Uint32(nbo[8:12])

	if value != 10 {
		t.Errorf("value mismatch: expected 10, got %d", value)
	}
	if fraction != 50000000 {
		t.Errorf("fraction mismatch: expected 50000000, got %d", fraction)
	}
}
