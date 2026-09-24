package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Built-in tool implementations
func ReadFile(input json.RawMessage) (string, error) {
	var params ReadFileInput
	err := json.Unmarshal(input, &params)
	if err != nil {
		return "", fmt.Errorf("invalid parameters: %w", err)
	}

	content, err := os.ReadFile(params.Path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), nil
}

func ListFiles(input json.RawMessage) (string, error) {
	var params ListFilesInput
	err := json.Unmarshal(input, &params)
	if err != nil {
		return "", fmt.Errorf("invalid parameters: %w", err)
	}

	if params.Path == "" {
		params.Path = "."
	}

	entries, err := os.ReadDir(params.Path)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	var result []string
	for _, entry := range entries {
		if entry.IsDir() {
			result = append(result, entry.Name()+"/")
		} else {
			result = append(result, entry.Name())
		}
	}

	return strings.Join(result, "\n"), nil
}

func EditFile(input json.RawMessage) (string, error) {
	var params EditFileInput
	err := json.Unmarshal(input, &params)
	if err != nil {
		return "", fmt.Errorf("invalid parameters: %w", err)
	}

	// Read current content
	content, err := os.ReadFile(params.Path)
	if err != nil {
		// If file doesn't exist, create it
		if os.IsNotExist(err) {
			err = os.WriteFile(params.Path, []byte(params.NewStr), 0644)
			if err != nil {
				return "", fmt.Errorf("failed to create file: %w", err)
			}
			return fmt.Sprintf("Created new file: %s", params.Path), nil
		}
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Replace content
	newContent := strings.ReplaceAll(string(content), params.OldStr, params.NewStr)
	
	// Write back to file
	err = os.WriteFile(params.Path, []byte(newContent), 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return fmt.Sprintf("Successfully edited file: %s", params.Path), nil
}

func BashCommand(input json.RawMessage) (string, error) {
	var params BashInput
	err := json.Unmarshal(input, &params)
	if err != nil {
		return "", fmt.Errorf("invalid parameters: %w", err)
	}

	cmd := exec.Command("bash", "-c", params.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Command failed: %s\nOutput: %s", err.Error(), string(output)), nil
	}

	return string(output), nil
}

func CodeSearch(input json.RawMessage) (string, error) {
	var params CodeSearchInput
	err := json.Unmarshal(input, &params)
	if err != nil {
		return "", fmt.Errorf("invalid parameters: %w", err)
	}

	if params.Path == "" {
		params.Path = "."
	}

	// Build grep command
	args := []string{"-r", "-n"}
	
	if !params.CaseSensitive {
		args = append(args, "-i")
	}
	
	if params.FileType != "" {
		args = append(args, "--include=*."+params.FileType)
	}
	
	args = append(args, params.Pattern, params.Path)

	cmd := exec.Command("grep", args...)
	output, err := cmd.Output()
	if err != nil {
		// grep returns exit code 1 when no matches found
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 {
			return "No matches found", nil
		}
		return "", fmt.Errorf("grep command failed: %w", err)
	}

	return string(output), nil
}
