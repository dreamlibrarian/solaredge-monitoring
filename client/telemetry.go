package client

import (
	"fmt"
	"time"

	"github.com/dreamlibrarian/solaredge-monitoring/api"
)

// Takes siteID and equipment SN.
const equipmentDataEndpointTemplate = "site/%s/%s/data"

func (c *Client) GetTelemetryForEquipment(siteID, serialNumber string, timeUnit string, startTime, endTime time.Time) ([]api.Telemetry, error) {
	var edd api.EquipmentDataDocument

	req := c.CreateRequestf(equipmentDataEndpointTemplate, siteID, serialNumber)
	req.SetTimeParams(timeUnit, startTime, endTime)

	resp, err := c.do(c.client.Get, req)
	if err != nil {
		return nil, fmt.Errorf("unable to get telemetry: %w", err)
	}

	return edd.Data.Telemetries, handleResponse(resp, &edd)
}

func (c *Client) GetTelemetryForAllInverters(siteID string, timeUnit string, startTime, endTime time.Time) (map[string][]api.Telemetry, error) {
	result := make(map[string][]api.Telemetry)

	inventory, err := c.GetSiteInventory(siteID)
	if err != nil {
		return nil, fmt.Errorf("unable to inventory site: %w", err)
	}

	for _, i := range inventory.Inverters {
		sn := i.SerialNumber
		telemetries, err := c.GetTelemetryForEquipment(siteID, sn, timeUnit, startTime, endTime)
		if err != nil {
			return nil, fmt.Errorf("error while querying site %s serial %s: %w", siteID, sn, err)
		}
		result[sn] = telemetries
	}

	return result, nil
}
