package worker

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/config"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/types"
	"github.com/figment-networks/mina-indexer/store"
)

func RunInit(cfg *config.Config, db *store.Store) error {
	n, err := db.Accounts.Count()
	if err != nil {
		return err
	}
	if n > 0 {
		return errors.New("accounts table should be empty before the genesis import")
	}

	genesis, err := graph.ReadGenesisFile(cfg.GenesisFile)
	if err != nil {
		return err
	}
	log.Info("genesis accounts found:", len(genesis.Ledger.Accounts))

	for _, a := range genesis.Ledger.Accounts {
		acc := model.Account{
			PublicKey:      a.PK,
			Balance:        types.NewInt64Amount(0),
			BalanceUnknown: types.NewInt64Amount(0),
			Stake:          types.NewFloatAmount(a.Balance),
			Delegate:       a.Delegate,
			StartHeight:    0,
			StartTime:      genesis.Config.Timestamp,
			LastHeight:     0,
			LastTime:       genesis.Config.Timestamp,
		}

		log.WithField("pk", a.PK).Info("importing account")
		if err := db.Accounts.Create(&acc); err != nil {
			return err
		}
	}

	return nil
}
