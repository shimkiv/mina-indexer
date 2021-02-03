package cli

import (
	"context"
	"fmt"
	"net/http"

	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/config"
)

func startStatus(cfg *config.Config) error {
	db, err := initStore(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	client := graph.NewClient(http.DefaultClient, cfg.CodaEndpoint)
	status, err := client.GetDaemonStatus(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("=== Node Status ===")
	fmt.Println("Sync status:", status.SyncStatus)
	fmt.Println("Uptime:", status.UptimeSecs)
	fmt.Println("Version:", status.CommitID)

	if status.SyncStatus != graph.SyncStatusBootstrap {
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
