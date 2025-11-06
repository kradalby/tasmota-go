package tasmota

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_GetNetworkConfig(t *testing.T) {
	mockResponse := `{
		"StatusNET": {
			"Hostname": "tasmota-test",
			"IPAddress": "192.168.1.100",
			"Gateway": "192.168.1.1",
			"Subnetmask": "255.255.255.0",
			"DNSServer": "192.168.1.1",
			"Mac": "AA:BB:CC:DD:EE:FF"
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

	config, err := client.GetNetworkConfig(context.Background())
	if err != nil {
		t.Fatalf("GetNetworkConfig() error: %v", err)
	}

	if config.Hostname != "tasmota-test" {
		t.Errorf("Hostname = %v, want tasmota-test", config.Hostname)
	}
	if config.IPAddress != "192.168.1.100" {
		t.Errorf("IPAddress = %v, want 192.168.1.100", config.IPAddress)
	}
	if config.Gateway != "192.168.1.1" {
		t.Errorf("Gateway = %v, want 192.168.1.1", config.Gateway)
	}
}

func TestClient_SetHostname(t *testing.T) {
	tests := []struct {
		name     string
		hostname string
		wantErr  bool
	}{
		{"valid hostname", "tasmota-device", false},
		{"empty hostname", "", true},
		{"too long", strings.Repeat("a", 33), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetHostname(context.Background(), tt.hostname)
				if err == nil {
					t.Error("SetHostname() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				cmd := r.URL.Query().Get("cmnd")
				if !strings.Contains(cmd, tt.hostname) {
					t.Errorf("command does not contain hostname %s", tt.hostname)
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"Hostname":"` + tt.hostname + `"}`))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.SetHostname(context.Background(), tt.hostname)
			if err != nil {
				t.Errorf("SetHostname() error: %v", err)
			}
		})
	}
}

func TestClient_SetStaticIP(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		gateway string
		subnet  string
		wantErr bool
	}{
		{
			name:    "valid config",
			ip:      "192.168.1.100",
			gateway: "192.168.1.1",
			subnet:  "255.255.255.0",
			wantErr: false,
		},
		{
			name:    "invalid IP",
			ip:      "invalid",
			gateway: "192.168.1.1",
			subnet:  "255.255.255.0",
			wantErr: true,
		},
		{
			name:    "invalid gateway",
			ip:      "192.168.1.100",
			gateway: "invalid",
			subnet:  "255.255.255.0",
			wantErr: true,
		},
		{
			name:    "invalid subnet",
			ip:      "192.168.1.100",
			gateway: "192.168.1.1",
			subnet:  "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetStaticIP(context.Background(), tt.ip, tt.gateway, tt.subnet)
				if err == nil {
					t.Error("SetStaticIP() expected error, got nil")
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

			err := client.SetStaticIP(context.Background(), tt.ip, tt.gateway, tt.subnet)
			if err != nil {
				t.Errorf("SetStaticIP() error: %v", err)
			}
		})
	}
}

func TestClient_EnableDHCP(t *testing.T) {
	t.Run("enable", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cmd := r.URL.Query().Get("cmnd")
			if !strings.Contains(cmd, "0.0.0.0") {
				t.Error("command should contain 0.0.0.0")
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"IPAddress1":"0.0.0.0"}`))
		}))
		defer server.Close()

		client := &Client{
			baseURL:    server.URL,
			httpClient: server.Client(),
		}

		err := client.EnableDHCP(context.Background(), true)
		if err != nil {
			t.Errorf("EnableDHCP(true) error: %v", err)
		}
	})

	t.Run("disable", func(t *testing.T) {
		client := &Client{}
		err := client.EnableDHCP(context.Background(), false)
		if err == nil {
			t.Error("EnableDHCP(false) expected error, got nil")
		}
	})
}

func TestClient_SetDNSServer(t *testing.T) {
	tests := []struct {
		name      string
		dnsServer string
		wantErr   bool
	}{
		{"valid DNS", "8.8.8.8", false},
		{"invalid DNS", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetDNSServer(context.Background(), tt.dnsServer)
				if err == nil {
					t.Error("SetDNSServer() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"IPAddress4":"` + tt.dnsServer + `"}`))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.SetDNSServer(context.Background(), tt.dnsServer)
			if err != nil {
				t.Errorf("SetDNSServer() error: %v", err)
			}
		})
	}
}

