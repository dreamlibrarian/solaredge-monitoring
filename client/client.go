package client

import "net/http"

const defaultBaseURL = "https://monitoringapi.solaredge.com"

type Client struct {
	http.Client
}

func NewClient(key string) {

}