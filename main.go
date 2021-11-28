package main

import (
	"fmt"
	"os"

	"github.com/dreamlibrarian/solaredge-monitoring/client"
)

func main() {

	apiKey := os.Getenv("APIKEY")

	c := client.NewClient(apiKey)

	list, err := c.GetSiteList()
	if err != nil {
		fmt.Printf("Couldn't list sites: %s", err)
		os.Exit(1)
	}

	fmt.Printf("%+v", list)

}
