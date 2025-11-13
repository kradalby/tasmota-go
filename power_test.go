package tasmota

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPowerResponse_IsOn(t *testing.T) {
	tests := []struct {
		name     string
		response PowerResponse
		relayNum int
		expected bool
	}{
		{
			name:     "POWER is ON",
			response: PowerResponse{Power: "ON"},
			relayNum: 0,
			expected: true,
		},
		{
			name:     "POWER is OFF",
			response: PowerResponse{Power: "OFF"},
			relayNum: 0,
			expected: false,
		},
		{
			name:     "POWER1 is ON",
			response: PowerResponse{Power1: "ON"},
			relayNum: 1,
			expected: true,
		},
		{
			name:     "POWER2 is OFF",
			response: PowerResponse{Power2: "OFF"},
			relayNum: 2,
			expected: false,
		},
		{
			name:     "lowercase on",
			response: PowerResponse{Power: "on"},
			relayNum: 0,
			expected: true,
		},
		{
			name:     "invalid relay number",
			response: PowerResponse{Power: "ON"},
			relayNum: 9,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.response.IsOn(tt.relayNum); got != tt.expected {
				t.Errorf("IsOn(%d) = %v, want %v", tt.relayNum, got, tt.expected)
			}
		})
	}
}

func TestPowerResponse_GetState(t *testing.T) {
	resp := PowerResponse{
		Power:  "ON",
		Power1: "OFF",
		Power2: "ON",
		Power3: "OFF",
		Power4: "ON",
		Power5: "OFF",
		Power6: "ON",
		Power7: "OFF",
		Power8: "ON",
	}

	tests := []struct {
		relayNum int
		expected string
	}{
		{0, "ON"},
		{1, "OFF"},
		{2, "ON"},
		{3, "OFF"},
		{4, "ON"},
		{5, "OFF"},
		{6, "ON"},
		{7, "OFF"},
		{8, "ON"},
		{9, ""}, // Invalid
	}

	for _, tt := range tests {
		t.Run(string(rune('0'+tt.relayNum)), func(t *testing.T) {
			if got := resp.GetState(tt.relayNum); got != tt.expected {
				t.Errorf("GetState(%d) = %v, want %v", tt.relayNum, got, tt.expected)
			}
		})
	}
}

func TestClient_Power(t *testing.T) {
	tests := []struct {
		name         string
		state        PowerState
		mockResponse string
		wantErr      bool
		wantPower    string
	}{
		{
			name:         "turn on",
			state:        PowerOn,
			mockResponse: `{"POWER":"ON"}`,
			wantErr:      false,
			wantPower:    "ON",
		},
		{
			name:         "turn off",
			state:        PowerOff,
			mockResponse: `{"POWER":"OFF"}`,
			wantErr:      false,
			wantPower:    "OFF",
		},
		{
			name:         "toggle",
			state:        PowerToggle,
			mockResponse: `{"POWER":"ON"}`,
			wantErr:      false,
			wantPower:    "ON",
		},
		{
			name:         "blink",
			state:        PowerBlink,
			mockResponse: `{"POWER":"ON"}`,
			wantErr:      false,
			wantPower:    "ON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify command contains the state
				if !strings.Contains(r.URL.Query().Get("cmnd"), string(tt.state)) {
					t.Errorf("command does not contain state %s", tt.state)
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			resp, err := client.Power(context.Background(), tt.state)
			if tt.wantErr {
				if err == nil {
					t.Error("Power() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Power() unexpected error: %v", err)
					return
				}
				if resp.Power != tt.wantPower {
					t.Errorf("Power() = %v, want %v", resp.Power, tt.wantPower)
				}
			}
		})
	}
}

