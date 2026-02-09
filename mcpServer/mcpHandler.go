package mcpServer

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// NewServer creates an MCP server configured with tools for mark3labs mcp-go.
func NewServer() *server.MCPServer {
	s := server.NewMCPServer("Andreas API MCP", "1.0.0", server.WithToolCapabilities(true))

	tool := mcp.NewTool(
		"hello",
		mcp.WithDescription("Hello MCP tool"),
		mcp.WithString("message", mcp.Description("Optional message to echo back")),
	)
	s.AddTool(tool, helloHandler)

	return s
}

func helloHandler(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	msg, _ := args["message"].(string)
	if msg == "" {
		msg = "(no message)"
	}
	return mcp.NewToolResultText(fmt.Sprintf("Hello MCP! You sent: %s", msg)), nil
}
