package cli

import (
	"errors"
	"log"
	"os"

	"github.com/figment-networks/mina-indexer/config"
	"github.com/figment-networks/mina-indexer/indexing"
	"github.com/figment-networks/mina-indexer/store"
)

func runCalculateRewards(cfg *config.Config) error {
	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	account := os.Getenv("REWARDS_ACCOUNT")
	if account == "" {
		return errors.New("account is not provided")
	}

	if err := db.Rewards.DeleteByValidator(account); err != nil {
		return err
	}

	if err := db.Rewards.DeleteByDelegator(account); err != nil {
		return err
	}

	search := &store.BlockSearch{
		Creator:   account,
		Canonical: "1",
		MinHeight: 0,
		Limit:     100,
		Sort:      "height",
		Order:     "asc",
	}

	for {
		if err := search.Validate(); err != nil {
			return err
		}

		log.Println("processing blocks from height:", search.MinHeight)

		blocks, err := db.Blocks.Search(search)
		if err != nil {
			return err
		}
		if len(blocks) == 0 {
			break
		}

		for _, block := range blocks {
			if err := indexing.RewardCalculation(db, block); err != nil {
				return err
			}
		}

		search.MinHeight = uint(blocks[len(blocks)-1].Height) + 1
	}

	return nil
}
