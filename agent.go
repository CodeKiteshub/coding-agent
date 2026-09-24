package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// NewAgent creates a new agent instance with all tools and MCP servers configured
func NewAgent(getUserMessage func() (string, bool), verbose bool) *Agent {
	agent := &Agent{
		getUserMessage: getUserMessage,
		verbose:        verbose,
		conversation:   []Message{},
		mcpServers:     []MCPServer{},
	}

	// Register built-in tools
	agent.tools = []ToolDefinition{
		{
			Name:        "read_file",
			Description: "Read the contents of a file",
			Function:    ReadFile,
		},
		{
			Name:        "list_files",
			Description: "List files and directories in a given path",
			Function:    ListFiles,
		},
		{
			Name:        "edit_file",
			Description: "Edit a file by replacing old_str with new_str",
			Function:    EditFile,
		},
		{
			Name:        "bash",
			Description: "Execute a bash command and return the output",
			Function:    BashCommand,
		},
		{
			Name:        "code_search",
			Description: "Search for code patterns in files using grep",
			Function:    CodeSearch,
		},
		{
			Name:        "mcp_tool",
			Description: "Execute a tool from an MCP server",
			Function:    agent.MCPToolCall,
		},
		{
			Name:        "list_mcp_servers",
			Description: "List available MCP servers and their tools",
			Function:    agent.ListMCPServers,
		},
	}

	// Initialize MCP servers
	agent.initializeMCPServers()

	return agent
}

// initializeMCPServers sets up all available MCP servers
func (a *Agent) initializeMCPServers() {
	// Create MCP servers using factory functions
	fsServer := createFilesystemServer()
	webServer := createWebSearchServer()
	dbServer := createDatabaseServer()
	dartServer := createDartServer()

	a.mcpServers = []MCPServer{fsServer, webServer, dbServer, dartServer}

	if a.verbose {
		log.Printf("Initialized %d MCP servers", len(a.mcpServers))
	}
}

// Run starts the agent's main conversation loop
func (a *Agent) Run(ctx context.Context) error {
	a.printWelcomeMessage()
	a.addSystemMessage()

	for {
		fmt.Print("You: ")
		userMessage, ok := a.getUserMessage()
		if !ok {
			break
		}

		if userMessage == "" {
			continue
		}

		// Add user message to conversation
		a.conversation = append(a.conversation, Message{
			Role:    "user",
			Content: userMessage,
		})

		if a.verbose {
			log.Printf("User message added to conversation: %s", userMessage)
		}

		// Process the conversation with tool handling
		err := a.processConversation(ctx)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			continue
		}
	}

	return nil
}

// printWelcomeMessage displays the startup information
func (a *Agent) printWelcomeMessage() {
	fmt.Println("🚀 Super Agent with MCP Support + Llama 3.2 (use 'ctrl-c' to quit)")
	fmt.Println("Built-in tools: read_file, list_files, edit_file, bash, code_search")
	fmt.Println("MCP servers: filesystem, web-search, database, dart")
	fmt.Println("Model: llama3.2:3b (running locally via Ollama)")
	fmt.Println()
}

// addSystemMessage creates and adds the system message with tool descriptions
func (a *Agent) addSystemMessage() {
	systemMsg := "You are a helpful AI assistant with access to the following built-in tools:\n"
	for _, tool := range a.tools {
		if tool.Name != "mcp_tool" && tool.Name != "list_mcp_servers" {
			systemMsg += fmt.Sprintf("- %s: %s\n", tool.Name, tool.Description)
		}
	}
	
	systemMsg += "\nYou also have access to MCP (Model Context Protocol) servers with additional tools:\n"
	for _, server := range a.mcpServers {
		systemMsg += fmt.Sprintf("- %s server:\n", server.Name)
		for _, tool := range server.Tools {
			systemMsg += fmt.Sprintf("  - %s: %s\n", tool.Name, tool.Description)
		}
	}

	systemMsg += "\nTo use built-in tools: tool: tool_name({\"param\": \"value\"})\n"
	systemMsg += "To use MCP tools: tool: mcp_tool({\"server_name\": \"server\", \"tool_name\": \"tool\", \"arguments\": {\"param\": \"value\"}})\n"
	systemMsg += "To list MCP servers: tool: list_mcp_servers({})\n"
	systemMsg += "You can use multiple tools in sequence to complete complex tasks."

	a.conversation = append(a.conversation, Message{
		Role:    "system",
		Content: systemMsg,
	})
}

