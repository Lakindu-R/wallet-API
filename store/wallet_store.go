package store

import (
	"encoding/json"
	"wallet-api/models"

	"github.com/cockroachdb/pebble"
)

func SaveWallet(wallet models.Wallet) error {
	data, err := json.Marshal(wallet)

	if err != nil {
		return err
	}

	return DB.Set(
		[]byte(wallet.ID),
		data,
		pebble.Sync,
	)
}
// Get wallet
func GetWallet(id string) (*models.Wallet, error) {

	data, closer, err := DB.Get([]byte(id))

	if err != nil {
		return nil, err
	}

	defer closer.Close()

	var wallet models.Wallet

	err = json.Unmarshal(data, &wallet)

	if err != nil {
		return nil, err
	}

	return &wallet, nil
}