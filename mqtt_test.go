package tasmota

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_GetMQTTConfig(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		if callCount == 0 {
			// First call: Status 6 (MQTT info)
			_, _ = w.Write([]byte(`{"StatusMQT":{"MqttHost":"mqtt.local","MqttPort":1883,"MqttUser":"user"}}`))
		} else {
			// Second call: Status 1 (Device info)
			_, _ = w.Write([]byte(`{"Status":{"Topic":"tasmota_test"}}`))
		}
		callCount++
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	config, err := client.GetMQTTConfig(context.Background())
	if err != nil {
		t.Fatalf("GetMQTTConfig() error: %v", err)
	}

	if config.Host != "mqtt.local" {
		t.Errorf("Host = %v, want mqtt.local", config.Host)
	}
	if config.Port != 1883 {
		t.Errorf("Port = %v, want 1883", config.Port)
	}
	if config.User != "user" {
		t.Errorf("User = %v, want user", config.User)
	}
	if config.Topic != "tasmota_test" {
		t.Errorf("Topic = %v, want tasmota_test", config.Topic)
	}
}

func TestClient_SetMQTTHost(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		wantErr bool
	}{
		{"valid host", "mqtt.example.com", false},
		{"valid IP", "192.168.1.100", false},
		{"empty host", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetMQTTHost(context.Background(), tt.host)
				if err == nil {
					t.Error("SetMQTTHost() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				cmd := r.URL.Query().Get("cmnd")
				if !strings.Contains(cmd, tt.host) {
					t.Errorf("command does not contain host %s", tt.host)
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"MqttHost":"` + tt.host + `"}`))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.SetMQTTHost(context.Background(), tt.host)
			if err != nil {
				t.Errorf("SetMQTTHost() error: %v", err)
			}
		})
	}
}

func TestClient_SetMQTTPort(t *testing.T) {
	tests := []struct {
		name    string
		port    int
		wantErr bool
	}{
		{"valid 1883", 1883, false},
		{"valid 8883", 8883, false},
		{"invalid 0", 0, true},
		{"invalid 70000", 70000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetMQTTPort(context.Background(), tt.port)
				if err == nil {
					t.Error("SetMQTTPort() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"MqttPort":1883}`))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.SetMQTTPort(context.Background(), tt.port)
			if err != nil {
				t.Errorf("SetMQTTPort() error: %v", err)
			}
		})
	}
}

func TestClient_SetMQTTUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"MqttUser":"testuser"}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	err := client.SetMQTTUser(context.Background(), "testuser")
	if err != nil {
		t.Errorf("SetMQTTUser() error: %v", err)
	}
}

func TestClient_SetMQTTPassword(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"MqttPassword":"****"}`))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	err := client.SetMQTTPassword(context.Background(), "secret")
	if err != nil {
		t.Errorf("SetMQTTPassword() error: %v", err)
	}
}

func TestClient_SetMQTTClient(t *testing.T) {
	tests := []struct {
		name       string
		clientName string
		wantErr    bool
	}{
		{"valid name", "tasmota-client", false},
		{"empty name", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetMQTTClient(context.Background(), tt.clientName)
				if err == nil {
					t.Error("SetMQTTClient() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"MqttClient":"` + tt.clientName + `"}`))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.SetMQTTClient(context.Background(), tt.clientName)
			if err != nil {
				t.Errorf("SetMQTTClient() error: %v", err)
			}
		})
	}
}

