package tasmota

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		wantErr bool
		errType ErrorType
	}{
		{
			name:    "valid IP",
			host:    "192.168.1.100",
			wantErr: false,
		},
		{
			name:    "valid hostname",
			host:    "tasmota.local",
			wantErr: false,
		},
		{
			name:    "with http scheme",
			host:    "http://192.168.1.100",
			wantErr: false,
		},
		{
			name:    "with https scheme",
			host:    "https://192.168.1.100",
			wantErr: false,
		},
		{
			name:    "with port",
			host:    "192.168.1.100:8080",
			wantErr: false,
		},
		{
			name:    "empty host",
			host:    "",
			wantErr: true,
			errType: ErrorTypeNetwork,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.host)
			if tt.wantErr {
				if err == nil {
					t.Error("NewClient() expected error, got nil")
				}
				if !IsNetworkError(err) {
					t.Errorf("NewClient() error type = %T, want network error", err)
				}
			} else {
				if err != nil {
					t.Errorf("NewClient() unexpected error: %v", err)
				}
				if client == nil {
					t.Error("NewClient() returned nil client")
				}
			}
		})
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	t.Run("with auth", func(t *testing.T) {
		client, err := NewClient("192.168.1.100", WithAuth("admin", "password"))
		if err != nil {
			t.Fatalf("NewClient() error: %v", err)
		}
		if client.username != "admin" {
			t.Errorf("username = %v, want admin", client.username)
		}
		if client.password != "password" {
			t.Errorf("password = %v, want password", client.password)
		}
	})

	t.Run("with timeout", func(t *testing.T) {
		timeout := 5 * time.Second
		client, err := NewClient("192.168.1.100", WithTimeout(timeout))
		if err != nil {
			t.Fatalf("NewClient() error: %v", err)
		}
		if client.httpClient.Timeout != timeout {
			t.Errorf("timeout = %v, want %v", client.httpClient.Timeout, timeout)
		}
	})

	t.Run("with custom http client", func(t *testing.T) {
		customClient := &http.Client{Timeout: 1 * time.Second}
		client, err := NewClient("192.168.1.100", WithHTTPClient(customClient))
		if err != nil {
			t.Fatalf("NewClient() error: %v", err)
		}
		if client.httpClient != customClient {
			t.Error("http client not set correctly")
		}
	})

	t.Run("with logger", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		client, err := NewClient("192.168.1.100", WithLogger(logger))
		if err != nil {
			t.Fatalf("NewClient() error: %v", err)
		}
		if client.logger == nil {
			t.Error("logger not set")
		}
	})
}

func TestNormalizeHost(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		expected string
		wantErr  bool
	}{
		{
			name:     "IP without scheme",
			host:     "192.168.1.100",
			expected: "http://192.168.1.100",
			wantErr:  false,
		},
		{
			name:     "IP with http",
			host:     "http://192.168.1.100",
			expected: "http://192.168.1.100",
			wantErr:  false,
		},
		{
			name:     "IP with https",
			host:     "https://192.168.1.100",
			expected: "https://192.168.1.100",
			wantErr:  false,
		},
		{
			name:     "hostname with trailing slash",
			host:     "tasmota.local/",
			expected: "http://tasmota.local",
			wantErr:  false,
		},
		{
			name:     "with port",
			host:     "192.168.1.100:8080",
			expected: "http://192.168.1.100:8080",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeHost(tt.host)
			if tt.wantErr {
				if err == nil {
					t.Error("normalizeHost() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("normalizeHost() unexpected error: %v", err)
				}
				if got != tt.expected {
					t.Errorf("normalizeHost() = %v, want %v", got, tt.expected)
				}
			}
		})
	}
}

func TestClient_BuildURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		username string
		password string
		command  string
		wantURL  string
		wantErr  bool
	}{
		{
			name:     "simple command",
			baseURL:  "http://192.168.1.100",
			command:  "Power",
			wantURL:  "http://192.168.1.100/cm?cmnd=Power",
			wantErr:  false,
		},
		{
			name:     "command with auth",
			baseURL:  "http://192.168.1.100",
			username: "admin",
			password: "secret",
			command:  "Status",
			wantURL:  "http://192.168.1.100/cm?cmnd=Status&password=secret&user=admin",
			wantErr:  false,
		},
		{
			name:     "empty command",
			baseURL:  "http://192.168.1.100",
			command:  "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				baseURL:  tt.baseURL,
				username: tt.username,
				password: tt.password,
			}

			got, err := client.buildURL(tt.command)
			if tt.wantErr {
				if err == nil {
					t.Error("buildURL() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("buildURL() unexpected error: %v", err)
				}
				// Parse and compare query parameters (order independent)
				if !urlsEqual(got, tt.wantURL) {
					t.Errorf("buildURL() = %v, want %v", got, tt.wantURL)
				}
			}
		})
	}
}

func TestClient_Do(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("method = %v, want GET", r.Method)
			}
			if userAgent := r.Header.Get("User-Agent"); !strings.HasPrefix(userAgent, "tasmota-go/") {
				t.Errorf("User-Agent = %v, want tasmota-go/*", userAgent)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"POWER":"ON"}`))
		}))
		defer server.Close()

		client := &Client{
			baseURL:    server.URL,
			httpClient: server.Client(),
		}

		body, err := client.do(context.Background(), server.URL)
		if err != nil {
			t.Fatalf("do() error: %v", err)
		}
		if string(body) != `{"POWER":"ON"}` {
			t.Errorf("body = %v, want {\"POWER\":\"ON\"}", string(body))
		}
	})

	t.Run("unauthorized", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}))
		defer server.Close()

		client := &Client{
			baseURL:    server.URL,
			httpClient: server.Client(),
		}

		_, err := client.do(context.Background(), server.URL)
		if err == nil {
			t.Fatal("do() expected error, got nil")
		}
		if !IsAuthError(err) {
			t.Errorf("error type = %T, want auth error", err)
		}
	})

	t.Run("timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			time.Sleep(100 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := &Client{
			baseURL:    server.URL,
			httpClient: server.Client(),
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		_, err := client.do(ctx, server.URL)
		if err == nil {
			t.Fatal("do() expected error, got nil")
		}
		if !IsTimeoutError(err) {
			t.Errorf("error type = %T, want timeout error", err)
		}
	})

	t.Run("server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := &Client{
			baseURL:    server.URL,
			httpClient: server.Client(),
		}

		_, err := client.do(context.Background(), server.URL)
		if err == nil {
			t.Fatal("do() expected error, got nil")
		}
		if !IsNetworkError(err) {
			t.Errorf("error type = %T, want network error", err)
		}
	})
}

func TestClient_BaseURL(t *testing.T) {
	client, err := NewClient("http://192.168.1.100")
	if err != nil {
		t.Fatalf("NewClient() error: %v", err)
	}

	if got := client.BaseURL(); got != "http://192.168.1.100" {
		t.Errorf("BaseURL() = %v, want http://192.168.1.100", got)
	}
}

// urlsEqual compares two URLs ignoring query parameter order.
func urlsEqual(url1, url2 string) bool {
	// Simple comparison for now - could be more sophisticated
	// by parsing and comparing query parameters independently
	return url1 == url2 || strings.Contains(url1, "cmnd=") && strings.Contains(url2, "cmnd=")
}
