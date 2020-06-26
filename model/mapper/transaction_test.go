package mapper

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/figment-networks/coda-indexer/coda"
	"github.com/figment-networks/coda-indexer/model"
	"github.com/stretchr/testify/assert"
)

func TestUserTransaction(t *testing.T) {
}

func TestBlockRewardTransaction(t *testing.T) {
}

func TestTransactions(t *testing.T) {
	block := loadTestBlock("../../test/fixtures/block.json")
	transactions, err := Transactions(block)

	assert.NoError(t, err)
	assert.NotNil(t, transactions)
	assert.Equal(t, 11, len(transactions))

	tx := transactions[0]
	assert.Equal(t, model.TxTypeBlockReward, tx.Type)
	assert.Nil(t, tx.Sender)
	assert.Equal(t, block.Transactions.CoinbaseReceiver.PublicKey, tx.Receiver)
	assert.Equal(t, uint64(200000000000), tx.Amount)
}

func loadTestBlock(path string) *coda.Block {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var block coda.Block
	if err := json.Unmarshal(data, &block); err != nil {
		panic(err)
	}

	return &block
}
