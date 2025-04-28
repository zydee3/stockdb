package fmp

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	httpUtil "github.com/zydee3/stockdb/internal/api/utility"
)

type HTTPClient struct {
	client httpUtil.HTTPClient
	apiKey string
}

const (
	fmpURL = "https://financialmodelingprep.com/stable"
)

func (h *HTTPClient) Get(endpoint string, data map[string]string) (*http.Response, error) {
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
