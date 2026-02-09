// LLM chat: OpenAI with MCP tools. Env: OPENAI_API_KEY (required); OPENAI_CHAT_MODEL (optional, default gpt-4o).
package services

import (
	"context"
	"encoding/json"
	"os"
	"sync"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

// ChatService runs an LLM (OpenAI) with access to MCP server tools.
type ChatService struct {
	openaiClient *openai.Client
	mcpServer    *server.MCPServer
	mcpClient    *client.Client
	mcpInitOnce  sync.Once
	mcpInitErr   error
}

// NewChatService creates a chat service that uses the given MCP server for tools.
// OpenAI API key is read from OPENAI_API_KEY env (set later by caller).
func NewChatService(mcpServer *server.MCPServer) *ChatService {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = "not-set" // client still created; calls will fail until env is set
	}
	return &ChatService{
		openaiClient: openai.NewClient(apiKey),
		mcpServer:    mcpServer,
		mcpClient:    nil,
	}
}

// ensureMCPClient initializes the in-process MCP client once (Start + Initialize).
func (s *ChatService) ensureMCPClient(ctx context.Context) error {
	if s.mcpClient != nil {
		return nil
	}
	if s.mcpServer == nil {
		return nil
	}
	s.mcpInitOnce.Do(func() {
		c, err := client.NewInProcessClient(s.mcpServer)
		if err != nil {
			s.mcpInitErr = err
			return
		}
		s.mcpClient = c
		if err := c.Start(ctx); err != nil {
			s.mcpInitErr = err
			return
		}
		initReq := mcp.InitializeRequest{}
		initReq.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
		initReq.Params.Capabilities = mcp.ClientCapabilities{}
		initReq.Params.ClientInfo = mcp.Implementation{Name: "LLMChat", Version: "1.0.0"}
		_, err = c.Initialize(ctx, initReq)
		if err != nil {
			s.mcpInitErr = err
			return
		}
	})
	return s.mcpInitErr
}

// mcpToolsToOpenAI converts MCP list_tools result to OpenAI tools slice.
func mcpToolsToOpenAI(tools []mcp.Tool) []openai.Tool {
	out := make([]openai.Tool, 0, len(tools))
	for _, t := range tools {
		out = append(out, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        t.Name,
				Description: t.Description,
				Parameters: jsonschema.Definition{
					Type:       jsonschema.Object,
					Properties: map[string]jsonschema.Definition{},
				},
			},
		})
	}
	return out
}

// Chat sends a single message to OpenAI with MCP tools and returns the reply (no thread history).
func (s *ChatService) Chat(ctx context.Context, message string) (string, error) {
	return s.ChatWithHistory(ctx, nil, message)
}

// ChatWithHistory sends existing conversation + new user message to OpenAI and returns the assistant reply.
// history is the prior messages in OpenAI format (user/assistant only; no tool_calls). Can be nil.
func (s *ChatService) ChatWithHistory(ctx context.Context, history []openai.ChatCompletionMessage, newUserMessage string) (string, error) {
	if err := s.ensureMCPClient(ctx); err != nil {
		return "", err
	}

	toolsResult, err := s.mcpClient.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return "", err
	}
	openaiTools := mcpToolsToOpenAI(toolsResult.Tools)

	messages := make([]openai.ChatCompletionMessage, 0, len(history)+1)
	if len(history) > 0 {
		messages = append(messages, history...)
	}
	messages = append(messages, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: newUserMessage})

	model := os.Getenv("OPENAI_CHAT_MODEL")
	if model == "" {
		model = openai.GPT4o
	}

	for {
		req := openai.ChatCompletionRequest{
			Model:    model,
			Messages: messages,
			Tools:    openaiTools,
		}
		resp, err := s.openaiClient.CreateChatCompletion(ctx, req)
		if err != nil {
			return "", err
		}
		if len(resp.Choices) == 0 {
			return "", nil
		}
		msg := resp.Choices[0].Message

		if len(msg.ToolCalls) == 0 {
			return msg.Content, nil
		}

		// Append assistant message with tool calls
		messages = append(messages, msg)

		// Execute each tool call and append tool results
		for _, tc := range msg.ToolCalls {
			var args map[string]interface{}
			if tc.Function.Arguments != "" {
				_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)
			}
			if args == nil {
				args = make(map[string]interface{})
			}
			callReq := mcp.CallToolRequest{}
			callReq.Params.Name = tc.Function.Name
			callReq.Params.Arguments = args

			result, err := s.mcpClient.CallTool(ctx, callReq)
			if err != nil {
				result = &mcp.CallToolResult{IsError: true, Content: []mcp.Content{mcp.NewTextContent(err.Error())}}
			}
			var contentStr string
			for _, c := range result.Content {
				contentStr += mcp.GetTextFromContent(c)
			}
			messages = append(messages, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    contentStr,
				ToolCallID: tc.ID,
			})
		}
	}
}