func TestClient_PowerN(t *testing.T) {
	tests := []struct {
		name         string
		relayNum     int
		state        PowerState
		mockResponse string
		wantErr      bool
		wantField    string
		wantValue    string
	}{
		{
			name:         "relay 1 on",
			relayNum:     1,
			state:        PowerOn,
			mockResponse: `{"POWER1":"ON"}`,
			wantErr:      false,
			wantField:    "POWER1",
			wantValue:    "ON",
		},
		{
			name:         "relay 3 off",
			relayNum:     3,
			state:        PowerOff,
			mockResponse: `{"POWER3":"OFF"}`,
			wantErr:      false,
			wantField:    "POWER3",
			wantValue:    "OFF",
		},
		{
			name:         "relay 8 toggle",
			relayNum:     8,
			state:        PowerToggle,
			mockResponse: `{"POWER8":"ON"}`,
			wantErr:      false,
			wantField:    "POWER8",
			wantValue:    "ON",
		},
		{
			name:     "invalid relay 0",
			relayNum: 0,
			state:    PowerOn,
			wantErr:  true,
		},
		{
			name:     "invalid relay 9",
			relayNum: 9,
			state:    PowerOn,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr && tt.mockResponse == "" {
				// Test validation errors without server
				client := &Client{}
				_, err := client.PowerN(context.Background(), tt.relayNum, tt.state)
				if err == nil {
					t.Error("PowerN() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				cmd := r.URL.Query().Get("cmnd")
				expectedCmd := strings.ToLower("Power" + string(rune('0'+tt.relayNum)))
				if !strings.Contains(strings.ToLower(cmd), expectedCmd) {
					t.Errorf("command = %s, want to contain %s", cmd, expectedCmd)
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			resp, err := client.PowerN(context.Background(), tt.relayNum, tt.state)
			if err != nil {
				t.Errorf("PowerN() unexpected error: %v", err)
				return
			}

			state := resp.GetState(tt.relayNum)
			if state != tt.wantValue {
				t.Errorf("PowerN() relay %d state = %v, want %v", tt.relayNum, state, tt.wantValue)
			}
		})
	}
}

func TestClient_GetPower(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cmd := r.URL.Query().Get("cmnd")
		if cmd != "Power" {
			t.Errorf("command = %s, want Power", cmd)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"POWER":"ON"}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	resp, err := client.GetPower(context.Background())
	if err != nil {
		t.Fatalf("GetPower() error: %v", err)
	}
	if resp.Power != "ON" {
		t.Errorf("GetPower() = %v, want ON", resp.Power)
	}
}

func TestClient_GetPowerN(t *testing.T) {
	tests := []struct {
		name     string
		relayNum int
		wantErr  bool
	}{
		{"relay 1", 1, false},
		{"relay 8", 8, false},
		{"invalid 0", 0, true},
		{"invalid 9", 9, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				_, err := client.GetPowerN(context.Background(), tt.relayNum)
				if err == nil {
					t.Error("GetPowerN() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				resp := `{"POWER` + string(rune('0'+tt.relayNum)) + `":"ON"}`
				_, _ = w.Write([]byte(resp))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			resp, err := client.GetPowerN(context.Background(), tt.relayNum)
			if err != nil {
				t.Errorf("GetPowerN() error: %v", err)
				return
			}
			if resp.GetState(tt.relayNum) != "ON" {
				t.Errorf("GetPowerN() relay %d state = %v, want ON", tt.relayNum, resp.GetState(tt.relayNum))
			}
		})
	}
}

func TestClient_IsPowerOn(t *testing.T) {
	tests := []struct {
		name         string
		relayNum     int
		mockResponse string
		expected     bool
	}{
		{
			name:         "relay 0 on",
			relayNum:     0,
			mockResponse: `{"POWER":"ON"}`,
			expected:     true,
		},
		{
			name:         "relay 0 off",
			relayNum:     0,
			mockResponse: `{"POWER":"OFF"}`,
			expected:     false,
		},
		{
			name:         "relay 2 on",
			relayNum:     2,
			mockResponse: `{"POWER2":"ON"}`,
			expected:     true,
		},
		{
			name:         "relay 2 off",
			relayNum:     2,
			mockResponse: `{"POWER2":"OFF"}`,
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			isOn, err := client.IsPowerOn(context.Background(), tt.relayNum)
			if err != nil {
				t.Fatalf("IsPowerOn() error: %v", err)
			}
			if isOn != tt.expected {
				t.Errorf("IsPowerOn(%d) = %v, want %v", tt.relayNum, isOn, tt.expected)
			}
		})
	}
}

func TestClient_SetPowerOn(t *testing.T) {
	tests := []struct {
		name     string
		relayNum int
	}{
		{"relay 0", 0},
		{"relay 1", 1},
		{"relay 5", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				cmd := r.URL.Query().Get("cmnd")
				if !strings.Contains(cmd, "ON") {
					t.Error("command does not contain ON")
				}
				w.WriteHeader(http.StatusOK)
				if tt.relayNum == 0 {
					_, _ = w.Write([]byte(`{"POWER":"ON"}`))
				} else {
					resp := `{"POWER` + string(rune('0'+tt.relayNum)) + `":"ON"}`
					_, _ = w.Write([]byte(resp))
				}
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.SetPowerOn(context.Background(), tt.relayNum)
			if err != nil {
				t.Errorf("SetPowerOn() error: %v", err)
			}
		})
	}
}

func TestClient_SetPowerOff(t *testing.T) {
	tests := []struct {
		name     string
		relayNum int
	}{
		{"relay 0", 0},
		{"relay 1", 1},
		{"relay 5", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				cmd := r.URL.Query().Get("cmnd")
				if !strings.Contains(cmd, "OFF") {
					t.Error("command does not contain OFF")
				}
				w.WriteHeader(http.StatusOK)
				if tt.relayNum == 0 {
					_, _ = w.Write([]byte(`{"POWER":"OFF"}`))
				} else {
					resp := `{"POWER` + string(rune('0'+tt.relayNum)) + `":"OFF"}`
					_, _ = w.Write([]byte(resp))
				}
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.SetPowerOff(context.Background(), tt.relayNum)
			if err != nil {
				t.Errorf("SetPowerOff() error: %v", err)
			}
		})
	}
}

func TestClient_TogglePower(t *testing.T) {
	tests := []struct {
		name     string
		relayNum int
	}{
		{"relay 0", 0},
		{"relay 3", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				cmd := r.URL.Query().Get("cmnd")
				if !strings.Contains(cmd, "TOGGLE") {
					t.Error("command does not contain TOGGLE")
				}
				w.WriteHeader(http.StatusOK)
				if tt.relayNum == 0 {
					_, _ = w.Write([]byte(`{"POWER":"ON"}`))
				} else {
					resp := `{"POWER` + string(rune('0'+tt.relayNum)) + `":"ON"}`
					_, _ = w.Write([]byte(resp))
				}
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.TogglePower(context.Background(), tt.relayNum)
			if err != nil {
				t.Errorf("TogglePower() error: %v", err)
			}
		})
	}
}