func TestClient_SetWiFi(t *testing.T) {
	tests := []struct {
		name     string
		ssid     string
		password string
		slot     int
		wantErr  bool
	}{
		{"valid slot 1", "MyWiFi", "password", 1, false},
		{"valid slot 2", "MyWiFi2", "password", 2, false},
		{"no password", "MyWiFi", "", 1, false},
		{"invalid slot 0", "MyWiFi", "password", 0, true},
		{"invalid slot 3", "MyWiFi", "password", 3, true},
		{"empty SSID", "", "password", 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetWiFi(context.Background(), tt.ssid, tt.password, tt.slot)
				if err == nil {
					t.Error("SetWiFi() expected error, got nil")
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

			err := client.SetWiFi(context.Background(), tt.ssid, tt.password, tt.slot)
			if err != nil {
				t.Errorf("SetWiFi() error: %v", err)
			}
		})
	}
}

func TestClient_GetSSID(t *testing.T) {
	mockResponse := `{"SSId1":"WiFi1","SSId2":"WiFi2"}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	ssids, err := client.GetSSID(context.Background())
	if err != nil {
		t.Fatalf("GetSSID() error: %v", err)
	}

	if len(ssids) != 2 {
		t.Errorf("GetSSID() returned %d SSIDs, want 2", len(ssids))
	}
	if ssids[0] != "WiFi1" {
		t.Errorf("SSId[0] = %v, want WiFi1", ssids[0])
	}
	if ssids[1] != "WiFi2" {
		t.Errorf("SSId[1] = %v, want WiFi2", ssids[1])
	}
}

func TestClient_SetAPMode(t *testing.T) {
	tests := []struct {
		name    string
		mode    int
		wantErr bool
	}{
		{"disable", 0, false},
		{"enable", 1, false},
		{"enable no auth", 2, false},
		{"invalid", 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetAPMode(context.Background(), tt.mode)
				if err == nil {
					t.Error("SetAPMode() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, `{"AP":%d}`, tt.mode)
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.SetAPMode(context.Background(), tt.mode)
			if err != nil {
				t.Errorf("SetAPMode() error: %v", err)
			}
		})
	}
}

func TestClient_WebPassword(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"WebPassword":"****"}`))
		}))
		defer server.Close()

		client := &Client{
			baseURL:    server.URL,
			httpClient: server.Client(),
		}

		err := client.SetWebPassword(context.Background(), "secret")
		if err != nil {
			t.Errorf("SetWebPassword() error: %v", err)
		}
	})

	t.Run("get", func(t *testing.T) {
		mockResponse := `{"WebPassword":1}`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(mockResponse))
		}))
		defer server.Close()

		client := &Client{
			baseURL:    server.URL,
			httpClient: server.Client(),
		}

		hasPassword, err := client.GetWebPassword(context.Background())
		if err != nil {
			t.Fatalf("GetWebPassword() error: %v", err)
		}
		if !hasPassword {
			t.Error("GetWebPassword() = false, want true")
		}
	})
}

