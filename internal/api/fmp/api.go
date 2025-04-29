package fmp

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	httpUtil "github.com/zydee3/stockdb/internal/api/utility"
)

// TODO: Oscar - This should be named FMPClient or something similar. If you're
// using this pattern for all APIs, then you should make it an interface so that
// the workers can use any client abstractly.

type HTTPClient struct {
	client httpUtil.HTTPClient
	apiKey string
}

func NewHTTPClient(client httpUtil.HTTPClient, apiKey string) *HTTPClient {
	return &HTTPClient{
		client: client,
		apiKey: apiKey,
	}
}

func (h *HTTPClient) Get(endpoint string, data map[string]string) (*http.Response, error) {
	const (
		fmpURL = "https://financialmodelingprep.com/stable"
	)

	if data == nil {
		data = map[string]string{}
	}

	data["api_key"] = h.apiKey

	var dataStringBuilder strings.Builder
	for key, value := range data {
		dataStringBuilder.WriteString(key)
		dataStringBuilder.WriteString("=")
		dataStringBuilder.WriteString(value)
		dataStringBuilder.WriteString("&")
	}

	endpoint = fmt.Sprintf("%s/%s?%s", fmpURL, endpoint, dataStringBuilder.String())

	ctx := context.Background()

	return h.client.Get(ctx, endpoint)
}
