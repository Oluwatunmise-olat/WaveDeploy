package http

import "net/http"

type BaseHttpClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Headers    map[string]string
}
