package staketab

import (
	"encoding/json"
	"net/http"
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

// GetProviders returns providers
func (c Client) GetProviders() (Providers, error) {
	resp, err := c.client.Get(c.endpoint + "/get_providers")
	if err != nil {
		return Providers{}, err
	}

	providers := &Providers{}
	err = json.NewDecoder(resp.Body).Decode(providers)

	return *providers, err
}
