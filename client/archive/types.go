package archive

const (
	LedgerTypeCurrent = "current"
	LedgerTypeStaged  = "staged"
)

type BlocksRequest struct {
	Canonical   *bool
	StartHeight uint
	Limit       uint
}

type Summary struct {
	BlocksCount           uint            `json:"blocks_count"`
	BlocksMinHeight       uint            `json:"blocks_min_height"`
	BlocksMaxHeight       uint            `json:"blocks_max_height"`
	BlocksMinTimestamp    int64           `json:"blocks_min_timestamp"`
	BlocksMaxTimestamp    int64           `json:"blocks_max_timestamp"`
	BlocksProducersCount  uint            `json:"blocks_producers_count"`
	PublicKeysCount       uint            `json:"public_keys_count"`
	InternalCommandsCount uint            `json:"internal_commands_count"`
	UserCommandsCount     uint            `json:"user_commands_count"`
	UserCommandsTypes     map[string]uint `json:"user_commands_types"`
	InternalCommandsTypes map[string]uint `json:"internal_commands_types"`
}

type Block struct {
	Height                 uint64            `json:"height"`
	StateHash              string            `json:"state_hash"`
	ParentHash             string            `json:"parent_hash"`
	LedgerHash             string            `json:"ledger_hash"`
	SnarkedLedgerHash      string            `json:"snarked_ledger_hash"`
	Creator                string            `json:"creator"`
	Winner                 string            `json:"winner"`
	Timestamp              int64             `json:"timestamp"`
	TimestampFormatted     string            `json:"timestamp_formatted"`
	GlobalSlotSinceGenesis uint              `json:"global_slot_since_genesis"`
	GlobalSlot             uint              `json:"global_slot"`
	InternalCommands       []InternalCommand `json:"internal_commands"`
	UserCommands           []UserCommand     `json:"user_commands"`
}

type InternalCommand struct {
	ID                  string `json:"id"`
	Hash                string `json:"hash"`
	Type                string `json:"type"`
	Fee                 int64  `json:"fee"`
	Token               int    `json:"token"`
	Receiver            string `json:"receiver"`
	SequenceNo          int    `json:"sequence_no"`
	SecondarySequenceNo int    `json:"secondary_sequence_no"`
}

type UserCommand struct {
	Hash                           string  `json:"hash"`
	Type                           string  `json:"type"`
	FeeToken                       int     `json:"fee_token"`
	Token                          int     `json:"token"`
	Nonce                          int     `json:"nonce"`
	Amount                         int64   `json:"amount"`
	Fee                            int64   `json:"fee"`
	ValidUntil                     *uint64 `json:"valid_until"`
	Memo                           string  `json:"memo"`
	Status                         string  `json:"status"`
	FailureReason                  *string `json:"failure_reason"`
	FeePayerAccountCreationFeePaid *uint64 `json:"fee_payer_account_creation_fee_paid"`
	ReceiverAccountCreationFeePaid *uint64 `json:"receiver_account_creation_fee_paid"`
	CreatedToken                   *uint64 `json:"created_token"`
	SequenceNo                     int     `json:"sequence_no"`
	FeePayer                       string  `json:"fee_payer"`
	Sender                         string  `json:"sender"`
	Receiver                       string  `json:"receiver"`
}

type StakingInfo struct {
	Pk       string `json:"pk"`
	Balance  string `json:"balance"`
	Delegate string `json:"delegate"`
	Timing   *struct {
		InitialMinimumBalance string `json:"initial_minimum_balance"`
		CliffTime             string `json:"cliff_time"`
		CliffAmount           string `json:"cliff_amount"`
		VestingPeriod         string `json:"vesting_period"`
		VestingIncrement      string `json:"vesting_increment"`
	} `json:"timing"`
	Token            string `json:"token"`
	TokenPermissions struct {
	} `json:"token_permissions"`
	ReceiptChainHash string `json:"receipt_chain_hash"`
	VotingFor        string `json:"voting_for"`
	Permissions      *struct {
		Stake              bool   `json:"stake"`
		EditState          string `json:"edit_state"`
		Send               string `json:"send"`
		SetDelegate        string `json:"set_delegate"`
		SetPermissions     string `json:"set_permissions"`
		SetVerificationKey string `json:"set_verification_key"`
	} `json:"permissions"`
}
