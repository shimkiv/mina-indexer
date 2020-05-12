package cli

import (
	"fmt"
	"net/http"

	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/config"
)

func startStatus(cfg *config.Config) error {
	client := coda.NewClient(http.DefaultClient, cfg.CodaEndpoint)

	status, err := client.GetDaemonStatus()
	if err != nil {
		terminate(err)
	}

	fmt.Println("Sync status:", status.SyncStatus)
	fmt.Println("Uptime:", status.UptimeSecs)
	fmt.Println("Version:", status.CommitID)
	fmt.Println("Blockchain Length:", *status.BlockchainLength)
	fmt.Println("Max length received:", status.HighestBlockLengthReceived)
	fmt.Println("State hash:", *status.StateHash)
	fmt.Println("Accounts:", *status.NumAccounts)

	return nil
}
