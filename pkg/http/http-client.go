package http

import (
	"bytes"
	"net/http"
)

func HttpClient(baseURL string) *BaseHttpClient {
	return &BaseHttpClient{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}
}

func (c *BaseHttpClient) SendGetRequest(uri string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.BaseURL+uri, nil)
	if err != nil {
		return nil, err
	}
	c._setCleanedHeaders(req)

	return c.HTTPClient.Do(req)
}

func (c *BaseHttpClient) SendPostRequest(uri, payload string) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.BaseURL+uri, bytes.NewBufferString(payload))
	if err != nil {
		return nil, err
	}
	c._setCleanedHeaders(req)

	return c.HTTPClient.Do(req)
}

func (c *BaseHttpClient) SendPatchRequest(uri, payload string) (*http.Response, error) {
	req, err := http.NewRequest("PATCH", c.BaseURL+uri, bytes.NewBufferString(payload))
	if err != nil {
		return nil, err
	}
	c._setCleanedHeaders(req)

	return c.HTTPClient.Do(req)
}

func (c *BaseHttpClient) _setCleanedHeaders(request *http.Request) {
	for key, value := range c.Headers {
		request.Header.Set(key, value)
	}
}
