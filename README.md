# Payment Platform MCP Server

MCP server for the payment platform API. This repository owns the MCP interface and calls the payment API over HTTP. It does not access PostgreSQL or provider APIs directly.

## Initial tools

- `get_payment`: read the current state of a payment.
- `pay_payment`: execute the POS-style create, authorize, capture, and ledger flow. Treat this as a high-risk mutation and require approval in an agent.

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

The API token is passed as a bearer token when configured. Payment API idempotency keys are generated when `pay_payment` does not receive one.

## Architecture

See the payment platform repository document `docs/08-ai-mcp-architecture.md` for repository boundaries, safety rules, and the agent delivery plan.

## Validation

```powershell
go test ./...
go vet ./...
```
