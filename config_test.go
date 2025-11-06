package tasmota

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_GetConfig(t *testing.T) {
	mockResponse := `{
		"Status": {
			"Module": 1,
			"DeviceName": "Tasmota",
			"FriendlyName": ["Living Room", "Kitchen"],
			"PowerOnState": 3,
			"LedState": 1,
			"ButtonRetain": 0,
			"SwitchRetain": 0,
			"SensorRetain": 0,
			"PowerRetain": 0
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	config, err := client.GetConfig(context.Background())
	if err != nil {
		t.Fatalf("GetConfig() error: %v", err)
	}

	if config.DeviceName != "Tasmota" {
		t.Errorf("DeviceName = %v, want Tasmota", config.DeviceName)
	}
	if len(config.FriendlyName) != 2 {
		t.Errorf("FriendlyName length = %v, want 2", len(config.FriendlyName))
	}
	if config.PowerOnState != 3 {
		t.Errorf("PowerOnState = %v, want 3", config.PowerOnState)
	}
}

func TestClient_SetDeviceName(t *testing.T) {
	tests := []struct {
		name    string
		devName string
		wantErr bool
	}{
		{"valid name", "MyDevice", false},
		{"empty name", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetDeviceName(context.Background(), tt.devName)
				if err == nil {
					t.Error("SetDeviceName() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				cmd := r.URL.Query().Get("cmnd")
				if !strings.Contains(cmd, tt.devName) {
					t.Errorf("command does not contain device name %s", tt.devName)
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"DeviceName":"` + tt.devName + `"}`))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.SetDeviceName(context.Background(), tt.devName)
			if err != nil {
				t.Errorf("SetDeviceName() error: %v", err)
			}
		})
	}
}

func TestClient_SetFriendlyName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cmd := r.URL.Query().Get("cmnd")
		if !strings.Contains(cmd, "FriendlyName1") {
			t.Error("command should contain FriendlyName1")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"FriendlyName1":"Test"}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	err := client.SetFriendlyName(context.Background(), "Test")
	if err != nil {
		t.Errorf("SetFriendlyName() error: %v", err)
	}
}

func TestClient_SetFriendlyNameN(t *testing.T) {
	tests := []struct {
		name    string
		index   int
		fname   string
		wantErr bool
	}{
		{"valid index 1", 1, "Name1", false},
		{"valid index 8", 8, "Name8", false},
		{"invalid index 0", 0, "Name", true},
		{"invalid index 9", 9, "Name", true},
		{"empty name", 1, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetFriendlyNameN(context.Background(), tt.index, tt.fname)
				if err == nil {
					t.Error("SetFriendlyNameN() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"FriendlyName":"` + tt.fname + `"}`))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.SetFriendlyNameN(context.Background(), tt.index, tt.fname)
			if err != nil {
				t.Errorf("SetFriendlyNameN() error: %v", err)
			}
		})
	}
}

func TestClient_SetPowerOnState(t *testing.T) {
	tests := []struct {
		name    string
		state   int
		wantErr bool
	}{
		{"valid 0", 0, false},
		{"valid 3", 3, false},
		{"valid 5", 5, false},
		{"invalid -1", -1, true},
		{"invalid 6", 6, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetPowerOnState(context.Background(), tt.state)
				if err == nil {
					t.Error("SetPowerOnState() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, `{"PowerOnState":%d}`, tt.state)
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.SetPowerOnState(context.Background(), tt.state)
			if err != nil {
				t.Errorf("SetPowerOnState() error: %v", err)
			}
		})
	}
}

func TestClient_SetLedState(t *testing.T) {
	tests := []struct {
		name    string
		state   int
		wantErr bool
	}{
		{"valid 0", 0, false},
		{"valid 4", 4, false},
		{"valid 8", 8, false},
		{"invalid -1", -1, true},
		{"invalid 9", 9, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetLedState(context.Background(), tt.state)
				if err == nil {
					t.Error("SetLedState() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, `{"LedState":%d}`, tt.state)
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.SetLedState(context.Background(), tt.state)
			if err != nil {
				t.Errorf("SetLedState() error: %v", err)
			}
		})
	}
}

func TestClient_SetSleep(t *testing.T) {
	tests := []struct {
		name     string
		duration int
		wantErr  bool
	}{
		{"valid 0", 0, false},
		{"valid 50", 50, false},
		{"valid 250", 250, false},
		{"invalid -1", -1, true},
		{"invalid 251", 251, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetSleep(context.Background(), tt.duration)
				if err == nil {
					t.Error("SetSleep() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, `{"Sleep":%d}`, tt.duration)
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.SetSleep(context.Background(), tt.duration)
			if err != nil {
				t.Errorf("SetSleep() error: %v", err)
			}
		})
	}
}

func TestClient_SetRetainFlags(t *testing.T) {
	tests := []struct {
		name   string
		method func(*Client, context.Context, bool) error
		retain bool
	}{
		{"ButtonRetain true", (*Client).SetButtonRetain, true},
		{"ButtonRetain false", (*Client).SetButtonRetain, false},
		{"SwitchRetain true", (*Client).SetSwitchRetain, true},
		{"SwitchRetain false", (*Client).SetSwitchRetain, false},
		{"SensorRetain true", (*Client).SetSensorRetain, true},
		{"SensorRetain false", (*Client).SetSensorRetain, false},
		{"PowerRetain true", (*Client).SetPowerRetain, true},
		{"PowerRetain false", (*Client).SetPowerRetain, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"Retain":1}`))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := tt.method(client, context.Background(), tt.retain)
			if err != nil {
				t.Errorf("%s error: %v", tt.name, err)
			}
		})
	}
}

func TestClient_ApplyConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *DeviceConfig
		wantErr bool
	}{
		{
			name: "full config",
			config: &DeviceConfig{
				DeviceName:   "TestDevice",
				FriendlyName: []string{"Room1", "Room2"},
				PowerOnState: 3,
				LedState:     1,
				Sleep:        50,
				ButtonRetain: 0,
				SwitchRetain: 0,
				SensorRetain: 0,
				PowerRetain:  0,
			},
			wantErr: false,
		},
		{
			name: "partial config",
			config: &DeviceConfig{
				DeviceName:   "TestDevice",
				PowerOnState: 1,
			},
			wantErr: false,
		},
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.ApplyConfig(context.Background(), tt.config)
				if err == nil {
					t.Error("ApplyConfig() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				cmd := r.URL.Query().Get("cmnd")
				if !strings.Contains(cmd, "Backlog") {
					t.Error("command should use Backlog")
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"Response":"Done"}`))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.ApplyConfig(context.Background(), tt.config)
			if err != nil {
				t.Errorf("ApplyConfig() error: %v", err)
			}
		})
	}
}

