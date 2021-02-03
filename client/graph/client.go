package graph

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	ErrBlockNotFound = errors.New("block not found")
	ErrBlockInvalid  = errors.New("block is invalid")
)

// Client is a GraphQL API client
type Client struct {
	endpoint string
	client   *http.Client
	debug    bool
}

// NewClient returns a new client for a given endpoint
func NewClient(client *http.Client, endpoint string) *Client {
	return &Client{
		client:   client,
		endpoint: endpoint,
	}
}

// NewDefaultClient returns a default client for a given endpoint
func NewDefaultClient(endpoint string) *Client {
	return &Client{
		client: &http.Client{
			Timeout: time.Minute * 5,
		},
		endpoint: endpoint,
	}
}

func (c *Client) SetDebug(enabled bool) {
	c.debug = enabled
}

// Execute make a GraphQL query and returns the response
func (c Client) Execute(ctx context.Context, q string) (*GraphResponse, error) {
	r := map[string]string{"query": q}
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, err
	}
	reqBody := bytes.NewReader(data)

	if c.debug {
		fmt.Printf("%s\n", q)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if c.debug {
		log.Debugf("client response: %s\n", respBody)
	}

	graphResp := GraphResponse{}
	if err := json.Unmarshal(respBody, &graphResp); err != nil {
		switch err.(type) {
		case *json.UnmarshalTypeError, *json.UnmarshalFieldError:
			return nil, errors.New(string(respBody))
		default:
			return nil, err
		}
	}

	if len(graphResp.Errors) > 0 {
		return nil, errors.New(graphResp.Errors[0].Message)
	}

	return &graphResp, nil
}

// Query executes the query and parses the result
func (c Client) Query(input string, out interface{}) error {
	resp, err := c.Execute(context.Background(), input)
	if err != nil {
		return err
	}
	return resp.Decode(out)
}

// QueryWithContext executes the query with context and parses the result
func (c Client) QueryWithContext(ctx context.Context, input string, out interface{}) error {
	resp, err := c.Execute(ctx, input)
	if err != nil {
		return err
	}
	return resp.Decode(out)
}

// GetDaemonStatus returns current node daemon status
func (c Client) GetDaemonStatus(ctx context.Context) (*DaemonStatus, error) {
	var result struct {
		DaemonStatus `json:"daemonStatus"`
	}
	if err := c.QueryWithContext(ctx, queryDaemonStatus, &result); err != nil {
		return nil, err
	}
	return &result.DaemonStatus, nil
}

// GetCurrentHeight returns the current blockchain height
func (c Client) GetCurrentHeight() (int64, error) {
	block, err := c.GetLastBlock()
	if err != nil {
		return 0, err
	}
	if block == nil {
		return 0, ErrBlockNotFound
	}
	if block.ProtocolState == nil {
		return 0, ErrBlockInvalid
	}
	if block.ProtocolState.ConsensusState == nil {
		return 0, ErrBlockInvalid
	}

	height := block.ProtocolState.ConsensusState.BlockHeight
	return strconv.ParseInt(height, 10, 64)
}

// GetBestChain returns the blocks from the canonical chain
func (c Client) GetBestChain() ([]Block, error) {
	var result struct {
		Blocks []Block `json:"bestChain"`
	}
	q := buildBestChainQuery()
	if err := c.Query(q, &result); err != nil {
		return nil, err
	}
	return result.Blocks, nil
}

// GetBlock returns a single block for the given state hash
func (c Client) GetBlock(hash string) (*Block, error) {
	q := fmt.Sprintf(queryBlock, hash, queryBlockFields)
	result := struct {
		Block Block `json:"block"`
	}{}

	if err := c.Query(q, &result); err != nil {
		return nil, err
	}

	return &result.Block, nil
}

// GetBlocks returns blocks for a filter
func (c Client) GetBlocks(filter string) ([]Block, error) {
	var result struct {
		Blocks struct {
			Nodes []Block `json:"nodes"`
		} `json:"blocks"`
	}

	q := buildBlocksQuery(filter)
	if err := c.Query(q, &result); err != nil {
		return nil, err
	}

	return result.Blocks.Nodes, nil
}

// GetSingleBlock returns a single block record from the result
func (c Client) GetSingleBlock(filter string) (*Block, error) {
	blocks, err := c.GetBlocks(filter)
	if err != nil {
		return nil, err
	}
	if len(blocks) == 0 {
		return nil, nil
	}
	return &blocks[0], nil
}

// GetFirstBlock returns the first block available in the chain node
func (c Client) GetFirstBlock() (*Block, error) {
	return c.GetSingleBlock("first:1")
}

// GetFirstBlocks returns the first n blocks
func (c Client) GetFirstBlocks(n int) ([]Block, error) {
	filter := fmt.Sprintf("first:%v", n)
	return c.GetBlocks(filter)
}

// GetLastBlock returns the last block available in the chain node
func (c Client) GetLastBlock() (*Block, error) {
	return c.GetSingleBlock("last:1")
}

// GetNextBlock returns the next block after the given block's hash
func (c Client) GetNextBlock(after string) (*Block, error) {
	if after == "" {
		return c.GetFirstBlock()
	}

	filter := fmt.Sprintf("after:%q,first:1", after)
	return c.GetSingleBlock(filter)
}

// GetNextBlocks returns a next N blocks after a given block hash
func (c Client) GetNextBlocks(after string, n int) ([]Block, error) {
	return c.GetBlocks(fmt.Sprintf("after:%q,first:%v", after, n))
}

// GetAccount returns account for a given public key
func (c Client) GetAccount(publicKey string) (*Account, error) {
	var result struct {
		Account Account `json:"account"`
	}
	if err := c.Query(buildAccountQuery(publicKey), &result); err != nil {
		return nil, err
	}
	return &result.Account, nil
}
