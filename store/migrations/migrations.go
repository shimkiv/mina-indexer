package migrations

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assets56014d3cd9fc1bf04559f8253464e04746c520f8 = "-- +goose Up\nCREATE TYPE CHAIN_TX_TYPE AS ENUM (\n  'payment',\n  'delegation',\n  'coinbase',\n  'fee_transfer_via_coinbase',\n  'fee_transfer',\n  'snark_fee'\n);\n\nCREATE TYPE CHAIN_TX_STATUS AS ENUM ('applied', 'failed');\n\nCREATE TABLE transactions (\n  id                        SERIAL NOT NULL,\n  type                      CHAIN_TX_TYPE,\n  hash                      TEXT NOT NULL,\n  block_hash                TEXT NOT NULL,\n  block_height              CHAIN_HEIGHT,\n  time                      CHAIN_TIME,\n  nonce                     INTEGER,\n  sender                    TEXT,\n  receiver                  TEXT NOT NULL,\n  amount                    CHAIN_CURRENCY,\n  fee                       CHAIN_CURRENCY,\n  memo                      TEXT,\n  status                    CHAIN_TX_STATUS,\n  failure_reason            TEXT,\n  sequence_number           INTEGER,\n  secondary_sequence_number INTEGER,\n  created_at                CHAIN_TIME,\n  updated_at                CHAIN_TIME\n);\n\nCREATE INDEX idx_transactions_type\n  ON transactions(type);\n\nCREATE INDEX idx_transactions_block_hash\n  ON transactions(block_hash);\n\nCREATE UNIQUE INDEX idx_transactions_hash\n  ON transactions(hash);\n\nCREATE INDEX idx_transactions_block_height\n  ON transactions(block_height);\n\nCREATE INDEX idx_transactions_time\n  ON transactions(time);\n\nCREATE INDEX idx_transactions_sender\n  ON transactions(sender);\n\nCREATE INDEX idx_transactions_receiver\n  ON transactions(receiver);\n\nCREATE INDEX idx_transactions_memo\n  ON transactions(LOWER(memo));\n\n-- +goose Down\nDROP TABLE transactions;\nDROP TYPE CHAIN_TX_TYPE;\nDROP TYPE CHAIN_TX_STATUS;\n"
var _Assetseb0dad53ad772fd424cf4a289a88fdc15078e06c = ""
var _Assetse412840a468eca59245b861d26ad081f0facbe0f = "-- +goose Up\nCREATE TABLE blocks (\n  id                   SERIAL NOT NULL,\n  canonical            BOOLEAN NOT NULL,\n  height               CHAIN_HEIGHT,\n  time                 CHAIN_TIME,\n  hash                 TEXT NOT NULL,\n  parent_hash          TEXT NOT NULL,\n  ledger_hash          TEXT NOT NULL,\n  snarked_ledger_hash  TEXT NOT NULL,\n  creator              TEXT NOT NULL,\n  coinbase             CHAIN_CURRENCY DEFAULT 0,\n  total_currency       CHAIN_CURRENCY DEFAULT 0,\n  epoch                INTEGER DEFAULT 0,\n  slot                 INTEGER DEFAULT 0,\n  transactions_count   INTEGER NOT NULL DEFAULT 0,\n  transactions_fees    CHAIN_CURRENCY DEFAULT 0,\n  snarkers_count       INTEGER NOT NULL DEFAULT 0,\n  snark_jobs_count     INTEGER NOT NULL DEFAULT 0,\n  snark_jobs_fees      CHAIN_CURRENCY,\n\n  PRIMARY KEY (id)\n);\n\nCREATE INDEX idx_blocks_canonic\n  ON blocks (canonical);\n\nCREATE INDEX idx_blocks_height\n  ON blocks (height);\n\nCREATE INDEX idx_blocks_time\n  ON blocks (time);\n\nCREATE UNIQUE INDEX idx_blocks_hash\n  ON blocks (hash);\n\nCREATE INDEX idx_blocks_creator\n  ON blocks (creator);\n\n-- +goose Down\nDROP TABLE IF EXISTS blocks;"
var _Assets95bb338fdbc7ab8063566a16b82143a8a1a14866 = "-- +goose Up\nCREATE TABLE validator_stats (\n  id                    CHAIN_UUID,\n  validator_id          INTEGER NOT NULL,\n  time                  CHAIN_TIME,\n  bucket                CHAIN_INTERVAL,\n  blocks_produced_count INTEGER DEFAULT 0,\n  delegations_count     INTEGER DEFAULT 0,\n  delegations_amount    CHAIN_CURRENCY\n);\n\nCREATE UNIQUE INDEX idx_validator_stats_bucket\n  ON validator_stats(time, bucket, validator_id);\n\n-- +goose Down\nDROP TABLE validator_stats;\n"
var _Assetsde9d94875307644d384ac39183586da629e79511 = "-- +goose Up\nCREATE TABLE IF NOT EXISTS validators (\n  id                 SERIAL NOT NULL,\n  public_key         TEXT NOT NULL,\n  identity_name      TEXT,\n  start_height       CHAIN_HEIGHT,\n  start_time         CHAIN_TIME,\n  last_height        CHAIN_HEIGHT,\n  last_time          CHAIN_TIME,\n  stake              CHAIN_CURRENCY,\n  blocks_created     INTEGER NOT NULL DEFAULT 0,\n  blocks_proposed    INTEGER NOT NULL DEFAULT 0,\n  delegated_accounts INTEGER NOT NULL DEFAULT 0,\n  delegated_balance  CHAIN_CURRENCY,\n  created_at         CHAIN_TIME,\n  updated_at         CHAIN_TIME,\n\n  PRIMARY KEY (id)\n);\n\nCREATE UNIQUE INDEX idx_validators_publix_key\n  ON validators (public_key);\n\n-- +goose Down\nDROP TABLE validators;\n"
var _Assets8f0ce2f6ab2de2fde8e39e373c1169f7098589dc = "-- +goose Up\nCREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";\n\nCREATE DOMAIN CHAIN_UUID     UUID NOT NULL DEFAULT uuid_generate_v4();\nCREATE DOMAIN CHAIN_CURRENCY DECIMAL(65, 0);\nCREATE DOMAIN CHAIN_TIME     TIMESTAMP WITH TIME ZONE NOT NULL;\nCREATE DOMAIN CHAIN_HEIGHT   INTEGER NOT NULL;\n\nCREATE TYPE CHAIN_INTERVAL AS ENUM ('h', 'd', 'w');\n\n-- +goose Down\nDROP EXTENSION IF EXISTS \"uuid-ossp\";\n\nDROP DOMAIN IF EXISTS CHAIN_UUID;\nDROP DOMAIN IF EXISTS CHAIN_CURRENCY;\nDROP DOMAIN IF EXISTS CHAIN_TIME;\nDROP DOMAIN IF EXISTS CHAIN_HEIGHT;\nDROP TYPE   IF EXISTS CHAIN_INTERVAL;"
var _Assets0e21e93d4c07a5a5a7f0de2a3a90dfdf60ddae84 = "-- +goose Up\nCREATE TABLE IF NOT EXISTS snarkers (\n  id           SERIAL NOT NULL,\n  account      TEXT NOT NULL,\n  start_height INTEGER NOT NULL,\n  start_time   TIMESTAMP WITH TIME ZONE NOT NULL,\n  last_height  INTEGER NOT NULL,\n  last_time    TIMESTAMP WITH TIME ZONE NOT NULL,\n  created_at   TIMESTAMP WITH TIME ZONE NOT NULL,\n  updated_at   TIMESTAMP WITH TIME ZONE NOT NULL,\n  fee          DECIMAL(65, 0) NOT NULL,\n  jobs_count   INTEGER DEFAULT 0,\n  works_count  INTEGER DEFAULT 0,\n\n  PRIMARY KEY (id)\n);\n\nCREATE UNIQUE INDEX idx_snarkers_account\n  ON snarkers(account);\n\n-- +goose Down\nDROP TABLE snarkers;\n"
var _Assets26cd7018beecca6820637c4b60d1b171cdd96ce5 = "-- +goose Up\nCREATE TABLE chain_stats (\n  time                CHAIN_TIME,\n  bucket              CHAIN_INTERVAL,\n\n  block_time_avg      NUMERIC,\n  blocks_count        INTEGER,\n  blocks_total_count  INTEGER,\n  transactions_count  INTEGER,\n  validators_count    INTEGER,\n  accounts_count      INTEGER,\n  epochs_count        INTEGER,\n  slots_count         INTEGER,\n  snarkers_count      INTEGER,\n  jobs_count          INTEGER,\n  coinbase            CHAIN_CURRENCY,\n  total_currency      CHAIN_CURRENCY,\n\n  PRIMARY KEY (time, bucket)\n);\n\n-- +goose Down\nDROP TABLE chain_stats;\n"
var _Assets4ba8ddba8971be822efa4ce3c8122eba05441ef3 = "-- +goose Up\nCREATE TABLE IF NOT EXISTS accounts (\n  id              SERIAL NOT NULL,\n  public_key      TEXT NOT NULL,\n  delegate        TEXT,\n  balance         CHAIN_CURRENCY,\n  balance_unknown CHAIN_CURRENCY,\n  nonce           INTEGER NOT NULL,\n  stake           CHAIN_CURRENCY,\n  start_height    CHAIN_HEIGHT,\n  start_time      CHAIN_TIME,\n  last_height     CHAIN_HEIGHT,\n  last_time       CHAIN_TIME,\n  created_at      CHAIN_TIME,\n  updated_at      CHAIN_TIME,\n\n  PRIMARY KEY (id)\n);\n\nCREATE UNIQUE INDEX idx_accounts_public_key\n  ON accounts(public_key);\n\n-- +goose Down\nDROP TABLE accounts;\n"
var _Assets291657202d99299958f048acc27d88d19c2c1136 = "-- +goose Up\nCREATE TABLE snark_jobs (\n  id          SERIAL NOT NULL,\n  height      CHAIN_HEIGHT,\n  time        CHAIN_TIME,\n  prover      TEXT NOT NULL,\n  fee         CHAIN_CURRENCY,\n  works_count INTEGER NOT NULL,\n  created_at  CHAIN_TIME,\n\n  PRIMARY KEY (id)\n);\n\nCREATE INDEX idx_snark_jobs_time\n  ON snark_jobs(time);\n\nCREATE INDEX idx_snark_jobs_height\n  ON snark_jobs(height);\n\nCREATE INDEX idx_snark_jobs_prover\n  ON snark_jobs(prover);\n\n-- +goose Down\nDROP TABLE snark_jobs;\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"store"}, "/store": []string{"migrations"}, "/store/migrations": []string{"010_create_validator_stats.sql", "004_create_transactions.sql", "007_create_snarkers.sql", "008_create_chain_stats.sql", "migrations.go", "005_create_accounts.sql", "006_create_snark_jobs.sql", "003_create_validators.sql", "002_create_blocks.sql", "001_create_objects.sql"}}, map[string]*assets.File{
	"/store/migrations/010_create_validator_stats.sql": &assets.File{
		Path:     "/store/migrations/010_create_validator_stats.sql",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1612237124, 1612237124874817839),
		Data:     []byte(_Assets95bb338fdbc7ab8063566a16b82143a8a1a14866),
	}, "/store/migrations/004_create_transactions.sql": &assets.File{
		Path:     "/store/migrations/004_create_transactions.sql",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1612236908, 1612236908788380078),
		Data:     []byte(_Assets56014d3cd9fc1bf04559f8253464e04746c520f8),
	}, "/store/migrations/migrations.go": &assets.File{
		Path:     "/store/migrations/migrations.go",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1612310424, 1612310424639941141),
		Data:     []byte(_Assetseb0dad53ad772fd424cf4a289a88fdc15078e06c),
	}, "/store/migrations/002_create_blocks.sql": &assets.File{
		Path:     "/store/migrations/002_create_blocks.sql",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1612236894, 1612236894882806817),
		Data:     []byte(_Assetse412840a468eca59245b861d26ad081f0facbe0f),
	}, "/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1612310425, 1612310425961936789),
		Data:     nil,
	}, "/store": &assets.File{
		Path:     "/store",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1611718308, 1611718308766199326),
		Data:     nil,
	}, "/store/migrations": &assets.File{
		Path:     "/store/migrations",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1611969797, 1611969797625711424),
		Data:     nil,
	}, "/store/migrations/006_create_snark_jobs.sql": &assets.File{
		Path:     "/store/migrations/006_create_snark_jobs.sql",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1612236921, 1612236921926779637),
		Data:     []byte(_Assets291657202d99299958f048acc27d88d19c2c1136),
	}, "/store/migrations/003_create_validators.sql": &assets.File{
		Path:     "/store/migrations/003_create_validators.sql",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1612236902, 1612236902600243736),
		Data:     []byte(_Assetsde9d94875307644d384ac39183586da629e79511),
	}, "/store/migrations/001_create_objects.sql": &assets.File{
		Path:     "/store/migrations/001_create_objects.sql",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1611969959, 1611969959862037523),
		Data:     []byte(_Assets8f0ce2f6ab2de2fde8e39e373c1169f7098589dc),
	}, "/store/migrations/007_create_snarkers.sql": &assets.File{
		Path:     "/store/migrations/007_create_snarkers.sql",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1600216739, 1600216739000000000),
		Data:     []byte(_Assets0e21e93d4c07a5a5a7f0de2a3a90dfdf60ddae84),
	}, "/store/migrations/008_create_chain_stats.sql": &assets.File{
		Path:     "/store/migrations/008_create_chain_stats.sql",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1612309089, 1612309089341190339),
		Data:     []byte(_Assets26cd7018beecca6820637c4b60d1b171cdd96ce5),
	}, "/store/migrations/005_create_accounts.sql": &assets.File{
		Path:     "/store/migrations/005_create_accounts.sql",
		FileMode: 0x1a4,
		Mtime:    time.Unix(1612304611, 1612304611547918183),
		Data:     []byte(_Assets4ba8ddba8971be822efa4ce3c8122eba05441ef3),
	}}, "")
