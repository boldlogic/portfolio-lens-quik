package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
}

func NewClient(cfg HttpClientConfig) *Client {
	tr := &http.Transport{
		MaxIdleConnsPerHost: 1,
		IdleConnTimeout:     30 * time.Second,
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(cfg.Timeout) * time.Second,
	}
	return &Client{client: httpClient}
}

type Response struct {
	StatusCode int
	Body       []byte
}

// SendRequest выполняет единичный HTTP-Запрос. Возвращает StatusCode int, Body []byte, error
func (c *Client) SendRequest(ctx context.Context, req *http.Request) (int, []byte, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("не удалось выполнить запрос: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("не удалось прочитать тело ответа: %w", err)
	}
	return resp.StatusCode, body, nil
}

// SendRequest выполняет HTTP-Запрос с повторами. Возвращает StatusCode int,Body []byte, attempts int, error
func (c *Client) SendWithRetry(ctx context.Context, req *http.Request, retryCount int) (int, []byte, int, error) {
	var statusCode int
	var body []byte
	var err error

	for attempt := 0; attempt <= retryCount; attempt++ {
		statusCode, body, err = c.SendRequest(ctx, req)

		if err == nil && statusCode == http.StatusOK {
			return statusCode, body, attempt + 1, nil
		}
	}
	return statusCode, nil, retryCount + 1, err
}