func TestClient_SetNetworkConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *NetworkConfig
		wantErr bool
	}{
		{
			name: "full static config",
			config: &NetworkConfig{
				Hostname:  "tasmota",
				IPAddress: "192.168.1.100",
				Gateway:   "192.168.1.1",
				Subnet:    "255.255.255.0",
				DNSServer: "8.8.8.8",
				SSID1:     "WiFi1",
				Password1: "pass1",
			},
			wantErr: false,
		},
		{
			name: "DHCP config",
			config: &NetworkConfig{
				Hostname: "tasmota",
				UseDHCP:  true,
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
			config:  &NetworkConfig{},
			wantErr: true,
		},
		{
			name: "invalid IP",
			config: &NetworkConfig{
				IPAddress: "invalid",
				Gateway:   "192.168.1.1",
				Subnet:    "255.255.255.0",
			},
			wantErr: true,
		},
		{
			name: "hostname too long",
			config: &NetworkConfig{
				Hostname: strings.Repeat("a", 33),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetNetworkConfig(context.Background(), tt.config)
				if err == nil {
					t.Error("SetNetworkConfig() expected error, got nil")
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

			err := client.SetNetworkConfig(context.Background(), tt.config)
			if err != nil {
				t.Errorf("SetNetworkConfig() error: %v", err)
			}
		})
	}
}

func TestClient_GetIPConfig(t *testing.T) {
	mockResponse := `{
		"StatusNET": {
			"IPAddress": "192.168.1.100",
			"Gateway": "192.168.1.1",
			"Subnetmask": "255.255.255.0",
			"DNSServer": "192.168.1.1"
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

	ip, gateway, subnet, dns, err := client.GetIPConfig(context.Background())
	if err != nil {
		t.Fatalf("GetIPConfig() error: %v", err)
	}

	if ip != "192.168.1.100" {
		t.Errorf("IP = %v, want 192.168.1.100", ip)
	}
	if gateway != "192.168.1.1" {
		t.Errorf("Gateway = %v, want 192.168.1.1", gateway)
	}
	if subnet != "255.255.255.0" {
		t.Errorf("Subnet = %v, want 255.255.255.0", subnet)
	}
	if dns != "192.168.1.1" {
		t.Errorf("DNS = %v, want 192.168.1.1", dns)
	}
}

func TestClient_GetMACAddress(t *testing.T) {
	mockResponse := `{
		"StatusNET": {
			"Mac": "AA:BB:CC:DD:EE:FF"
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

	mac, err := client.GetMACAddress(context.Background())
	if err != nil {
		t.Fatalf("GetMACAddress() error: %v", err)
	}

	if mac != "AA:BB:CC:DD:EE:FF" {
		t.Errorf("MAC = %v, want AA:BB:CC:DD:EE:FF", mac)
	}
}

func TestClient_WiFiPower(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		mockResponse := `{
			"StatusNET": {
				"WifiPower": 17.0
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

		power, err := client.GetWiFiPower(context.Background())
		if err != nil {
			t.Fatalf("GetWiFiPower() error: %v", err)
		}
		if power != 17.0 {
			t.Errorf("GetWiFiPower() = %v, want 17.0", power)
		}
	})

	t.Run("set valid", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"WifiPower":17.0}`))
		}))
		defer server.Close()

		client := &Client{
			baseURL:    server.URL,
			httpClient: server.Client(),
		}

		err := client.SetWiFiPower(context.Background(), 17.0)
		if err != nil {
			t.Errorf("SetWiFiPower() error: %v", err)
		}
	})

	t.Run("set invalid", func(t *testing.T) {
		client := &Client{}

		tests := []float64{-1.0, 25.0}
		for _, val := range tests {
			err := client.SetWiFiPower(context.Background(), val)
			if err == nil {
				t.Errorf("SetWiFiPower(%v) expected error, got nil", val)
			}
		}
	})
}

func TestClient_SetWiFiConfig(t *testing.T) {
	tests := []struct {
		name    string
		mode    int
		wantErr bool
	}{
		{"restart", 0, false},
		{"smartconfig", 1, false},
		{"manager", 2, false},
		{"wps", 3, false},
		{"retry disable", 4, false},
		{"retry enable", 5, false},
		{"invalid", 6, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				client := &Client{}
				err := client.SetWiFiConfig(context.Background(), tt.mode)
				if err == nil {
					t.Error("SetWiFiConfig() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = fmt.Fprintf(w, `{"WifiConfig":%d}`, tt.mode)
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			err := client.SetWiFiConfig(context.Background(), tt.mode)
			if err != nil {
				t.Errorf("SetWiFiConfig() error: %v", err)
			}
		})
	}
}

func TestClient_Ping(t *testing.T) {
	tests := []struct {
		name         string
		host         string
		mockResponse string
		expected     bool
		wantErr      bool
	}{
		{
			name:         "successful ping",
			host:         "8.8.8.8",
			mockResponse: `{"Ping":"Success"}`,
			expected:     true,
			wantErr:      false,
		},
		{
			name:         "failed ping",
			host:         "192.168.1.200",
			mockResponse: `{"Ping":"Timeout"}`,
			expected:     false,
			wantErr:      false,
		},
		{
			name:    "empty host",
			host:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr && tt.host == "" {
				client := &Client{}
				_, err := client.Ping(context.Background(), tt.host)
				if err == nil {
					t.Error("Ping() expected error, got nil")
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client := &Client{
				baseURL:    server.URL,
				httpClient: server.Client(),
			}

			success, err := client.Ping(context.Background(), tt.host)
			if err != nil {
				t.Errorf("Ping() error: %v", err)
				return
			}
			if success != tt.expected {
				t.Errorf("Ping() = %v, want %v", success, tt.expected)
			}
		})
	}
}
