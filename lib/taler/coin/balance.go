package libcoin

import (
	"sort"

	"codeberg.org/reiver/go-erorr"
)

// CurrencyBalance represents the total balance for a single currency,
// aggregated across all exchanges.
type CurrencyBalance struct {
	Currency string
	Total    Amount
}

// ExchangeBalance represents the balance for a single currency
// at a single exchange.
type ExchangeBalance struct {
	Currency        string
	ExchangeBaseURL string
	Total           Amount
}

// BalanceByCurrency aggregates coins into per-currency totals.
// Used for Screen 1 (balance list).
func BalanceByCurrency(coins []Coin) ([]CurrencyBalance, error) {
	totals := make(map[string]Amount)

	for _, coin := range coins {
		currency := coin.CurrentAmount.Currency

		existing, ok := totals[currency]
		if ok {
			sum, err := existing.Add(coin.CurrentAmount)
			if nil != err {
				return nil, erorr.Wrap(err, "could not aggregate balance for currency: "+currency)
			}
			totals[currency] = sum
		} else {
			totals[currency] = coin.CurrentAmount
		}
	}

	var balances []CurrencyBalance
	for currency, total := range totals {
		balances = append(balances, CurrencyBalance{
			Currency: currency,
			Total:    total,
		})
	}

	sort.Slice(balances, func(i, j int) bool {
		return balances[i].Currency < balances[j].Currency
	})

	return balances, nil
}

// BalanceByCurrencyAndExchange aggregates coins by currency + exchange.
// Used for Screen 2 (currency detail).
func BalanceByCurrencyAndExchange(coins []Coin, currency string) ([]ExchangeBalance, error) {
	totals := make(map[string]Amount)

	for _, coin := range coins {
		if currency != coin.CurrentAmount.Currency {
			continue
		}

		exchange := coin.ExchangeBaseURL

		existing, ok := totals[exchange]
		if ok {
			sum, err := existing.Add(coin.CurrentAmount)
			if nil != err {
				return nil, erorr.Wrap(err, "could not aggregate balance for exchange: "+exchange)
			}
			totals[exchange] = sum
		} else {
			totals[exchange] = coin.CurrentAmount
		}
	}

	var balances []ExchangeBalance
	for exchange, total := range totals {
		balances = append(balances, ExchangeBalance{
			Currency:        currency,
			ExchangeBaseURL: exchange,
			Total:           total,
		})
	}

	sort.Slice(balances, func(i, j int) bool {
		return balances[i].ExchangeBaseURL < balances[j].ExchangeBaseURL
	})

	return balances, nil
}
