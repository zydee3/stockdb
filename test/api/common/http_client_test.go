package common_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	common2 "github.com/zydee3/stockdb/internal/api/utility"
)

type MockRoundTripper struct {
	Resp *http.Response
	Err  error
}

func (m *MockRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return m.Resp, m.Err
}

type MultiMockRoundTripper struct {
	Mocks []MockRoundTripper
}

func (m *MultiMockRoundTripper) RoundTrip(rq *http.Request) (*http.Response, error) {
	poped := m.Mocks[len(m.Mocks)-1]
	m.Mocks = m.Mocks[:len(m.Mocks)-1]
	return poped.RoundTrip(rq)
}

func TestHttpClient(t *testing.T) {
	t.Run("Test Basic Request", func(t *testing.T) {
		mock := &MockRoundTripper{Resp: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`OK`)),
		}}

		client := common2.HTTPClient{
			Client:        &http.Client{Transport: mock},
			RetryCount:    2,
			RetryWaitTime: time.Millisecond * 500,
		}

		ctx := context.Background()

		get, err := client.Get(ctx, "test.com")
		if err != nil {
			t.Errorf("Error while calling test.com: %v", err)
		}

		if get.StatusCode != http.StatusOK {
			t.Errorf("Invalid status code: %d", get.StatusCode)
		}

		if get.Body == nil {
			t.Errorf("Invalid response body")
		}

		body, err := io.ReadAll(get.Body)
		if err != nil {
			t.Errorf("Error while reading body: %v", err)
		}

		if string(body) != `OK` {
			t.Errorf("Invalid response body: %s", string(body))
		}
	})
}
