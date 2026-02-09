// MCP client that connects to the Andreas API MCP server (Streamable HTTP)
// and lists tools / calls the hello tool.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

const defaultBaseURL = "http://localhost:8081/mcp"

func main() {
	baseURL := flag.String("url", defaultBaseURL, "MCP server base URL (Streamable HTTP)")
	message := flag.String("message", "from MCP client", "message to send to the hello tool")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c, err := client.NewStreamableHttpClient(*baseURL)
	if err != nil {
		log.Fatalf("Failed to create MCP client: %v", err)
	}
	defer c.Close()

	if err := c.Start(ctx); err != nil {
		log.Fatalf("Failed to start client: %v", err)
	}

	initReq := mcp.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initReq.Params.Capabilities = mcp.ClientCapabilities{}
	initReq.Params.ClientInfo = mcp.Implementation{
		Name:    "Andreas API MCP Client",
		Version: "1.0.0",
	}

	serverInfo, err := c.Initialize(ctx, initReq)
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}

	fmt.Printf("Connected to: %s (version %s)\n", serverInfo.ServerInfo.Name, serverInfo.ServerInfo.Version)

	if serverInfo.Capabilities.Tools == nil {
		log.Fatal("Server does not support tools")
	}

	toolsResult, err := c.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}

	fmt.Printf("Available tools (%d):\n", len(toolsResult.Tools))
	for _, t := range toolsResult.Tools {
		fmt.Printf("  - %s: %s\n", t.Name, t.Description)
	}

	callReq := mcp.CallToolRequest{}
	callReq.Params.Name = "hello"
	callReq.Params.Arguments = map[string]interface{}{
		"message": *message,
	}

	result, err := c.CallTool(ctx, callReq)
	if err != nil {
		log.Fatalf("Failed to call tool: %v", err)
	}

	if result.IsError {
		fmt.Fprintf(os.Stderr, "Tool error: %s\n", result.Content)
		os.Exit(1)
	}

	fmt.Println("Result:")
	for _, content := range result.Content {
		fmt.Println(mcp.GetTextFromContent(content))
	}
}
