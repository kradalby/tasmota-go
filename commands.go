package tasmota

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// ExecuteCommand sends a command to the Tasmota device and returns the raw JSON response.
func (c *Client) ExecuteCommand(ctx context.Context, command string) (json.RawMessage, error) {
	if command == "" {
		return nil, NewError(ErrorTypeCommand, "command cannot be empty", nil)
	}

	urlStr, err := c.buildURL(command)
	if err != nil {
		return nil, err
	}

	body, err := c.do(ctx, urlStr)
	if err != nil {
		return nil, err
	}

	// Validate that the response is JSON
	var raw json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, NewError(ErrorTypeParse, "invalid JSON response", err)
	}

	return raw, nil
}

// ExecuteBacklog executes multiple commands in sequence using the Backlog command.
// Tasmota will execute each command sequentially with a 1-second delay between them.
// The commands are separated by semicolons.
func (c *Client) ExecuteBacklog(ctx context.Context, commands ...string) (json.RawMessage, error) {
	if len(commands) == 0 {
		return nil, NewError(ErrorTypeCommand, "no commands provided", nil)
	}

	if len(commands) > 30 {
		return nil, NewError(ErrorTypeCommand, "backlog supports maximum 30 commands", nil)
	}

	// Filter out empty commands
	var validCommands []string
	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)
		if cmd != "" {
			validCommands = append(validCommands, cmd)
		}
	}

	if len(validCommands) == 0 {
		return nil, NewError(ErrorTypeCommand, "no valid commands provided", nil)
	}

	// Build the backlog command
	backlogCmd := fmt.Sprintf("Backlog %s", strings.Join(validCommands, "; "))

	return c.ExecuteCommand(ctx, backlogCmd)
}
