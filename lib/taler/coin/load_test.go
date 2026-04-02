package libcoin

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromDirectoryWithTestdata(t *testing.T) {
	// find testdata relative to this test file
	testdataDir := filepath.Join("..", "..", "..", "testdata")

	coins, err := LoadFromDirectory(testdataDir)
	if nil != err {
		t.Fatalf("unexpected error: %v", err)
	}

	if 5 != len(coins) {
		t.Fatalf("expected 5 fresh coins, got %d", len(coins))
	}

	// verify all coins are fresh
	for _, coin := range coins {
		if "fresh" != coin.Status {
			t.Errorf("expected status fresh, got %s for coin %s", coin.Status, coin.CoinPub)
		}
	}
}

func TestLoadFromDirectoryFiltersDormant(t *testing.T) {
	// create a temp directory with one fresh and one dormant coin
	tmpDir := t.TempDir()

	freshCoin := `{
		"coin_pub": "FRESH1",
		"coin_priv": "PRIV1",
		"denom_pub_hash": "HASH1",
		"denom_sig": {"cipher": "RSA", "rsa_signature": "SIG1"},
		"current_amount": {"currency": "EUR", "value": 5, "fraction": 0},
		"exchange_base_url": "https://exchange.example/",
		"blinding_key": "BLIND1",
		"coin_ev_hash": "EV1",
		"status": "fresh",
		"coin_source": {"type": "withdraw", "reserve_pub": "RES1"}
	}`

	dormantCoin := `{
		"coin_pub": "DORMANT1",
		"coin_priv": "PRIV2",
		"denom_pub_hash": "HASH2",
		"denom_sig": {"cipher": "RSA", "rsa_signature": "SIG2"},
		"current_amount": {"currency": "EUR", "value": 2, "fraction": 0},
		"exchange_base_url": "https://exchange.example/",
		"blinding_key": "BLIND2",
		"coin_ev_hash": "EV2",
		"status": "dormant",
		"coin_source": {"type": "withdraw", "reserve_pub": "RES2"}
	}`

	err := os.WriteFile(filepath.Join(tmpDir, "fresh.talercoin"), []byte(freshCoin), 0644)
	if nil != err {
		t.Fatalf("could not write fresh coin: %v", err)
	}

	err = os.WriteFile(filepath.Join(tmpDir, "dormant.talercoin"), []byte(dormantCoin), 0644)
	if nil != err {
		t.Fatalf("could not write dormant coin: %v", err)
	}

	coins, err := LoadFromDirectory(tmpDir)
	if nil != err {
		t.Fatalf("unexpected error: %v", err)
	}

	if 1 != len(coins) {
		t.Fatalf("expected 1 fresh coin, got %d", len(coins))
	}

	if "FRESH1" != coins[0].CoinPub {
		t.Errorf("expected fresh coin, got %s", coins[0].CoinPub)
	}
}

func TestLoadFromDirectoryEmptyDir(t *testing.T) {
	tmpDir := t.TempDir()

	coins, err := LoadFromDirectory(tmpDir)
	if nil != err {
		t.Fatalf("unexpected error: %v", err)
	}

	if 0 != len(coins) {
		t.Errorf("expected empty slice, got %d coins", len(coins))
	}
}