func TestClient_GetCurrentPower(t *testing.T) {
	tests := []struct {
		name         string
		mockResponse string
		expected     float64
		wantErr      bool
	}{
		{
			name:         "float power value",
			mockResponse: `{"StatusSNS":{"ENERGY":{"Power":123.45}}}`,
			expected:     123.45,
			wantErr:      false,
		},
		{
			name:         "int power value",
			mockResponse: `{"StatusSNS":{"ENERGY":{"Power":100}}}`,
			expected:     100.0,
			wantErr:      false,
		},
		{
			name:         "string power value",
			mockResponse: `{"StatusSNS":{"ENERGY":{"Power":"75.5"}}}`,
			expected:     75.5,
			wantErr:      false,
		},
		{
			name:         "invalid JSON",
			mockResponse: `invalid`,
			wantErr:      true,
		},
		{
			name:         "missing power field",
			mockResponse: `{"StatusSNS":{"ENERGY":{}}}`,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				cmd := r.URL.Query().Get("cmnd")
				if cmd != "Status 10" {
					t.Errorf("command = %s, want Status 10", cmd)
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			power, err := client.GetCurrentPower(context.Background())
			if tt.wantErr {
				if err == nil {
					t.Error("GetCurrentPower() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("GetCurrentPower() unexpected error: %v", err)
					return
				}
				if power != tt.expected {
					t.Errorf("GetCurrentPower() = %v, want %v", power, tt.expected)
				}
			}
		})
	}
}
