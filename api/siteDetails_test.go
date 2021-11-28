package api

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed testdata/site-siteid-details.json
var siteDetailsData []byte

func TestSiteDetailsParse(t *testing.T) {
	var err error
	var siteDetailsDocument SiteDetailsDocument
	err = json.Unmarshal(siteDetailsData, &siteDetailsDocument)
	if assert.NoError(t, err, "unable to parse site details") {
		siteDetails := siteDetailsDocument.Details
		if assert.NotNil(t, siteDetails.LastUpdateTime, "lastUpdateTime should not be nil") {
			assert.False(t, siteDetails.LastUpdateTime.IsZero(), "lastUpdateTime should not be zero")
		}
		if assert.NotNil(t, siteDetails.InstallationDate, "installationDate should not be nil") {
			assert.False(t, siteDetails.InstallationDate.IsZero(), "installationDate should not be zero")
		}
	}

}
