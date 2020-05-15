package server

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/figment-networks/coda-indexer/model"
	"github.com/figment-networks/coda-indexer/store"
)

// Server handles HTTP requests
type Server struct {
	*gin.Engine
	db *store.Store
}

// New returns a new server instance
func New(db *store.Store) *Server {
	s := &Server{
		db:     db,
		Engine: gin.Default(),
	}

	s.GET("/health", s.Health)
	s.GET("/height", s.GetCurrentHeight)
	s.GET("/block", s.GetCurrentBlock)
	s.GET("/blocks", s.GetBlocks)
	s.GET("/blocks/:id", s.GetBlock)
	s.GET("/block_times", s.GetBlockTimes)
	s.GET("/block_times_interval", s.GetBlockTimesInterval)
	s.GET("/transactions", s.GetTransactions)
	s.GET("/transactions/:id", s.GetTransaction)
	s.GET("/accounts", s.GetAccounts)
	s.GET("/accounts/:id", s.GetAccount)

	return s
}

// Health reports the system health
func (s *Server) Health(c *gin.Context) {
	if err := s.db.Test(); err != nil {
		c.String(500, "ERROR")
		return
	}
	c.String(200, "OK")
}

// GetCurrentHeight returns the current blockchain height
func (s *Server) GetCurrentHeight(c *gin.Context) {
	block, err := s.db.Blocks.Recent()
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, gin.H{
		"height": block.Height,
		"time":   block.Time,
	})
}

// GetCurrentBlock returns the current blockchain height
func (s *Server) GetCurrentBlock(c *gin.Context) {
	block, err := s.db.Blocks.Recent()
	if shouldReturn(c, err) {
		return
	}
	jsonOk(c, block)
}

// GetBlock returns a single block
func (s *Server) GetBlock(c *gin.Context) {
	var block *model.Block
	var err error

	id := resourceID(c, "id")
	if id.IsNumeric() {
		if id.Number() <= 0 {
			badRequest(c, errors.New("height must be greater than 0"))
			return
		}
		block, err = s.db.Blocks.FindByHeight(id.Number())
	} else {
		block, err = s.db.Blocks.FindByHash(id.String())
	}
	if shouldReturn(c, err) {
		return
	}

	creator, err := s.db.Accounts.FindByPublicKey(block.Creator)
	if shouldReturn(c, err) {
		return
	}

	transactions, err := s.db.Transactions.ByHeight(block.Height)
	if shouldReturn(c, err) {
		return
	}

	jobs, err := s.db.Jobs.ByHeight(block.Height)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, gin.H{
		"block":        block,
		"creator":      creator,
		"transactions": transactions,
		"jobs":         jobs,
	})
}

// GetBlocks returns a list of available blocks matching the filter
func (s *Server) GetBlocks(c *gin.Context) {
	search := store.BlockSearch{}

	if err := c.BindQuery(&search); err != nil {
		badRequest(c, err)
		return
	}

	blocks, err := s.db.Blocks.Search(search)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, blocks)
}

// GetBlockTimes returns avg block times info
func (s *Server) GetBlockTimes(c *gin.Context) {
	params := blockTimesParams{}

	if err := c.BindQuery(&params); err != nil {
		badRequest(c, err)
		return
	}
	params.setDefaults()

	result, err := s.db.Blocks.AvgRecentTimes(params.Limit)
	if err != nil {
		badRequest(c, err)
		return
	}

	jsonOk(c, result)
}

// GetBlockTimesInterval returns block stats for an interval
func (s *Server) GetBlockTimesInterval(c *gin.Context) {
	params := blockTimesIntervalParams{}

	if err := c.BindQuery(&params); err != nil {
		badRequest(c, err)
		return
	}
	params.setDefaults()

	result, err := s.db.Blocks.AvgTimesForInterval(params.Interval, params.Period)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, result)
}

// GetTransaction returns a single transaction details
func (s *Server) GetTransaction(c *gin.Context) {
	var tran *model.Transaction
	var err error

	id := resourceID(c, "id")
	if id.IsNumeric() {
		tran, err = s.db.Transactions.FindByID(id.Number())
	} else {
		tran, err = s.db.Transactions.FindByHash(id.String())
	}
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, tran)
}

// GetTransactions returns transactions by height
func (s *Server) GetTransactions(c *gin.Context) {
	search := store.TransactionSearch{}
	if err := c.BindQuery(&search); err != nil {
		badRequest(c, err)
		return
	}

	transactions, err := s.db.Transactions.Search(search)
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, transactions)
}

// GetAccounts returns all accounts
func (s *Server) GetAccounts(c *gin.Context) {
	params := accountsIndexParams{}
	if err := c.BindQuery(&params); err != nil {
		badRequest(c, err)
		return
	}

	if params.Height <= 0 {
		badRequest(c, errors.New("height is required"))
		return
	}

	accounts, err := s.db.Accounts.ByHeight(params.Height)
	if err != nil {
		badRequest(c, err)
	}

	jsonOk(c, accounts)
}

// GetAccount returns account for by hash or ID
func (s *Server) GetAccount(c *gin.Context) {
	var acc *model.Account
	var err error

	id := resourceID(c, "id")
	if id.IsNumeric() {
		acc, err = s.db.Accounts.FindByID(id.Number())
	} else {
		acc, err = s.db.Accounts.FindByPublicKey(id.String())
	}
	if shouldReturn(c, err) {
		return
	}

	jsonOk(c, acc)
}
