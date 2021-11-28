package api

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/site-inventory.json
var inventoryData []byte

//go:embed testdata/site-inventory-from-documentation.json
var documentationInventoryData []byte

func TestInventoryParse(t *testing.T) {
	var err error

	var inventory Inventory
	err = json.Unmarshal(inventoryData, &inventory)
	assert.NoError(t, err, "unable to parse actual inventory as Inventory")

	var documentationInventory Inventory
	err = json.Unmarshal(documentationInventoryData, &documentationInventory)
	assert.NoError(t, err, "unable to parse documentation-based inventory as Inventory")

}
