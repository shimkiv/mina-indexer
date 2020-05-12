package coda

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	ErrBlockNotFound = errors.New("block not found")
	ErrBlockInvalid  = errors.New("block in not valid")
)

// Client is a GraphQL API client
type Client struct {
	endpoint string
	client   *http.Client
}

// NewClient returns a new client for a given endpoint
func NewClient(client *http.Client, endpoint string) *Client {
	return &Client{
		client:   client,
		endpoint: endpoint,
	}
}

// Execute make a GraphQL query and returns the response
func (c Client) Execute(q string) (*GraphResponse, error) {
	r := map[string]string{"query": q}
	data, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	reqBody := bytes.NewReader(data)

	req, err := http.NewRequest(http.MethodPost, c.endpoint, reqBody)
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
	resp, err := c.Execute(input)
	if err != nil {
		return err
	}
	return resp.Decode(out)
}

// GetDaemonStatus returns current node daemon status
func (c Client) GetDaemonStatus() (*DaemonStatus, error) {
	var result struct {
		DaemonStatus `json:"daemonStatus"`
	}
	if err := c.Query(queryDaemonStatus, &result); err != nil {
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
