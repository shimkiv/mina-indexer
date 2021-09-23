package indexing

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

type Identity struct {
	PublicKey string   `json:"public_key"`
	Name      string   `json:"name"`
	Fee       *float64 `json:"fee"`
}

func ReadIdentityFile(src string, handler func(Identity) error) error {
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

func importFromCSV(f *os.File, handler func(Identity) error) error {
	reader := csv.NewReader(f)

	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		currentIdentity := Identity{
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

func importFromJSON(f *os.File, handler func(Identity) error) error {
	identities := []Identity{}
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
