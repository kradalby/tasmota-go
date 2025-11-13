package tasmota

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_Status(t *testing.T) {
	tests := []struct {
		name         string
		category     int
		mockResponse string
		wantErr      bool
		checkField   string
	}{
		{
			name:         "status 0 all",
			category:     0,
			mockResponse: `{"Status":{"Module":1}}`,
			wantErr:      false,
			checkField:   "Status",
		},
		{
			name:         "status 1 device",
			category:     1,
			mockResponse: `{"Status":{"Module":1,"DeviceName":"Test"}}`,
			wantErr:      false,
			checkField:   "Status",
		},
		{
			name:         "status 2 firmware",
			category:     2,
			mockResponse: `{"StatusFWR":{"Version":"13.1.0"}}`,
			wantErr:      false,
			checkField:   "StatusFWR",
		},
		{
			name:         "status 5 network",
			category:     5,
			mockResponse: `{"StatusNET":{"Hostname":"tasmota"}}`,
			wantErr:      false,
			checkField:   "StatusNET",
		},
		{
			name:         "status 6 mqtt",
			category:     6,
			mockResponse: `{"StatusMQT":{"MqttHost":"localhost"}}`,
			wantErr:      false,
			checkField:   "StatusMQT",
		},
		{
			name:         "status 10 sensor",
			category:     10,
			mockResponse: `{"StatusSNS":{"Time":"2024-01-01T00:00:00"}}`,
			wantErr:      false,
			checkField:   "StatusSNS",
		},
		{
			name:         "status 11 state",
			category:     11,
			mockResponse: `{"StatusSTS":{"Uptime":"1T00:00:00"}}`,
			wantErr:      false,
			checkField:   "StatusSTS",
		},
		{
			name:     "invalid category negative",
			category: -1,
			wantErr:  true,
		},
		{
			name:     "invalid category too high",
			category: 12,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr && tt.mockResponse == "" {
				// Test validation without server
				client := &Client{}
				_, err := client.Status(context.Background(), tt.category)
				if err == nil {
					t.Error("Status() expected error, got nil")
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

			resp, err := client.Status(context.Background(), tt.category)
			if err != nil {
				t.Errorf("Status() unexpected error: %v", err)
				return
			}

			// Verify the expected field is present
			switch tt.checkField {
			case "Status":
				if resp.Status == nil {
					t.Error("Status field is nil")
				}
			case "StatusFWR":
				if resp.StatusFWR == nil {
					t.Error("StatusFWR field is nil")
				}
			case "StatusNET":
				if resp.StatusNET == nil {
					t.Error("StatusNET field is nil")
				}
			case "StatusMQT":
				if resp.StatusMQT == nil {
					t.Error("StatusMQT field is nil")
				}
			case "StatusSNS":
				if resp.StatusSNS == nil {
					t.Error("StatusSNS field is nil")
				}
			case "StatusSTS":
				if resp.StatusSTS == nil {
					t.Error("StatusSTS field is nil")
				}
			}
		})
	}
}

func TestClient_GetDeviceInfo(t *testing.T) {
	mockResponse := `{
		"Status": {
			"Module": 1,
			"DeviceName": "Tasmota",
			"FriendlyName": ["Living Room", "Kitchen"],
			"Topic": "tasmota_test",
			"Power": 1,
			"PowerOnState": 3,
			"LedState": 1
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

	info, err := client.GetDeviceInfo(context.Background())
	if err != nil {
		t.Fatalf("GetDeviceInfo() error: %v", err)
	}

	if info.DeviceName != "Tasmota" {
		t.Errorf("DeviceName = %v, want Tasmota", info.DeviceName)
	}
	if info.Module != 1 {
		t.Errorf("Module = %v, want 1", info.Module)
	}
	if len(info.FriendlyName) != 2 {
		t.Errorf("FriendlyName length = %v, want 2", len(info.FriendlyName))
	}
	if info.Topic != "tasmota_test" {
		t.Errorf("Topic = %v, want tasmota_test", info.Topic)
	}
}

func TestClient_GetFirmwareInfo(t *testing.T) {
	mockResponse := `{
		"StatusFWR": {
			"Version": "13.1.0",
			"BuildDateTime": "2024-01-01T12:00:00",
			"Core": "esp8266",
			"SDK": "2.2.2"
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

	info, err := client.GetFirmwareInfo(context.Background())
	if err != nil {
		t.Fatalf("GetFirmwareInfo() error: %v", err)
	}

	if info.Version != "13.1.0" {
		t.Errorf("Version = %v, want 13.1.0", info.Version)
	}
	if info.Core != "esp8266" {
		t.Errorf("Core = %v, want esp8266", info.Core)
	}
}

func TestClient_GetNetworkInfo(t *testing.T) {
	mockResponse := `{
		"StatusNET": {
			"Hostname": "tasmota-test",
			"IPAddress": "192.168.1.100",
			"Gateway": "192.168.1.1",
			"Subnetmask": "255.255.255.0",
			"DNSServer": "192.168.1.1",
			"Mac": "AA:BB:CC:DD:EE:FF",
			"Webserver": 2,
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

	info, err := client.GetNetworkInfo(context.Background())
	if err != nil {
		t.Fatalf("GetNetworkInfo() error: %v", err)
	}

	if info.Hostname != "tasmota-test" {
		t.Errorf("Hostname = %v, want tasmota-test", info.Hostname)
	}
	if info.IPAddress.String() != "192.168.1.100" {
		t.Errorf("IPAddress = %v, want 192.168.1.100", info.IPAddress)
	}
	if info.Mac.String() != "aa:bb:cc:dd:ee:ff" {
		t.Errorf("Mac = %v, want AA:BB:CC:DD:EE:FF", info.Mac)
	}
	if info.WifiPower != 17.0 {
		t.Errorf("WifiPower = %v, want 17.0", info.WifiPower)
	}
}

func TestClient_GetMQTTInfo(t *testing.T) {
	mockResponse := `{
		"StatusMQT": {
			"MqttHost": "mqtt.local",
			"MqttPort": 1883,
			"MqttClient": "tasmota-test",
			"MqttUser": "admin",
			"MqttCount": 5
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

	info, err := client.GetMQTTInfo(context.Background())
	if err != nil {
		t.Fatalf("GetMQTTInfo() error: %v", err)
	}

	if info.MqttHost != "mqtt.local" {
		t.Errorf("MqttHost = %v, want mqtt.local", info.MqttHost)
	}
	if info.MqttPort != 1883 {
		t.Errorf("MqttPort = %v, want 1883", info.MqttPort)
	}
	if info.MqttUser != "admin" {
		t.Errorf("MqttUser = %v, want admin", info.MqttUser)
	}
}

func TestClient_GetSensorData(t *testing.T) {
	mockResponse := `{
		"StatusSNS": {
			"Time": "2024-01-01T12:00:00",
			"ENERGY": {
				"TotalStartTime": "2024-01-01T00:00:00",
				"Total": 123.456,
				"Yesterday": 12.34,
				"Today": 5.67,
				"Power": 100.5,
				"ApparentPower": 105.2,
				"ReactivePower": 20.1,
				"Factor": 0.95,
				"Voltage": 230.0,
				"Current": 0.44
			}
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

	sensor, err := client.GetSensorData(context.Background())
	if err != nil {
		t.Fatalf("GetSensorData() error: %v", err)
	}

	if sensor.Time != "2024-01-01T12:00:00" {
		t.Errorf("Time = %v, want 2024-01-01T12:00:00", sensor.Time)
	}

	if sensor.Energy == nil {
		t.Fatal("Energy is nil")
	}

	if sensor.Energy.Power != 100.5 {
		t.Errorf("Power = %v, want 100.5", sensor.Energy.Power)
	}
	if sensor.Energy.Voltage != 230.0 {
		t.Errorf("Voltage = %v, want 230.0", sensor.Energy.Voltage)
	}
	if sensor.Energy.Current != 0.44 {
		t.Errorf("Current = %v, want 0.44", sensor.Energy.Current)
	}
	if sensor.Energy.Total != 123.456 {
		t.Errorf("Total = %v, want 123.456", sensor.Energy.Total)
	}
}

func TestClient_GetState(t *testing.T) {
	mockResponse := `{
		"StatusSTS": {
			"Time": "2024-01-01T12:00:00",
			"Uptime": "1T05:30:15",
			"UptimeSec": 106215,
			"Heap": 28,
			"SleepMode": "Dynamic",
			"Sleep": 50,
			"LoadAvg": 19,
			"MqttCount": 1,
			"POWER": "ON",
			"POWER1": "ON",
			"POWER2": "OFF",
			"Wifi": {
				"AP": 1,
				"SSId": "MyWiFi",
				"BSSId": "AA:BB:CC:DD:EE:FF",
				"Channel": 6,
				"Mode": "11n",
				"RSSI": -65,
				"Signal": 70,
				"LinkCount": 1,
				"Downtime": "0T00:00:03"
			}
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

	state, err := client.GetState(context.Background())
	if err != nil {
		t.Fatalf("GetState() error: %v", err)
	}

	if state.Uptime != "1T05:30:15" {
		t.Errorf("Uptime = %v, want 1T05:30:15", state.Uptime)
	}
	if state.UptimeSec != 106215 {
		t.Errorf("UptimeSec = %v, want 106215", state.UptimeSec)
	}
	if state.POWER != "ON" {
		t.Errorf("POWER = %v, want ON", state.POWER)
	}
	if state.POWER1 != "ON" {
		t.Errorf("POWER1 = %v, want ON", state.POWER1)
	}
	if state.POWER2 != "OFF" {
		t.Errorf("POWER2 = %v, want OFF", state.POWER2)
	}

	if state.Wifi == nil {
		t.Fatal("Wifi is nil")
	}
	if state.Wifi.SSId != "MyWiFi" {
		t.Errorf("SSId = %v, want MyWiFi", state.Wifi.SSId)
	}
	if state.Wifi.RSSI != -65 {
		t.Errorf("RSSI = %v, want -65", state.Wifi.RSSI)
	}
	if state.Wifi.Signal != 70 {
		t.Errorf("Signal = %v, want 70", state.Wifi.Signal)
	}
}

func TestClient_GetUptime(t *testing.T) {
	mockResponse := `{
		"StatusSTS": {
			"Time": "2024-01-01T12:00:00",
			"Uptime": "1T00:00:00",
			"UptimeSec": 86400
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

	uptime, err := client.GetUptime(context.Background())
	if err != nil {
		t.Fatalf("GetUptime() error: %v", err)
	}

	expected := 86400 * time.Second
	if uptime != expected {
		t.Errorf("GetUptime() = %v, want %v", uptime, expected)
	}
}

func TestClient_GetWiFiSignal(t *testing.T) {
	tests := []struct {
		name         string
		mockResponse string
		expected     int
		wantErr      bool
	}{
		{
			name: "with wifi info",
			mockResponse: `{
				"StatusSTS": {
					"Time": "2024-01-01T12:00:00",
					"Uptime": "1T00:00:00",
					"UptimeSec": 86400,
					"Wifi": {
						"RSSI": -55,
						"Signal": 90
					}
				}
			}`,
			expected: -55,
			wantErr:  false,
		},
		{
			name: "without wifi info",
			mockResponse: `{
				"StatusSTS": {
					"Time": "2024-01-01T12:00:00",
					"Uptime": "1T00:00:00",
					"UptimeSec": 86400
				}
			}`,
			wantErr: true,
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

			rssi, err := client.GetWiFiSignal(context.Background())
			if tt.wantErr {
				if err == nil {
					t.Error("GetWiFiSignal() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("GetWiFiSignal() unexpected error: %v", err)
					return
				}
				if rssi != tt.expected {
					t.Errorf("GetWiFiSignal() = %v, want %v", rssi, tt.expected)
				}
			}
		})
	}
}

func TestClient_GetDeviceInfo_MissingField(t *testing.T) {
	mockResponse := `{}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := &Client{
		baseURL:    server.URL,
		httpClient: server.Client(),
	}

	_, err := client.GetDeviceInfo(context.Background())
	if err == nil {
		t.Error("GetDeviceInfo() expected error for missing field, got nil")
	}
	if !IsParseError(err) {
		t.Errorf("expected parse error, got %T", err)
	}
}
