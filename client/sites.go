package client

import (
	"fmt"

	"github.com/dreamlibrarian/solaredge-monitoring/api"
)

const siteListEndpoint = "/sites/list"

func (c *Client) GetSiteList() ([]api.SiteDetails, error) {
	var result api.SiteListDocument

	url := c.baseURL
	url.Path = siteListEndpoint

	response, err := c.client.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("unable to list sites: %w", err)
	}

	return result.Sites.Sites, handleResponse(response, &result)
}
