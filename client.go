// Package tasmota provides a Go client library for controlling and configuring Tasmota smart devices.
package tasmota

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultConnectTimeout is the default connection timeout.
	DefaultConnectTimeout = 10 * time.Second
	// DefaultResponseTimeout is the default response timeout.
	DefaultResponseTimeout = 20 * time.Second
)

// Client represents a Tasmota device client.
type Client struct {
	baseURL    string
	httpClient *http.Client
	username   string
	password   string
	logger     *slog.Logger
	timeout    *time.Duration
}

// ClientOption is a functional option for configuring the Client.
type ClientOption func(*Client)

// WithAuth configures authentication credentials.
func WithAuth(username, password string) ClientOption {
	return func(c *Client) {
		c.username = username
		c.password = password
	}
}

// WithTimeout configures the HTTP client timeout. It takes effect regardless
// of option order and never modifies a client passed to WithHTTPClient.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = &timeout
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithLogger sets a custom slog.Logger for request/response logging.
// If not set, no logging will be performed.
func WithLogger(logger *slog.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// NewClient creates a new Tasmota client for the specified host.
// The host can be an IP address (192.168.1.100) or hostname with optional port.
// If no scheme is provided, http:// will be used.
func NewClient(host string, opts ...ClientOption) (*Client, error) {
	if host == "" {
		return nil, NewError(ErrorTypeNetwork, "host cannot be empty", nil)
	}

	// Parse and normalize the host
	baseURL, err := normalizeHost(host)
	if err != nil {
		return nil, NewError(ErrorTypeNetwork, "invalid host", err)
	}

	client := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: DefaultResponseTimeout,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: DefaultConnectTimeout,
				}).DialContext,
			},
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	// Set on a copy: the client from WithHTTPClient may be shared, even
	// http.DefaultClient.
	if client.timeout != nil {
		hc := *client.httpClient
		hc.Timeout = *client.timeout
		client.httpClient = &hc
	}

	return client, nil
}

// normalizeHost ensures the host has a scheme and returns a clean base URL.
func normalizeHost(host string) (string, error) {
	// Remove trailing slashes
	host = strings.TrimRight(host, "/")

	// Add scheme if missing
	if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
		host = "http://" + host
	}

	// Parse to validate
	u, err := url.Parse(host)
	if err != nil {
		return "", err
	}

	if u.Host == "" {
		return "", fmt.Errorf("invalid host: %s", host)
	}

	return host, nil
}

// buildURL constructs the full URL for a command.
func (c *Client) buildURL(command string) (string, error) {
	if command == "" {
		return "", NewError(ErrorTypeCommand, "command cannot be empty", nil)
	}

	u, err := url.Parse(c.baseURL)
	if err != nil {
		return "", NewError(ErrorTypeNetwork, "invalid base URL", err)
	}

	u.Path = "/cm"
	q := u.Query()
	q.Set("cmnd", command)

	// Add authentication if configured
	if c.username != "" {
		q.Set("user", c.username)
	}
	if c.password != "" {
		q.Set("password", c.password)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

// do executes an HTTP GET request and returns the response body.
// classifyTransportError maps a transport-level failure onto the library's
// error taxonomy. A timeout cannot be recognised from the context alone:
// http.Client.Timeout -- set by WithTimeout, and by DefaultResponseTimeout on
// every client built with NewClient -- and the dialer timeout both fire
// without ever populating ctx.Err(). net.Error is what actually reports them,
// and it does so for a stalled response body just as much as for a request
// that never got a reply, so both call sites share this.
func classifyTransportError(ctx context.Context, err error, networkMessage string) error {
	if netErr, ok := errors.AsType[net.Error](err); ok && netErr.Timeout() {
		return NewError(ErrorTypeTimeout, "request timeout", err)
	}
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return NewError(ErrorTypeTimeout, "request timeout", err)
	}
	return NewError(ErrorTypeNetwork, networkMessage, err)
}

func (c *Client) do(ctx context.Context, urlStr string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, NewError(ErrorTypeNetwork, "failed to create request", err)
	}

	req.Header.Set("User-Agent", UserAgent)

	if c.logger != nil {
		c.logger.Debug("sending request",
			"method", req.Method,
			"url", urlStr)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if c.logger != nil {
			c.logger.Error("request failed",
				"url", urlStr,
				"error", err)
		}
		return nil, classifyTransportError(ctx, err, "request failed")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		if c.logger != nil {
			c.logger.Error("failed to read response body",
				"error", err)
		}
		return nil, classifyTransportError(ctx, err, "failed to read response")
	}

	if c.logger != nil {
		c.logger.Debug("received response",
			"status_code", resp.StatusCode,
			"body_length", len(body),
			"body", string(body))
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			return nil, NewError(ErrorTypeAuth, "authentication failed", nil)
		}
		return nil, NewError(ErrorTypeNetwork,
			fmt.Sprintf("unexpected status code: %d", resp.StatusCode), nil)
	}

	return body, nil
}

// BaseURL returns the base URL of the client.
func (c *Client) BaseURL() string {
	return c.baseURL
}
