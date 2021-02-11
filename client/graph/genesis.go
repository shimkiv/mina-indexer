package graph

import (
	"encoding/json"
	"os"
	"time"
)

type Genesis struct {
	Config struct {
		Timestamp time.Time `json:"genesis_state_timestamp"`
	} `json:"genesis"`
	Ledger struct {
		Name        string           `json:"name"`
		NumAccounts int              `json:"num_accounts"`
		Accounts    []GenesisAccount `json:"accounts"`
	} `json:"ledger"`
}

type GenesisAccount struct {
	PK       string  `json:"pk"`
	SK       *string `json:"sk"`
	Balance  string  `json:"balance"`
	Delegate *string `json:"delegate"`
}

// ReadGenesisFile reads and returns genesis file records
func ReadGenesisFile(path string) (*Genesis, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	genesis := &Genesis{}

	if err := json.NewDecoder(f).Decode(genesis); err != nil {
		return nil, err
	}
	return genesis, nil
}
