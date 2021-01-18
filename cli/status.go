package cli

import (
	"fmt"
	"net/http"

	"github.com/figment-networks/mina-indexer/coda"
	"github.com/figment-networks/mina-indexer/config"
)

func startStatus(cfg *config.Config) error {
	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	fmt.Println("=== Height Indexing ===")

	height, err := db.Heights.LastSuccessful()
	if err != nil {
		fmt.Println("cant fetch last synced height:", err)
	} else {
		fmt.Println("Last height:", height.Height)
	}

	heightStatuses, err := db.Heights.StatusCounts()
	if err == nil {
		for _, s := range heightStatuses {
			fmt.Printf("Status: %s, Count: %d\n", s.Status, s.Num)
		}
	} else {
		fmt.Println("cant fetch sync counts:", err)
	}

	client := coda.NewClient(http.DefaultClient, cfg.CodaEndpoint)
	status, err := client.GetDaemonStatus()
	if err != nil {
		return err
	}

	fmt.Println("=== Node Status ===")
	fmt.Println("Sync status:", status.SyncStatus)
	fmt.Println("Uptime:", status.UptimeSecs)
	fmt.Println("Version:", status.CommitID)

	if status.SyncStatus != coda.SyncStatusBootstrap {
		if status.StateHash != nil {
			fmt.Println("State hash:", *status.StateHash)
		}
		if status.BlockchainLength != nil {
			fmt.Println("Blockchain Length:", *status.BlockchainLength)
		}
		fmt.Println("Max length received:", status.HighestBlockLengthReceived)
		if status.NumAccounts != nil {
			fmt.Println("Accounts:", *status.NumAccounts)
		}
	} else {
		fmt.Println("Node details excluded because of bootstrap status")
	}

	return nil
}
