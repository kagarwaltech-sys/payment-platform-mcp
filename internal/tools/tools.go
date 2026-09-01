package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"payment-platform-mcp/internal/paymentapi"
)

type GetPaymentInput struct {
	PaymentID string `json:"payment_id" jsonschema:"the local payment ID"`
}

type PayPaymentInput struct {
	Amount         int64  `json:"amount" jsonschema:"amount in minor currency units"`
	Currency       string `json:"currency" jsonschema:"three-letter ISO currency code"`
	CaptureMethod  string `json:"capture_method" jsonschema:"payment capture method, currently manual"`
	Reference      string `json:"reference,omitempty" jsonschema:"merchant order reference"`
	IdempotencyKey string `json:"idempotency_key,omitempty" jsonschema:"stable retry key for this payment request"`
}

func Register(server *mcp.Server, client *paymentapi.Client) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_payment",
		Description: "Retrieve the current state of a payment from the payment platform.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input GetPaymentInput) (*mcp.CallToolResult, paymentapi.Payment, error) {
		payment, err := client.GetPayment(ctx, input.PaymentID)
		return nil, payment, err
	})
	mcp.AddTool(server, &mcp.Tool{
		Name:        "pay_payment",
		Description: "Create, authorize, capture, and record a payment. This is a high-risk mutation.",
	}, func(ctx context.Context, _ *mcp.CallToolRequest, input PayPaymentInput) (*mcp.CallToolResult, paymentapi.Payment, error) {
		captureMethod := input.CaptureMethod
		if captureMethod == "" {
			captureMethod = "manual"
		}
		payment, err := client.PayPayment(ctx, input.IdempotencyKey, input.Amount, input.Currency, captureMethod, input.Reference)
		return nil, payment, err
	})
}
