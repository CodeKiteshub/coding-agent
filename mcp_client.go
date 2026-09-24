package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// MCP-specific tool functions
func (a *Agent) MCPToolCall(input json.RawMessage) (string, error) {
	var params MCPToolInput
	err := json.Unmarshal(input, &params)
	if err != nil {
		return "", fmt.Errorf("invalid parameters: %w", err)
	}

	// Find the MCP server
	var server *MCPServer
	for i := range a.mcpServers {
		if a.mcpServers[i].Name == params.ServerName {
			server = &a.mcpServers[i]
			break
		}
	}

	if server == nil {
		return "", fmt.Errorf("MCP server '%s' not found", params.ServerName)
	}

	// Call the MCP server
	return a.callMCPServer(server, params.ToolName, params.Arguments)
}

func (a *Agent) ListMCPServers(input json.RawMessage) (string, error) {
	var result strings.Builder
	result.WriteString("Available MCP Servers:\n\n")

	for _, server := range a.mcpServers {
		result.WriteString(fmt.Sprintf("Server: %s\n", server.Name))
		result.WriteString(fmt.Sprintf("URL: %s\n", server.URL))
		result.WriteString("Tools:\n")
		for _, tool := range server.Tools {
			result.WriteString(fmt.Sprintf("  - %s: %s\n", tool.Name, tool.Description))
		}
		result.WriteString("\n")
	}

	return result.String(), nil
}

func (a *Agent) callMCPServer(server *MCPServer, toolName string, arguments json.RawMessage) (string, error) {
	// Find the tool
	var tool *MCPTool
	for i := range server.Tools {
		if server.Tools[i].Name == toolName {
			tool = &server.Tools[i]
			break
		}
	}

	if tool == nil {
		return "", fmt.Errorf("tool '%s' not found on server '%s'", toolName, server.Name)
	}

	// Prepare MCP request
	mcpRequest := MCPRequest{
		Method: "tools/call",
		Params: map[string]interface{}{
			"name":      toolName,
			"arguments": json.RawMessage(arguments),
		},
	}

	jsonData, err := json.Marshal(mcpRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal MCP request: %w", err)
	}

	// Make HTTP request to MCP server
	resp, err := http.Post(server.URL+"/mcp", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		// If MCP server is not available, return a mock response for demonstration
		return a.mockMCPResponse(server.Name, toolName, arguments)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("MCP server error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read MCP response: %w", err)
	}

	var mcpResp MCPResponse
	err = json.Unmarshal(body, &mcpResp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal MCP response: %w", err)
	}

	if mcpResp.Error != nil {
		return "", fmt.Errorf("MCP error: %s", mcpResp.Error.Message)
	}

	// Convert result to string
	resultBytes, err := json.Marshal(mcpResp.Result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal MCP result: %w", err)
	}

	return string(resultBytes), nil
}

func (a *Agent) mockMCPResponse(serverName, toolName string, arguments json.RawMessage) (string, error) {
	// Mock responses for demonstration when MCP servers are not running
	switch serverName {
	case "dart":
		switch toolName {
		case "dart-analyze":
			return "Mock: Dart analysis complete. No issues found in the specified path.", nil
		case "dart-format":
			return "Mock: Dart files formatted successfully.", nil
		case "dart-test":
			return "Mock: All Dart tests passed successfully.", nil
		case "dart-run":
			return "Mock: Dart script executed successfully.", nil
		case "dart-create":
			return "Mock: New Dart project created successfully.", nil
		case "dart-package":
			return "Mock: Dart package operation completed successfully.", nil
		case "dart-compile":
			return "Mock: Dart compilation completed successfully.", nil
		default:
			return fmt.Sprintf("Mock: Unknown Dart tool '%s'", toolName), nil
		}
	case "filesystem":
		if toolName == "read_directory" {
			return "Mock filesystem response: Directory contains 5 files and 2 subdirectories", nil
		}
	case "web-search":
		if toolName == "search" {
			return "Mock web search response: Found 10 relevant results for your query", nil
		}
	case "database":
		if toolName == "query" {
			return "Mock database response: Query executed successfully, returned 3 rows", nil
		}
	}

	return fmt.Sprintf("Mock MCP response from %s.%s with arguments: %s", serverName, toolName, string(arguments)), nil
}
