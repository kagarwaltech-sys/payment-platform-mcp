package main

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"payment-platform-mcp/internal/paymentapi"
	"payment-platform-mcp/internal/tools"
)

func main() {
	apiURL := os.Getenv("PAYMENT_API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}
	server := mcp.NewServer(&mcp.Implementation{Name: "payment-platform-mcp", Version: "0.1.0"}, nil)
	tools.Register(server, paymentapi.New(apiURL, os.Getenv("PAYMENT_API_TOKEN")))
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
