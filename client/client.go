package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
)

const defaultBaseURL = "https://monitoringapi.solaredge.com"

type Client struct {
	client  http.Client
	baseURL url.URL
}

func NewClient(key string) *Client {
	log.Debug().Str("baseURL", defaultBaseURL).Msg("Setting up client")

	baseURL, err := url.Parse(defaultBaseURL)
	if err != nil {
		panic(err)
	}

	client := &Client{
		client: http.Client{
			Transport: AuthenticatingRoundTripper{
				key: key,
			},
		},
		baseURL: *baseURL,
	}

	return client
}

func handleResponse(response *http.Response, out interface{}) error {
	if response.StatusCode == 404 {
		return errors.New("no document found at endpoint")
	}
	if response.StatusCode == 429 {
		return errors.New("query limit exceeded, time to write that backoff logic")
	}
	if response.StatusCode == 401 {
		return errors.New("unauthorized response, check key validity")
	}
	if response.StatusCode == 403 {
		return errors.New("unauthorized response, check key permissions")
	}
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error parsing response body: %w", err)
	}
	if response.StatusCode > 299 {
		return fmt.Errorf("unexpected response code %d: %s", response.StatusCode, respBody)
	}

	log.Debug().Str("response body", string(respBody)).Msg("Response Body")
	err = json.Unmarshal(respBody, out)
	if err != nil {
		return fmt.Errorf("error unmarshalling response to %t: %w", out, err)
	}

	return nil
}

// FIXME: My happy place would be a good httpclient middleware framework, not a handjam like this.
type AuthenticatingRoundTripper struct {
	transport http.RoundTripper
	key       string
}

func (a AuthenticatingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	transport := a.transport
	if transport == nil {
		transport = http.DefaultTransport
	}

	log.Debug().Str("url", req.URL.String()).Msg("url pre-apikey injection")

	a.addAuthenticationKey(req)

	log.Trace().Str("url", req.URL.String()).Msg("url with api key")

	return transport.RoundTrip(req)
}

func (a *AuthenticatingRoundTripper) addAuthenticationKey(req *http.Request) {
	values := req.URL.Query()
	values.Set("api_key", a.key)
	req.URL.RawQuery = values.Encode()
	log.Trace().Str("url", req.URL.String()).Msg("url in addAuthenticationKey")
}
