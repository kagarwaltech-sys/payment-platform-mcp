package paymentapi

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func New(baseURL, token string) *Client {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), token: token, http: http.DefaultClient}
}

type Payment struct {
	ID                 string `json:"id"`
	Amount             int64  `json:"amount"`
	Currency           string `json:"currency"`
	AuthorizedAmount   int64  `json:"authorized_amount"`
	CapturedAmount     int64  `json:"captured_amount"`
	Status             string `json:"status"`
	CaptureMethod      string `json:"capture_method"`
	Reference          string `json:"reference,omitempty"`
	Processor          string `json:"processor,omitempty"`
	ProcessorPaymentID string `json:"processor_payment_id,omitempty"`
	ProcessorCaptureID string `json:"processor_capture_id,omitempty"`
}

func (c *Client) GetPayment(ctx context.Context, id string) (Payment, error) {
	var result Payment
	err := c.do(ctx, http.MethodGet, "/api/v1/payments/"+id, "", nil, &result)
	return result, err
}

func (c *Client) PayPayment(ctx context.Context, idempotencyKey string, amount int64, currency, captureMethod, reference string) (Payment, error) {
	if idempotencyKey == "" {
		var err error
		idempotencyKey, err = newIdempotencyKey()
		if err != nil {
			return Payment{}, err
		}
	}
	body := map[string]any{"amount": amount, "currency": currency, "capture_method": captureMethod, "reference": reference}
	var result Payment
	err := c.do(ctx, http.MethodPost, "/api/v1/payments/pay", idempotencyKey, body, &result)
	return result, err
}

func (c *Client) do(ctx context.Context, method, path, idempotencyKey string, body any, result any) error {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	if idempotencyKey != "" {
		req.Header.Set("Idempotency-Key", idempotencyKey)
	}
	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("payment API returned HTTP %d: %s", res.StatusCode, strings.TrimSpace(string(data)))
	}
	if err := json.Unmarshal(data, result); err != nil {
		return fmt.Errorf("decode payment API response: %w", err)
	}
	return nil
}

func newIdempotencyKey() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return "mcp-" + hex.EncodeToString(value), nil
}
