package main

import (
	"encoding/json"
)

// Core agent structures
type Agent struct {
	getUserMessage func() (string, bool)
	verbose        bool
	conversation   []Message
	tools          []ToolDefinition
	mcpServers     []MCPServer
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ToolDefinition struct {
	Name        string
	Description string
	Function    func(input json.RawMessage) (string, error)
}

// MCP Server structures
type MCPServer struct {
	Name        string
	URL         string
	Tools       []MCPTool
	Resources   []MCPResource
}

type MCPTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

type MCPResource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MimeType    string `json:"mimeType"`
}

type MCPRequest struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type MCPResponse struct {
	Result interface{} `json:"result"`
	Error  *MCPError   `json:"error"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Ollama API structures
type OllamaRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type OllamaResponse struct {
	Message Message `json:"message"`
	Done    bool    `json:"done"`
}

// Tool input structures
type ReadFileInput struct {
	Path string `json:"path"`
}

type ListFilesInput struct {
	Path string `json:"path"`
}

type EditFileInput struct {
	Path   string `json:"path"`
	OldStr string `json:"old_str"`
	NewStr string `json:"new_str"`
}

type BashInput struct {
	Command string `json:"command"`
}

type CodeSearchInput struct {
	Pattern       string `json:"pattern"`
	Path          string `json:"path"`
	FileType      string `json:"file_type"`
	CaseSensitive bool   `json:"case_sensitive"`
}

type MCPToolInput struct {
	ServerName string          `json:"server_name"`
	ToolName   string          `json:"tool_name"`
	Arguments  json.RawMessage `json:"arguments"`
}
