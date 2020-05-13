package server

import (
	"errors"

	"github.com/gin-gonic/gin"

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
	s.GET("/blocks", s.GetBlocks)
	s.GET("/blocks/:hash", s.GetBlock)
	s.GET("/block", s.GetBlockByHeight)
	s.GET("/block_times", s.GetBlockTimes)
	s.GET("/block_times_interval", s.GetBlockTimesInterval)
	s.GET("/transactions", s.GetTransactions)
	s.GET("/transactions/:hash", s.GetTransaction)
	s.GET("/accounts", s.GetAccounts)

	return s
}

// Health reports the system health
func (s *Server) Health(c *gin.Context) {
	if err := s.db.Test(); err != nil {
		c.String(400, "ERROR")
		return
	}
	c.String(200, "OK")
}

// GetCurrentHeight returns the current blockchain height
func (s *Server) GetCurrentHeight(c *gin.Context) {
	height, err := s.db.Syncables.GetMostRecentHeight()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"height": height})
}

// GetBlock returns a single block
func (s *Server) GetBlock(c *gin.Context) {
	block, err := s.db.Blocks.FindByHash(c.Param("hash"))
	if err != nil {
		if err == store.ErrNotFound {
			notFound(c, err)
			return
		}
		badRequest(c, err)
		return
	}

	jsonOk(c, block)
}

//GetBlockByHeight returns a block for a given height
func (s *Server) GetBlockByHeight(c *gin.Context) {
	params := blockByHeightParams{}
	if err := c.BindQuery(&params); err != nil {
		badRequest(c, err)
		return
	}
	if params.Height <= 0 {
		badRequest(c, "height must greater than 0")
		return
	}

	block, err := s.db.Blocks.FindByHeight(params.Height)
	if err != nil {
		if err == store.ErrNotFound {
			notFound(c, err)
			return
		}
		badRequest(c, err)
		return
	}

	transactions, err := s.db.Transactions.ListByHeight(params.Height)
	if err != nil {
		badRequest(c, err)
		return
	}

	jsonOk(c, gin.H{
		"block":        block,
		"transactions": transactions,
	})
}

// GetBlocks returns a list of available blocks matching the filter
func (s *Server) GetBlocks(c *gin.Context) {
	params := store.BlockIndexParams{}

	if err := c.BindQuery(&params); err != nil {
		badRequest(c, err)
		return
	}

	blocks, err := s.db.Blocks.Index(params)
	if err != nil {
		badRequest(c, err)
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
	if err != nil {
		badRequest(c, err)
		return
	}

	jsonOk(c, result)
}

// GetTransaction returns a single transaction details
func (s *Server) GetTransaction(c *gin.Context) {
	tran, err := s.db.Transactions.FindByHash(c.Param("hash"))
	if err != nil {
		if err == store.ErrNotFound {
			notFound(c, err)
			return
		}
		badRequest(c, err)
		return
	}
	jsonOk(c, tran)
}

// GetTransactions returns transactions by height
func (s *Server) GetTransactions(c *gin.Context) {
	params := transactionsIndexParams{}
	if err := c.BindQuery(&params); err != nil {
		badRequest(c, err)
		return
	}

	transactions, err := s.db.Transactions.ListByHeight(params.Height)
	if err != nil {
		badRequest(c, err)
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

	accounts, err := s.db.Accounts.ListByHeight(params.Height)
	if err != nil {
		badRequest(c, err)
	}
	jsonOk(c, accounts)
}
