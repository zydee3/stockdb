package utility

import (
	"context"
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

func (h *HTTPClient) Get(ctx context.Context, url string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return h.Do(request)
}

func (h *HTTPClient) Post(ctx context.Context, url string, contentType string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", contentType)
	return h.Do(request)
}

func (h *HTTPClient) Put(ctx context.Context, url string, contentType string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", contentType)
	return h.Do(request)
}

func (h *HTTPClient) Patch(
	ctx context.Context,
	url string,
	contentType string,
	body io.Reader,
) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	return h.Do(request)
}

func (h *HTTPClient) Delete(ctx context.Context, url string) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	return h.Do(request)
}
