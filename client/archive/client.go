package archive

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Client interacts with the Archive API
type Client struct {
	endpoint string
	client   *http.Client
}

// NewClient returns a new archive service client
func NewClient(httpClient *http.Client, endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		client:   httpClient,
	}
}

// NewDefaultClient returns a default archive service client
func NewDefaultClient(endpoint string) *Client {
	return NewClient(http.DefaultClient, endpoint)
}

// Summary returns archive summary
func (c Client) Summary() (*Summary, error) {
	resp, err := c.client.Get(c.endpoint + "/")
	if err != nil {
		return nil, err
	}

	summary := &Summary{}
	err = json.NewDecoder(resp.Body).Decode(summary)

	return summary, err
}

// Blocks returns blocks matching the request parameters
func (c Client) Blocks(blocksReq *BlocksRequest) ([]Block, error) {
	req, err := http.NewRequest(http.MethodGet, c.endpoint+"/blocks", nil)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("start_height", fmt.Sprintf("%v", blocksReq.StartHeight))
	if blocksReq.Canonical != nil {
		params.Add("canonical", fmt.Sprintf("%v", *blocksReq.Canonical))
	}
	params.Add("limit", fmt.Sprintf("%v", blocksReq.Limit))
	req.URL.RawQuery = params.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := []Block{}
	err = json.NewDecoder(resp.Body).Decode(&result)

	return result, err
}

// Block returns block for a given hash
func (c Client) Block(hash string) (*Block, error) {
	resp, err := c.client.Get(fmt.Sprintf("%s/blocks/%s", c.endpoint, hash))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	block := &Block{}
	err = json.NewDecoder(resp.Body).Decode(block)

	return block, err
}

// StakingLedger returns the staking ledger records
func (c Client) StakingLedger(ledgerType string) ([]StakingInfo, error) {
	path := fmt.Sprintf("%s/staking_ledger?type=%s", c.endpoint, ledgerType)

	resp, err := c.client.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := []StakingInfo{}
	err = json.NewDecoder(resp.Body).Decode(&result)

	return result, err
}
