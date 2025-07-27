# groundcover Go SDK Examples

This directory contains practical examples demonstrating how to use the groundcover Go SDK for common monitoring use cases.

## Getting Started

Before running any examples, make sure you have set the required environment variables:

```bash
export GC_API_KEY="your-api-key"
export GC_BACKEND_ID="your-backend-id"
```

Optionally, you can override the default base URL:

```bash
export GC_BASE_URL="https://your-custom-api.groundcover.com"
```

You can also set an optional traceparent for distributed tracing:

```bash
export GC_TRACEPARENT="your-traceparent-header"
```

## Running Examples

Each example is a standalone Go program. To run an example:

```bash
cd examples
go run logs/main.go
```

Or build and run:

```bash
cd examples/logs
go build -o example
./example
```

## Available Examples

### **[logs](./logs/)** - Search and filter logs
Search for error-level logs with color-coded output formatting.

### **[metrics](./metrics/)** - Query metrics using PromQL
Execute PromQL queries to get cluster metrics like node count.

### **[events](./events/)** - Event detection
Monitor and display Kubernetes events with detailed metadata. This example demonstrates searching for OOM (Out of Memory) events.

## Prerequisites

- Go 1.24 or higher
- Valid groundcover API credentials
- Access to a groundcover environment

## Documentation

For more detailed information, refer to:
- [groundcover API Documentation](https://docs.groundcover.com/use-groundcover/remote-access-and-apis/apis)
- [SDK README](../README.md) 