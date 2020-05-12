package coda

import "fmt"

var (
	// Get the node status
	queryDaemonStatus = `
		query {
  		daemonStatus {
				syncStatus
				uptimeSecs
				stateHash
				commitId
				peers
				highestBlockLengthReceived
				blockchainLength
				numAccounts
			}
		}
	`

	// Get block details
	queryBlocks = `
		query {
			blocks(%s) {
				nodes {
					%s
				}
			}
		}
	`

	// Block details fields
	queryBlockFields = `
		stateHash
		stateHashField
		creator
		creatorAccount {
			delegate
      publicKey
			balance {
      	blockHeight
        total
        unknown
      }
		}
		protocolState {
			blockchainState {
				date
				utcDate
				stagedLedgerHash
				snarkedLedgerHash
			}
			consensusState {
        blockHeight
        blockchainLength
        epoch
        epochCount
        hasAncestorInSameCheckpointWindow
        lastVrfOutput
				totalCurrency
				minWindowDensity
        slot
        stakingEpochData {
          epochLength
          lockCheckpoint
          seed
          startCheckpoint
        }
      }
			previousStateHash
		}
		snarkJobs {
			fee
			prover
			workIds
		}
		transactions {
			coinbase
			feeTransfer {
				recipient
				fee
			}
			userCommands {
				amount
				fee
				from
				id
				isDelegation
				memo
				nonce
				to
				fromAccount {
					delegate
      		publicKey
					balance {
      			blockHeight
        		total
        		unknown
      		}
				}
				toAccount {
					delegate
      		publicKey
					balance {
      			blockHeight
        		total
        		unknown
      		}
				}
			}
		}
	`

	queryAccount = `
		query {
			account(publicKey: "%s") {
    		nonce
    		inferredNonce
    		receiptChainHash
    		delegate
    		votingFor
    		locked
    		balance {
      		unknown
      		total
      		blockHeight
				}
			}
  }	
`
)

func buildBlocksQuery(filter string) string {
	return fmt.Sprintf(queryBlocks, filter, queryBlockFields)
}

func buildAccountQuery(filter string) string {
	return fmt.Sprintf(queryAccount, filter)
}
