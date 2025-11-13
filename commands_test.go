package tasmota

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_ExecuteCommand(t *testing.T) {
	tests := []struct {
		name       string
		command    string
		response   string
		statusCode int
		wantErr    bool
		errType    ErrorType
	}{
		{
			name:       "successful command",
			command:    "Power",
			response:   `{"POWER":"ON"}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "status command",
			command:    "Status 0",
			response:   `{"Status":{"Module":1,"DeviceName":"Tasmota"}}`,
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name:    "empty command",
			command: "",
			wantErr: true,
			errType: ErrorTypeCommand,
		},
		{
			name:       "invalid JSON response",
			command:    "Power",
			response:   `not json`,
			statusCode: http.StatusOK,
			wantErr:    true,
			errType:    ErrorTypeParse,
		},
		{
			name:       "server error",
			command:    "Power",
			response:   "",
			statusCode: http.StatusInternalServerError,
			wantErr:    true,
			errType:    ErrorTypeNetwork,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify command is in query params
				if tt.command != "" && !strings.Contains(r.URL.RawQuery, "cmnd=") {
					t.Error("command not found in query params")
				}

				w.WriteHeader(tt.statusCode)
				if tt.response != "" {
					_, _ = w.Write([]byte(tt.response))
				}
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			resp, err := client.ExecuteCommand(context.Background(), tt.command)
			if tt.wantErr {
				if err == nil {
					t.Error("ExecuteCommand() expected error, got nil")
					return
				}

				// Check error type
				switch tt.errType {
				case ErrorTypeCommand:
					if !IsCommandError(err) {
						t.Errorf("error type = %T, want command error", err)
					}
				case ErrorTypeParse:
					if !IsParseError(err) {
						t.Errorf("error type = %T, want parse error", err)
					}
				case ErrorTypeNetwork:
					if !IsNetworkError(err) {
						t.Errorf("error type = %T, want network error", err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("ExecuteCommand() unexpected error: %v", err)
					return
				}

				// Verify response is valid JSON
				var v interface{}
				if err := json.Unmarshal(resp, &v); err != nil {
					t.Errorf("response is not valid JSON: %v", err)
				}

				// Verify response matches expected
				if string(resp) != tt.response {
					t.Errorf("response = %v, want %v", string(resp), tt.response)
				}
			}
		})
	}
}

func TestClient_ExecuteBacklog(t *testing.T) {
	tests := []struct {
		name     string
		commands []string
		wantErr  bool
		errType  ErrorType
		wantCmd  string
	}{
		{
			name:     "single command",
			commands: []string{"Power ON"},
			wantErr:  false,
			wantCmd:  "Backlog Power ON",
		},
		{
			name:     "multiple commands",
			commands: []string{"Power1 ON", "Power2 OFF", "Delay 10"},
			wantErr:  false,
			wantCmd:  "Backlog Power1 ON; Power2 OFF; Delay 10",
		},
		{
			name:     "commands with spaces",
			commands: []string{"  Power ON  ", "Status 0"},
			wantErr:  false,
			wantCmd:  "Backlog Power ON; Status 0",
		},
		{
			name:     "empty commands filtered",
			commands: []string{"Power ON", "", "Status 0", "  "},
			wantErr:  false,
			wantCmd:  "Backlog Power ON; Status 0",
		},
		{
			name:     "no commands",
			commands: []string{},
			wantErr:  true,
			errType:  ErrorTypeCommand,
		},
		{
			name:     "all empty commands",
			commands: []string{"", "  ", ""},
			wantErr:  true,
			errType:  ErrorTypeCommand,
		},
		{
			name:     "too many commands",
			commands: make([]string, 31),
			wantErr:  true,
			errType:  ErrorTypeCommand,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize slice commands for too many commands test
			if tt.name == "too many commands" {
				for i := range tt.commands {
					tt.commands[i] = "Power"
				}
			}

			var receivedCommand string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				receivedCommand = r.URL.Query().Get("cmnd")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"Response":"Done"}`))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			resp, err := client.ExecuteBacklog(context.Background(), tt.commands...)
			if tt.wantErr {
				if err == nil {
					t.Error("ExecuteBacklog() expected error, got nil")
					return
				}

				// Check error type
				if tt.errType == ErrorTypeCommand && !IsCommandError(err) {
					t.Errorf("error type = %T, want command error", err)
				}
			} else {
				if err != nil {
					t.Errorf("ExecuteBacklog() unexpected error: %v", err)
					return
				}

				if resp == nil {
					t.Error("ExecuteBacklog() returned nil response")
					return
				}

				// Verify the command sent to server
				if receivedCommand != tt.wantCmd {
					t.Errorf("command sent = %v, want %v", receivedCommand, tt.wantCmd)
				}
			}
		})
	}
}

func TestClient_ExecuteBacklog_Integration(t *testing.T) {
	commandsReceived := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cmd := r.URL.Query().Get("cmnd")
		commandsReceived = append(commandsReceived, cmd)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"Response":"Done"}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	// Test multiple backlog executions
	_, err := client.ExecuteBacklog(context.Background(), "Power1 ON", "Delay 5", "Power2 ON")
	if err != nil {
		t.Fatalf("ExecuteBacklog() error: %v", err)
	}

	if len(commandsReceived) != 1 {
		t.Errorf("expected 1 command, got %d", len(commandsReceived))
	}

	expected := "Backlog Power1 ON; Delay 5; Power2 ON"
	if commandsReceived[0] != expected {
		t.Errorf("command = %v, want %v", commandsReceived[0], expected)
	}
}
