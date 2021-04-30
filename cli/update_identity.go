package cli

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/figment-networks/mina-indexer/config"
	"github.com/figment-networks/mina-indexer/store"
	"github.com/sirupsen/logrus"
)

type identity struct {
	PublicKey string `json:"public_key"`
	Name      string `json:"name"`
}

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

	return readIdentityFile(cfg.IdentityFile, func(item identity) error {
		err := db.Validators.UpdateIdentity(item.PublicKey, item.Name)

		logrus.
			WithField("pk", item.PublicKey).
			WithField("name", item.Name).
			WithError(err).
			Info("identity updated")

		return err
	})
}

func readIdentityFile(src string, handler func(identity) error) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	switch filepath.Ext(filepath.Base(src)) {
	case ".csv":
		reader := csv.NewReader(f)

		for {
			row, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}

			if err := handler(identity{row[0], row[1]}); err != nil {
				return err
			}
		}
	case ".json":
		identities := []identity{}
		if err := json.NewDecoder(f).Decode(&identities); err != nil {
			return err
		}
		for _, identityItem := range identities {
			if err := handler(identityItem); err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupports file extension")
	}

	return nil
}
