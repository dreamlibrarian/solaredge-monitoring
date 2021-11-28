package client

import (
	"fmt"
	"time"

	"github.com/dreamlibrarian/solaredge-monitoring/api"
)

const energyUsageEndpointTemplate = "site/%s/energy"

func (c *Client) GetEnergyUsage(siteID, timeUnit string, startTime, endTime time.Time) (*api.Energy, error) {
	result := &api.Energy{}

	req := c.CreateRequest(fmt.Sprintf(energyUsageEndpointTemplate, siteID))
	req.SetTimeParams(timeUnit, startTime, endTime)

	resp, err := c.do(c.client.Get, req)
	if err != nil {
		return nil, fmt.Errorf("unable to list sites: %w", err)
	}

	return result, handleResponse(resp, result)
}
