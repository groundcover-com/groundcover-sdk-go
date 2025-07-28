// Package groundcover provides the groundcover SDK client.
package groundcover

import (
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/option"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/transport"
)

// NewClient creates a new groundcover SDK client with a simplified API.
// It automatically reads configuration from environment variables unless overridden by options.
//
// Environment variables:
//   - GC_API_KEY: Your groundcover API key (required)
//   - GC_BACKEND_ID: Your groundcover Backend ID (required)
//   - GC_BASE_URL: The base URL of the groundcover API (optional, defaults to https://api.groundcover.com)
//
// Example usage:
//
//	client := groundcover.NewClient()
//	client := groundcover.NewClient(option.WithAPIKey("custom-key"))
func NewClient(options ...option.Option) (*client.GroundcoverAPI, error) {
	return transport.NewClient(options...)
}
