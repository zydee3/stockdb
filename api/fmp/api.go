package fmp

import (
	"fmt"
	"github.com/zydee3/stockdb/api/common"
	"net/http"
	"strings"
)

type HTTPClient struct {
	HttpClient common.HTTPClient
	ApiKey     string
}

const (
	fmpUrl = "https://financialmodelingprep.com/stable"
)

func (h *HTTPClient) Get(endpoint string, data map[string]string) (*http.Response, error) {
	if data == nil {
		data = map[string]string{}
	}

	data["api_key"] = h.ApiKey

	var dataStringBuilder strings.Builder
	for key, value := range data {
		dataStringBuilder.WriteString(key)
		dataStringBuilder.WriteString("=")
		dataStringBuilder.WriteString(value)
		dataStringBuilder.WriteString("&")
	}

	endpoint = fmt.Sprintf("%s/%s?%s", fmpUrl, endpoint, dataStringBuilder.String())

	return h.HttpClient.Get(endpoint)
}
