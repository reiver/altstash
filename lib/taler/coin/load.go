package libcoin

import (
	"encoding/json"
	"os"
	"path/filepath"

	"codeberg.org/reiver/go-erorr"
	"codeberg.org/reiver/go-field"
)

// LoadFromDirectory reads all .talercoin files from the given directory path.
// Returns only coins with status "fresh".
func LoadFromDirectory(dirPath string) ([]Coin, error) {
	const starDotTalerCoin string = "*"+FileExtension // *.talercoin

	pattern := filepath.Join(dirPath, starDotTalerCoin)

	matches, err := filepath.Glob(pattern)
	if nil != err {
		return nil, erorr.Wrap(err, "could not glob for talercoin files",
			field.String("pattern", pattern),
		)
	}

	var coins []Coin

	for _, match := range matches {
		coin, err := loadCoinFile(match)
		if nil != err {
			return nil, erorr.Wrap(err, "could not load coin file: "+match)
		}

		if "fresh" == coin.Status {
			coins = append(coins, coin)
		}
	}

	return coins, nil
}

func loadCoinFile(filePath string) (Coin, error) {
	data, err := os.ReadFile(filePath)
	if nil != err {
		return Coin{}, erorr.Wrap(err, "could not read file",
			field.String("file-path", filePath),
		)
	}

	var coin Coin

	err = json.Unmarshal(data, &coin)
	if nil != err {
		return Coin{}, erorr.Wrap(err, "could not parse coin JSON",
			field.String("file-path", filePath),
		)
	}

	return coin, nil
}
