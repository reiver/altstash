package libcoin

// Coin represents a single Taler ecash coin stored locally.
// Field names follow Taler's wire format conventions.
type Coin struct {
	CoinPub         string     `json:"coin_pub"`
	CoinPriv        string     `json:"coin_priv"`
	DenomPubHash    string     `json:"denom_pub_hash"`
	DenomSig        DenomSig   `json:"denom_sig"`
	CurrentAmount   Amount     `json:"current_amount"`
	ExchangeBaseURL string     `json:"exchange_base_url"`
	BlindingKey     string     `json:"blinding_key"`
	CoinEvHash      string     `json:"coin_ev_hash"`
	Status          string     `json:"status"`
	CoinSource      CoinSource `json:"coin_source"`
}
