package paymentapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRefundPayment(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/payments/pay-123/refund" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Idempotency-Key"); got != "refund-key" {
			t.Fatalf("idempotency key = %q", got)
		}
		var body map[string]int64
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["amount"] != 4000 {
			t.Fatalf("body = %#v", body)
		}
		_, _ = w.Write([]byte(`{"id":"pay-123","captured_amount":10000,"refunded_amount":4000,"status":"CAPTURED","processor_refund_id":"re-123"}`))
	}))
	defer server.Close()

	payment, err := New(server.URL, "token").RefundPayment(context.Background(), "pay-123", "refund-key", 4000)
	if err != nil {
		t.Fatal(err)
	}
	if payment.RefundedAmount != 4000 || payment.ProcessorRefundID != "re-123" {
		t.Fatalf("payment = %#v", payment)
	}
}

func TestRefundPaymentGeneratesIdempotencyKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(r.Header.Get("Idempotency-Key")) == 0 {
			t.Fatal("missing generated idempotency key")
		}
		_, _ = w.Write([]byte(`{"id":"pay-123"}`))
	}))
	defer server.Close()

	if _, err := New(server.URL, "").RefundPayment(context.Background(), "pay-123", "", 1); err != nil {
		t.Fatal(err)
	}
}
