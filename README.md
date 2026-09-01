# Payment Platform MCP Server

MCP server for the payment platform API. This repository owns the MCP interface and calls the payment API over HTTP. It does not access PostgreSQL or provider APIs directly.

## Tools

- `get_payment`: read the current state of a payment.
- `pay_payment`: execute the POS-style create, authorize, capture, and ledger flow. Treat this as a high-risk mutation and require approval in an agent.
- `refund_payment`: refund part or all of a captured payment. Treat this as a high-risk mutation and require explicit approval in an agent.

## Run locally

Requires Go 1.26+.

```powershell
$env:PAYMENT_API_URL='http://localhost:8080'
go run ./cmd/server
```

The server uses stdio transport. Configure it in an MCP client such as VS Code using `.vscode/mcp.json`.

## Configuration

```env
PAYMENT_API_URL=http://localhost:8080
PAYMENT_API_TOKEN=
```

The API token is passed as a bearer token when configured. Payment API idempotency keys are generated when `pay_payment` or `refund_payment` does not receive one.

## Codebase

```text
payment-platform-mcp/
├── cmd/server/main.go             # stdio MCP server entrypoint
├── internal/tools/tools.go        # tool schemas and handlers
├── internal/paymentapi/client.go  # authenticated REST client and payment models
└── internal/paymentapi/client_test.go
```

The API remains responsible for validation, idempotency, processor calls, payment state, and ledger behavior. This component only translates MCP calls into authenticated HTTP requests and returns structured results.

## Docker

```powershell
docker build -t payment-platform-mcp:local .
docker run --rm -i -e PAYMENT_API_URL=http://host.docker.internal:8080 payment-platform-mcp:local
```

The image contains a non-root, statically built server. It still uses stdio, so the MCP host must keep stdin/stdout attached.

## Security notes

- Never pass raw card data to an MCP tool.
- Require explicit approval in the agent before invoking payment or refund mutations.
- Prefer restricted payment API credentials.
- Treat API responses as authoritative; do not infer payment success from a tool request alone.

## Architecture

See the payment platform repository document `docs/08-ai-mcp-architecture.md` for repository boundaries, safety rules, and the agent delivery plan.

## Validation

```powershell
go test ./...
go vet ./...
```
