package cli

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/figment-networks/mina-indexer/config"
	"github.com/figment-networks/mina-indexer/store"
	"github.com/sirupsen/logrus"
)

type identity struct {
	PublicKey string   `json:"public_key"`
	Name      string   `json:"name"`
	Fee       *float64 `json:"fee"`
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
		err = importFromCSV(f, handler)
	case ".json":
		err = importFromJSON(f, handler)
	default:
		err = errors.New("unsupported file extension")
	}

	return err
}

func importFromCSV(f *os.File, handler func(identity) error) error {
	reader := csv.NewReader(f)

	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		currentIdentity := identity{
			PublicKey: row[0],
			Name:      row[1],
		}
		if len(row) > 2 {
			fee, err := strconv.ParseFloat(row[2], 32)
			if err != nil {
				return err
			}
			currentIdentity.Fee = &fee
		}

		if err := handler(currentIdentity); err != nil {
			return err
		}
	}

	return nil
}

func importFromJSON(f *os.File, handler func(identity) error) error {
	identities := []identity{}
	if err := json.NewDecoder(f).Decode(&identities); err != nil {
		return err
	}

	for _, identityItem := range identities {
		if err := handler(identityItem); err != nil {
			return err
		}
	}

	return nil
}
