// Package tools defines the MCP tools the server exposes.
package tools

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/tredy-ai/rogue-goose-go/internal/orbitalwiki"
)

// Register adds every tool to the server.
func Register(s *mcp.Server, api *orbitalwiki.Client) {
	registerSearchTool(s, api)
	// Part 2: register your tool here.
}
