package cli

import (
	"encoding/csv"
	"errors"
	"io"
	"os"

	"github.com/figment-networks/mina-indexer/config"
	"github.com/figment-networks/mina-indexer/store"
	"github.com/sirupsen/logrus"
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

	db.SetDebugMode(true)

	f, err := os.Open(cfg.IdentityFile)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)

	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if err := db.Validators.UpdateIdentity(row[1], row[0]); err != nil {
			logrus.
				WithField("pk", row[1]).
				WithField("name", row[0]).
				WithError(err).
				Error("cant update validator identity")
			continue
		}

		logrus.WithField("pk", row[1]).WithField("name", row[0]).Info("identity updated")
	}

	return nil
}
