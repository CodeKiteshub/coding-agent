package main

// MCP Server configurations
func createFilesystemServer() MCPServer {
	return MCPServer{
		Name: "filesystem",
		URL:  "http://localhost:3001",
		Tools: []MCPTool{
			{
				Name:        "read_directory",
				Description: "Read directory contents with detailed metadata",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "Directory path to read",
						},
					},
					"required": []string{"path"},
				},
			},
		},
	}
}

func createWebSearchServer() MCPServer {
	return MCPServer{
		Name: "web-search",
		URL:  "http://localhost:3002",
		Tools: []MCPTool{
			{
				Name:        "search",
				Description: "Search the web for information",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "Search query",
						},
					},
					"required": []string{"query"},
				},
			},
		},
	}
}

func createDatabaseServer() MCPServer {
	return MCPServer{
		Name: "database",
		URL:  "http://localhost:3003",
		Tools: []MCPTool{
			{
				Name:        "query",
				Description: "Execute SQL query on database",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"sql": map[string]interface{}{
							"type":        "string",
							"description": "SQL query to execute",
						},
					},
					"required": []string{"sql"},
				},
			},
		},
	}
}

func createDartServer() MCPServer {
	return MCPServer{
		Name: "dart",
		URL:  "localhost:3004",
		Tools: []MCPTool{
			{
				Name:        "dart-analyze",
				Description: "Analyze Dart code for issues and warnings",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "Directory or file to analyze",
						},
						"options": map[string]interface{}{
							"type":        "array",
							"description": "Additional options for the dart analyze command",
							"items":       map[string]interface{}{"type": "string"},
						},
					},
					"required": []string{"path"},
				},
			},
			{
				Name:        "dart-format",
				Description: "Format Dart code files",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"paths": map[string]interface{}{
							"type":        "array",
							"description": "Files or directories to format",
							"items":       map[string]interface{}{"type": "string"},
						},
						"options": map[string]interface{}{
							"type":        "array",
							"description": "Additional format options",
							"items":       map[string]interface{}{"type": "string"},
						},
					},
					"required": []string{"paths"},
				},
			},
			{
				Name:        "dart-test",
				Description: "Run Dart tests",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "Path to the test file or directory",
						},
						"workingDir": map[string]interface{}{
							"type":        "string",
							"description": "Working directory for the command",
						},
						"options": map[string]interface{}{
							"type":        "array",
							"description": "Additional test options",
							"items":       map[string]interface{}{"type": "string"},
						},
					},
				},
			},
			{
				Name:        "dart-run",
				Description: "Run a Dart script",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"script": map[string]interface{}{
							"type":        "string",
							"description": "Path to the Dart script to run",
						},
						"args": map[string]interface{}{
							"type":        "array",
							"description": "Arguments to pass to the script",
							"items":       map[string]interface{}{"type": "string"},
						},
						"workingDir": map[string]interface{}{
							"type":        "string",
							"description": "Working directory for the command",
						},
					},
					"required": []string{"script"},
				},
			},
			{
				Name:        "dart-create",
				Description: "Create a new Dart project",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"projectName": map[string]interface{}{
							"type":        "string",
							"description": "Name of the project to create",
						},
						"output": map[string]interface{}{
							"type":        "string",
							"description": "Directory where to create the project",
						},
						"template": map[string]interface{}{
							"type":        "string",
							"description": "Template to use for project generation",
							"enum":        []string{"console", "package", "server-shelf", "web"},
							"default":     "package",
						},
						"options": map[string]interface{}{
							"type":        "array",
							"description": "Additional project creation options",
							"items":       map[string]interface{}{"type": "string"},
						},
					},
					"required": []string{"projectName", "output"},
				},
			},
			{
				Name:        "dart-package",
				Description: "Manage Dart packages (pub commands)",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"command": map[string]interface{}{
							"type":        "string",
							"description": "Pub subcommand to execute",
							"enum":        []string{"get", "upgrade", "outdated", "add", "remove", "publish", "deps", "downgrade", "cache", "run", "global"},
						},
						"args": map[string]interface{}{
							"type":        "array",
							"description": "Arguments for the pub subcommand",
							"items":       map[string]interface{}{"type": "string"},
						},
						"workingDir": map[string]interface{}{
							"type":        "string",
							"description": "Working directory for the command",
						},
					},
					"required": []string{"command"},
				},
			},
			{
				Name:        "dart-compile",
				Description: "Compile Dart code to various formats",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{
							"type":        "string",
							"description": "Path to the Dart file to compile",
						},
						"output": map[string]interface{}{
							"type":        "string",
							"description": "Output file path",
						},
						"format": map[string]interface{}{
							"type":        "string",
							"description": "Output format for the compilation",
							"enum":        []string{"exe", "aot-snapshot", "jit-snapshot", "kernel", "js"},
						},
						"options": map[string]interface{}{
							"type":        "array",
							"description": "Additional compilation options",
							"items":       map[string]interface{}{"type": "string"},
						},
					},
					"required": []string{"path", "output", "format"},
				},
			},
		},
	}
}
