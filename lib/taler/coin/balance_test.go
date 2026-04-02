package libcoin

import (
	"testing"
)

func TestBalanceByCurrencyEmpty(t *testing.T) {
	actual, err := BalanceByCurrency(nil)
	if nil != err {
		t.Fatalf("unexpected error: %v", err)
	}

	if 0 != len(actual) {
		t.Errorf("the actual number of items is not as expected")
		t.Logf("EXPECTED: %d", 0)
		t.Logf("ACTUAL:   %d", len(actual))
	}
}

func TestBalanceByCurrencyAggregates(t *testing.T) {
	coins := []Coin{
		{CurrentAmount: Amount{Currency: "EUR", Value: 3, Fraction: 0}, ExchangeBaseURL: "https://a.example/"},
		{CurrentAmount: Amount{Currency: "EUR", Value: 2, Fraction: 0}, ExchangeBaseURL: "https://a.example/"},
		{CurrentAmount: Amount{Currency: "EUR", Value: 5, Fraction: 0}, ExchangeBaseURL: "https://b.example/"},
		{CurrentAmount: Amount{Currency: "KUDOS", Value: 12, Fraction: 50000000}, ExchangeBaseURL: "https://c.example/"},
		{CurrentAmount: Amount{Currency: "KUDOS", Value: 5, Fraction: 0}, ExchangeBaseURL: "https://d.example/"},
	}

	actual, err := BalanceByCurrency(coins)
	if nil != err {
		t.Fatalf("unexpected error: %v", err)
	}

	if 2 != len(actual) {
		t.Fatalf("the actual number of currencies is not as expected")
		t.Logf("EXPECTED: %d", 2)
		t.Logf("ACTUAL:   %d", len(actual))
	}

	// sorted alphabetically: EUR first, KUDOS second
	if "EUR" != actual[0].Currency {
		t.Errorf("the actual first currency is not as expected")
		t.Logf("EXPECTED: %s", "EUR")
		t.Logf("ACTUAL:   %s", actual[0].Currency)
	}
	if 10 != actual[0].Total.Value {
		t.Errorf("the actual EUR total 'value' is not as expected")
		t.Logf("EXPECTED: %d", 10)
		t.Logf("ACTUAL:   %d", actual[0].Total.Value)
	}
	if 0 != actual[0].Total.Fraction {
		t.Errorf("the actual EUR total 'fraction' is not as expected")
		t.Logf("EXPECTED: %d", 0)
		t.Logf("ACTUAL:   %d", actual[0].Total.Fraction)
	}

	if "KUDOS" != actual[1].Currency {
		t.Errorf("the actual second currency is not as expected")
		t.Logf("EXPECTED: %s", "KUDOS")
		t.Logf("ACTUAL:   %s", actual[1].Currency)
	}
	if 17 != actual[1].Total.Value {
		t.Errorf("the actual KUDOS total 'value' is not as expected")
		t.Logf("EXPECTED: %d", 17)
		t.Logf("ACTUAL:   %d", actual[1].Total.Value)
	}
	if 50000000 != actual[1].Total.Fraction {
		t.Errorf("the actual KUDOS total 'fraction' is not as expected")
		t.Logf("EXPECTED: %d", 50000000)
		t.Logf("ACTUAL:   %d", actual[1].Total.Fraction)
	}
}

func TestBalanceByCurrencyAndExchangeEmpty(t *testing.T) {
	actual, err := BalanceByCurrencyAndExchange(nil, "EUR")
	if nil != err {
		t.Fatalf("unexpected error: %v", err)
	}

	if 0 != len(actual) {
		t.Errorf("the actual number of items is not as expected")
		t.Logf("EXPECTED: %d", 0)
		t.Logf("ACTUAL:   %d", len(actual))
	}
}

func TestBalanceByCurrencyAndExchangeFiltersAndAggregates(t *testing.T) {
	coins := []Coin{
		{CurrentAmount: Amount{Currency: "EUR", Value: 3, Fraction: 0}, ExchangeBaseURL: "https://a.example/"},
		{CurrentAmount: Amount{Currency: "EUR", Value: 2, Fraction: 0}, ExchangeBaseURL: "https://a.example/"},
		{CurrentAmount: Amount{Currency: "EUR", Value: 5, Fraction: 0}, ExchangeBaseURL: "https://b.example/"},
		{CurrentAmount: Amount{Currency: "KUDOS", Value: 12, Fraction: 50000000}, ExchangeBaseURL: "https://c.example/"},
	}

	actual, err := BalanceByCurrencyAndExchange(coins, "EUR")
	if nil != err {
		t.Fatalf("unexpected error: %v", err)
	}

	if 2 != len(actual) {
		t.Fatalf("the actual number of exchanges is not as expected")
		t.Logf("EXPECTED: %d", 2)
		t.Logf("ACTUAL:   %d", len(actual))
	}

	// sorted alphabetically: a.example first, b.example second
	if "https://a.example/" != actual[0].ExchangeBaseURL {
		t.Errorf("the actual first exchange URL is not as expected")
		t.Logf("EXPECTED: %s", "https://a.example/")
		t.Logf("ACTUAL:   %s", actual[0].ExchangeBaseURL)
	}
	if 5 != actual[0].Total.Value {
		t.Errorf("the actual a.example total 'value' is not as expected")
		t.Logf("EXPECTED: %d", 5)
		t.Logf("ACTUAL:   %d", actual[0].Total.Value)
	}
	if 0 != actual[0].Total.Fraction {
		t.Errorf("the actual a.example total 'fraction' is not as expected")
		t.Logf("EXPECTED: %d", 0)
		t.Logf("ACTUAL:   %d", actual[0].Total.Fraction)
	}

	if "https://b.example/" != actual[1].ExchangeBaseURL {
		t.Errorf("the actual second exchange URL is not as expected")
		t.Logf("EXPECTED: %s", "https://b.example/")
		t.Logf("ACTUAL:   %s", actual[1].ExchangeBaseURL)
	}
	if 5 != actual[1].Total.Value {
		t.Errorf("the actual b.example total 'value' is not as expected")
		t.Logf("EXPECTED: %d", 5)
		t.Logf("ACTUAL:   %d", actual[1].Total.Value)
	}
	if 0 != actual[1].Total.Fraction {
		t.Errorf("the actual b.example total 'fraction' is not as expected")
		t.Logf("EXPECTED: %d", 0)
		t.Logf("ACTUAL:   %d", actual[1].Total.Fraction)
	}
}

func TestBalanceByCurrencyAndExchangeNoMatch(t *testing.T) {
	coins := []Coin{
		{CurrentAmount: Amount{Currency: "EUR", Value: 3, Fraction: 0}, ExchangeBaseURL: "https://a.example/"},
	}

	actual, err := BalanceByCurrencyAndExchange(coins, "KUDOS")
	if nil != err {
		t.Fatalf("unexpected error: %v", err)
	}

	if 0 != len(actual) {
		t.Errorf("the actual number of items is not as expected for non-matching currency")
		t.Logf("EXPECTED: %d", 0)
		t.Logf("ACTUAL:   %d", len(actual))
	}
}
