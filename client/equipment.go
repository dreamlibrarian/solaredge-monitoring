package client

import (
	"fmt"

	"github.com/dreamlibrarian/solaredge-monitoring/api"
)

const (
	// expects siteID
	siteInventoryEndpointTemplate = "site/%s/inventory"
)

func (c *Client) GetSiteInventory(siteID string) (*api.Inventory, error) {

	result := &api.InventoryDocument{}

	req := c.CreateRequestf(siteInventoryEndpointTemplate, siteID)

	resp, err := c.do(c.client.Get, req)
	if err != nil {
		return nil, fmt.Errorf("unable to get site inventory: %w", err)
	}

	return &result.Inventory, handleResponse(resp, result)

}