func TestClient_Restart(t *testing.T) {
	tests := []struct {
		name    string
		reason  int
		wantErr bool
	}{
		{"normal restart", 1, false},
		{"reset defaults", 99, false},
		{"invalid reason", 2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.Restart(context.Background(), tt.reason)
				if err == nil {
					t.Error("Restart() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"Restart":"Restarting"}`))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.Restart(context.Background(), tt.reason)
			if err != nil {
				t.Errorf("Restart() error: %v", err)
			}
		})
	}
}

func TestClient_Reset(t *testing.T) {
	tests := []struct {
		name    string
		level   int
		wantErr bool
	}{
		{"reset relay", 1, false},
		{"reset except wifi", 2, false},
		{"reset except wifi mqtt", 3, false},
		{"reset wifi", 4, false},
		{"erase flash keep wifi", 5, false},
		{"erase all", 6, false},
		{"reset and reboot", 99, false},
		{"invalid level", 7, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.Reset(context.Background(), tt.level)
				if err == nil {
					t.Error("Reset() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"Reset":"Done"}`))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.Reset(context.Background(), tt.level)
			if err != nil {
				t.Errorf("Reset() error: %v", err)
			}
		})
	}
}

func TestClient_GetModule(t *testing.T) {
	mockResponse := `{
		"Status": {
			"Module": 42
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	module, err := client.GetModule(context.Background())
	if err != nil {
		t.Fatalf("GetModule() error: %v", err)
	}
	if module != 42 {
		t.Errorf("GetModule() = %v, want 42", module)
	}
}

func TestClient_SetOption(t *testing.T) {
	tests := []struct {
		name    string
		option  int
		value   interface{}
		wantErr bool
	}{
		{"bool true", 1, true, false},
		{"bool false", 2, false, false},
		{"int", 3, 10, false},
		{"string", 4, "test", false},
		{"negative option", -1, 1, true},
		{"unsupported type", 1, 3.14, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetOption(context.Background(), tt.option, tt.value)
				if err == nil {
					t.Error("SetOption() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"SetOption":"Done"}`))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.SetOption(context.Background(), tt.option, tt.value)
			if err != nil {
				t.Errorf("SetOption() error: %v", err)
			}
		})
	}
}

func TestClient_TelePeriod(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		mockResponse := `{"TelePeriod":300}`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(mockResponse))
		}))
		defer server.Close()

		client := &Client{
			baseURL:    server.URL,
			httpClient: server.Client(),
		}

		period, err := client.GetTelePeriod(context.Background())
		if err != nil {
			t.Fatalf("GetTelePeriod() error: %v", err)
		}
		if period != 300 {
			t.Errorf("GetTelePeriod() = %v, want 300", period)
		}
	})

	t.Run("set valid", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"TelePeriod":600}`))
		}))
		defer server.Close()

		client := &Client{
			baseURL:    server.URL,
			httpClient: server.Client(),
		}

		err := client.SetTelePeriod(context.Background(), 600)
		if err != nil {
			t.Errorf("SetTelePeriod() error: %v", err)
		}
	})

	t.Run("set invalid", func(t *testing.T) {
		client := &Client{}

		tests := []int{5, 5000}
		for _, val := range tests {
			err := client.SetTelePeriod(context.Background(), val)
			if err == nil {
				t.Errorf("SetTelePeriod(%d) expected error, got nil", val)
			}
		}
	})
}

func TestClient_SetTemplate(t *testing.T) {
	tests := []struct {
		name     string
		template string
		wantErr  bool
	}{
		{"valid template", `{"NAME":"Test"}`, false},
		{"empty template", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetTemplate(context.Background(), tt.template)
				if err == nil {
					t.Error("SetTemplate() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"Template":"Done"}`))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.SetTemplate(context.Background(), tt.template)
			if err != nil {
				t.Errorf("SetTemplate() error: %v", err)
			}
		})
	}
}
