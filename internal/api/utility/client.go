package utility

import (
	"io"
	"net/http"
	"time"
)

type HTTPClient struct {
	Client        *http.Client
	RetryCount    int
	RetryWaitTime time.Duration
}

func (h *HTTPClient) Do(request *http.Request) (*http.Response, error) {
	var response *http.Response
	var err error

	retries := h.RetryCount
	for retries > 0 {
		response, err = h.Client.Do(request)
		if err == nil {
			break
		}
		time.Sleep(h.RetryWaitTime)
		retries--
	}

	return response, err
}

func (h *HTTPClient) Get(url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return h.Do(request)
}

func (h *HTTPClient) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", contentType)
	return h.Do(request)
}

func (h *HTTPClient) Put(url string, contentType string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", contentType)
	return h.Do(request)
}

func (h *HTTPClient) Patch(url string, contentType string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest("PATCH", url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	return h.Do(request)
}

func (h *HTTPClient) Delete(url string) (*http.Response, error) {
	request, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	return h.Do(request)
}
