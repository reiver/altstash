package libcoin

// CoinSource tracks the provenance of a coin: how it was obtained.
// Type is one of "withdraw", "refresh", or "tip".
type CoinSource struct {
	Type       string `json:"type"`
	ReservePub string `json:"reserve_pub,omitempty"`
	OldCoinPub string `json:"old_coin_pub,omitempty"`
}