// processConversation handles the conversation flow with tool execution
func (a *Agent) processConversation(ctx context.Context) error {
	for {
		// Get response from Ollama
		response, err := a.runInference(ctx)
		if err != nil {
			return err
		}

		// Check if response contains tool calls
		if strings.Contains(response, "tool:") {
			// Execute tools and continue conversation
			toolResults := a.executeTools(response)
			
			// Add assistant response with tool calls
			a.conversation = append(a.conversation, Message{
				Role:    "assistant",
				Content: response,
			})

			// Add tool results
			a.conversation = append(a.conversation, Message{
				Role:    "user",
				Content: "Tool results:\n" + toolResults,
			})

			if a.verbose {
				log.Printf("Tool results added to conversation")
			}
		} else {
			// Final response, add to conversation and display
			a.conversation = append(a.conversation, Message{
				Role:    "assistant",
				Content: response,
			})

			fmt.Printf("Agent: %s\n\n", response)
			break
		}
	}

	return nil
}

// executeTools parses and executes tool calls from the response
func (a *Agent) executeTools(response string) string {
	lines := strings.Split(response, "\n")
	var results []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "tool:") {
			// Parse tool call: tool: tool_name({"param": "value"})
			toolCall := strings.TrimPrefix(line, "tool:")
			toolCall = strings.TrimSpace(toolCall)

			// Extract tool name and parameters
			parenIndex := strings.Index(toolCall, "(")
			if parenIndex == -1 {
				results = append(results, "Error: Invalid tool call format")
				continue
			}

			toolName := strings.TrimSpace(toolCall[:parenIndex])
			paramsStr := toolCall[parenIndex+1:]
			if strings.HasSuffix(paramsStr, ")") {
				paramsStr = paramsStr[:len(paramsStr)-1]
			}

			// Find and execute the tool
			var toolResult string
			found := false
			for _, tool := range a.tools {
				if tool.Name == toolName {
					result, err := tool.Function(json.RawMessage(paramsStr))
					if err != nil {
						toolResult = fmt.Sprintf("Error executing %s: %s", toolName, err.Error())
					} else {
						toolResult = result
					}
					found = true
					break
				}
			}

			if !found {
				toolResult = fmt.Sprintf("Error: Unknown tool '%s'", toolName)
			}

			results = append(results, fmt.Sprintf("%s: %s", toolName, toolResult))

			if a.verbose {
				log.Printf("Executed tool %s with result length: %d", toolName, len(toolResult))
			}

			// Also print tool execution for user visibility
			fmt.Printf("tool: %s(%s)\n", toolName, paramsStr)
			fmt.Printf("result: %s\n", toolResult)
		}
	}

	return strings.Join(results, "\n")
}

// runInference makes a request to the Ollama API
func (a *Agent) runInference(ctx context.Context) (string, error) {
	if a.verbose {
		log.Printf("Making Ollama API call with %d messages", len(a.conversation))
	}

	request := OllamaRequest{
		Model:    "llama3.2:3b",
		Messages: a.conversation,
		Stream:   false,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to make request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var ollamaResp OllamaResponse
	err = json.Unmarshal(body, &ollamaResp)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if a.verbose {
		log.Printf("Ollama API response received, content length: %d characters", len(ollamaResp.Message.Content))
	}

	return ollamaResp.Message.Content, nil
}
