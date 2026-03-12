package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const defaultTurnBaseURL = "https://whatsapp.turn.io"

// TurnClient is a REST client for the turn.io WhatsApp API with rate limit tracking.
type TurnClient struct {
	mu                 sync.Mutex
	rateLimitRemaining int
	baseURL            string
	apiToken           string
	httpClient         *http.Client
}

// NewTurnClient creates a new TurnClient.
func NewTurnClient(apiToken string, baseURL string) *TurnClient {
	if baseURL == "" {
		baseURL = defaultTurnBaseURL
	}
	return &TurnClient{
		baseURL:    baseURL,
		apiToken:   apiToken,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Do performs an HTTP request and returns the raw response body.
func (c *TurnClient) Do(ctx context.Context, method, path string, body any) (json.RawMessage, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Track rate limit budget from response headers.
	if remaining := resp.Header.Get("X-Ratelimit-Remaining"); remaining != "" {
		if n, err := strconv.Atoi(remaining); err == nil {
			c.mu.Lock()
			c.rateLimitRemaining = n
			c.mu.Unlock()
		}
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		return nil, fmt.Errorf("rate limited (HTTP 429), retry after: %s", retryAfter)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return json.RawMessage(respBody), nil
}

// DoInto performs an HTTP request and unmarshals the response into target.
func (c *TurnClient) DoInto(ctx context.Context, method, path string, body, target any) error {
	data, err := c.Do(ctx, method, path, body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

// RateLimitRemaining returns the last known rate limit budget.
func (c *TurnClient) RateLimitRemaining() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.rateLimitRemaining
}