func TestClient_SetTopic(t *testing.T) {
	tests := []struct {
		name    string
		topic   string
		wantErr bool
	}{
		{"valid topic", "living-room", false},
		{"empty topic", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetTopic(context.Background(), tt.topic)
				if err == nil {
					t.Error("SetTopic() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"Topic":"` + tt.topic + `"}`))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.SetTopic(context.Background(), tt.topic)
			if err != nil {
				t.Errorf("SetTopic() error: %v", err)
			}
		})
	}
}

func TestClient_SetFullTopic(t *testing.T) {
	tests := []struct {
		name      string
		fullTopic string
		wantErr   bool
	}{
		{"valid topic", "%prefix%/%topic%/", false},
		{"empty topic", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetFullTopic(context.Background(), tt.fullTopic)
				if err == nil {
					t.Error("SetFullTopic() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"FullTopic":"` + tt.fullTopic + `"}`))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.SetFullTopic(context.Background(), tt.fullTopic)
			if err != nil {
				t.Errorf("SetFullTopic() error: %v", err)
			}
		})
	}
}

func TestClient_SetGroupTopic(t *testing.T) {
	tests := []struct {
		name       string
		groupTopic string
		wantErr    bool
	}{
		{"valid topic", "tasmotas", false},
		{"empty topic", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetGroupTopic(context.Background(), tt.groupTopic)
				if err == nil {
					t.Error("SetGroupTopic() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"GroupTopic":"` + tt.groupTopic + `"}`))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.SetGroupTopic(context.Background(), tt.groupTopic)
			if err != nil {
				t.Errorf("SetGroupTopic() error: %v", err)
			}
		})
	}
}

func TestClient_SetPrefix(t *testing.T) {
	tests := []struct {
		name      string
		prefixNum int
		prefix    string
		wantErr   bool
	}{
		{"valid prefix 1", 1, "cmnd", false},
		{"valid prefix 2", 2, "stat", false},
		{"valid prefix 3", 3, "tele", false},
		{"invalid prefix 0", 0, "test", true},
		{"invalid prefix 4", 4, "test", true},
		{"empty prefix", 1, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetPrefix(context.Background(), tt.prefixNum, tt.prefix)
				if err == nil {
					t.Error("SetPrefix() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"Prefix":"` + tt.prefix + `"}`))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.SetPrefix(context.Background(), tt.prefixNum, tt.prefix)
			if err != nil {
				t.Errorf("SetPrefix() error: %v", err)
			}
		})
	}
}

func TestClient_EnableMQTT(t *testing.T) {
	tests := []struct {
		name   string
		enable bool
	}{
		{"enable", true},
		{"disable", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"SetOption3":0}`))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.EnableMQTT(context.Background(), tt.enable)
			if err != nil {
				t.Errorf("EnableMQTT() error: %v", err)
			}
		})
	}
}

func TestClient_SetMQTTConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *MQTTConfig
		wantErr bool
	}{
		{
			name: "full config",
			config: &MQTTConfig{
				Host:       "mqtt.local",
				Port:       1883,
				User:       "user",
				Password:   "pass",
				Topic:      "tasmota",
				FullTopic:  "%prefix%/%topic%/",
				GroupTopic: "tasmotas",
				Retain:     true,
				TelePeriod: 300,
			},
			wantErr: false,
		},
		{
			name: "partial config",
			config: &MQTTConfig{
				Host: "mqtt.local",
				Port: 1883,
			},
			wantErr: false,
		},
		{
			name:    "nil config",
			config:  nil,
			wantErr: true,
		},
		{
			name:    "empty config",
			config:  &MQTTConfig{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetMQTTConfig(context.Background(), tt.config)
				if err == nil {
					t.Error("SetMQTTConfig() expected error, got nil")
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

			err := client.SetMQTTConfig(context.Background(), tt.config)
			if err != nil {
				t.Errorf("SetMQTTConfig() error: %v", err)
			}
		})
	}
}

func TestClient_MQTTFingerprint(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		mockResponse := `{"MqttFingerprint":"AA BB CC DD EE FF"}`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(mockResponse))
		}))
		defer server.Close()

		client := &Client{
			baseURL:    server.URL,
			httpClient: server.Client(),
		}

		fp, err := client.GetMQTTFingerprint(context.Background())
		if err != nil {
			t.Fatalf("GetMQTTFingerprint() error: %v", err)
		}
		if fp != "AA BB CC DD EE FF" {
			t.Errorf("GetMQTTFingerprint() = %v, want AA BB CC DD EE FF", fp)
		}
	})

	t.Run("set valid", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"MqttFingerprint":"AA BB CC DD EE FF"}`))
		}))
		defer server.Close()

		client := &Client{
			baseURL:    server.URL,
			httpClient: server.Client(),
		}

		err := client.SetMQTTFingerprint(context.Background(), "AA BB CC DD EE FF")
		if err != nil {
			t.Errorf("SetMQTTFingerprint() error: %v", err)
		}
	})

	t.Run("set empty", func(t *testing.T) {
		client := &Client{}
		err := client.SetMQTTFingerprint(context.Background(), "")
		if err == nil {
			t.Error("SetMQTTFingerprint() expected error, got nil")
		}
	})
}

func TestClient_MQTTRetry(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		mockResponse := `{"MqttRetry":30}`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(mockResponse))
		}))
		defer server.Close()

		client := &Client{
			baseURL:    server.URL,
			httpClient: server.Client(),
		}

		retry, err := client.GetMQTTRetry(context.Background())
		if err != nil {
			t.Fatalf("GetMQTTRetry() error: %v", err)
		}
		if retry != 30 {
			t.Errorf("GetMQTTRetry() = %v, want 30", retry)
		}
	})

	t.Run("set valid", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"MqttRetry":60}`))
		}))
		defer server.Close()

		client := &Client{
			baseURL:    server.URL,
			httpClient: server.Client(),
		}

		err := client.SetMQTTRetry(context.Background(), 60)
		if err != nil {
			t.Errorf("SetMQTTRetry() error: %v", err)
		}
	})

	t.Run("set invalid", func(t *testing.T) {
		client := &Client{}

		tests := []int{5, 50000}
		for _, val := range tests {
			err := client.SetMQTTRetry(context.Background(), val)
			if err == nil {
				t.Errorf("SetMQTTRetry(%d) expected error, got nil", val)
			}
		}
	})
}

func TestClient_TestMQTTConnection(t *testing.T) {
	tests := []struct {
		name         string
		mockResponse string
		wantErr      bool
	}{
		{
			name:         "connected",
			mockResponse: `{"StatusMQT":{"MqttHost":"mqtt.local","MqttCount":5}}`,
			wantErr:      false,
		},
		{
			name:         "not connected",
			mockResponse: `{"StatusMQT":{"MqttHost":"mqtt.local","MqttCount":0}}`,
			wantErr:      true,
		},
		{
			name:         "not configured",
			mockResponse: `{"StatusMQT":{"MqttHost":""}}`,
			wantErr:      true,
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

			err := client.TestMQTTConnection(context.Background())
			if tt.wantErr {
				if err == nil {
					t.Error("TestMQTTConnection() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("TestMQTTConnection() unexpected error: %v", err)
				}
			}
		})
	}
}
