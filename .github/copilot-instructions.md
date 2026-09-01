# MCP Server Instructions

This repository uses the official Go MCP SDK:

- https://github.com/modelcontextprotocol/go-sdk
- https://pkg.go.dev/github.com/modelcontextprotocol/go-sdk/mcp

Keep the MCP server independent from payment providers and PostgreSQL. Add payment capabilities by calling the payment platform REST API. Preserve idempotency keys, structured errors, explicit approval requirements for mutations, and never accept raw card data through MCP tools.
