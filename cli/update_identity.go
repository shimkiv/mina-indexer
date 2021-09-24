package cli

import (
	"errors"

	"github.com/sirupsen/logrus"

	"github.com/figment-networks/mina-indexer/config"
	"github.com/figment-networks/mina-indexer/indexing"
	"github.com/figment-networks/mina-indexer/store"
)

func runUpdateIdentity(cfg *config.Config) error {
	if cfg.IdentityFile == "" {
		return errors.New("identity file is not provided")
	}

	db, err := store.New(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	return indexing.ReadIdentityFile(cfg.IdentityFile, func(item indexing.Identity) error {
		if item.Fee != nil {
			err := db.Validators.UpdateFee(item.PublicKey, *item.Fee)

			logrus.
				WithField("pk", item.PublicKey).
				WithField("fee", *item.Fee).
				WithError(err).
				Info("fee updated")

			if err != nil {
				return err
			}
		}

		err := db.Validators.UpdateIdentity(item.PublicKey, item.Name)

		logrus.
			WithField("pk", item.PublicKey).
			WithField("name", item.Name).
			WithError(err).
			Info("identity updated")

		return err
	})
}
