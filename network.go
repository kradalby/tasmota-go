package tasmota

import (
	"context"
	"fmt"
	"net"
	"strings"
)

// NetworkConfig represents network configuration settings.
type NetworkConfig struct {
	Hostname  string
	IPAddress string
	Gateway   string
	Subnet    string
	DNSServer string
	SSID1     string
	SSID2     string
	Password1 string
	Password2 string
	UseDHCP   bool
}

// GetNetworkConfig retrieves the current network configuration.
func (c *Client) GetNetworkConfig(ctx context.Context) (*NetworkConfig, error) {
	netInfo, err := c.GetNetworkInfo(ctx)
	if err != nil {
		return nil, err
	}

	config := &NetworkConfig{
		Hostname:  netInfo.Hostname,
		IPAddress: netInfo.IPAddress,
		Gateway:   netInfo.Gateway,
		Subnet:    netInfo.Subnetmask,
		DNSServer: netInfo.DNSServer,
		UseDHCP:   netInfo.IPAddress == "0.0.0.0",
	}

	return config, nil
}

// SetHostname sets the device hostname.
func (c *Client) SetHostname(ctx context.Context, hostname string) error {
	if hostname == "" {
		return NewError(ErrorTypeCommand, "hostname cannot be empty", nil)
	}
	if len(hostname) > 32 {
		return NewError(ErrorTypeCommand, "hostname cannot exceed 32 characters", nil)
	}
	cmd := fmt.Sprintf("Hostname %s", hostname)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetStaticIP configures a static IP address.
func (c *Client) SetStaticIP(ctx context.Context, ip, gateway, subnet string) error {
	// Validate IP address
	if net.ParseIP(ip) == nil {
		return NewError(ErrorTypeCommand, "invalid IP address", nil)
	}
	if net.ParseIP(gateway) == nil {
		return NewError(ErrorTypeCommand, "invalid gateway address", nil)
	}
	if net.ParseIP(subnet) == nil {
		return NewError(ErrorTypeCommand, "invalid subnet mask", nil)
	}

	var commands []string
	commands = append(commands, fmt.Sprintf("IPAddress1 %s", ip))
	commands = append(commands, fmt.Sprintf("IPAddress2 %s", gateway))
	commands = append(commands, fmt.Sprintf("IPAddress3 %s", subnet))

	_, err := c.ExecuteBacklog(ctx, commands...)
	return err
}

// EnableDHCP enables or disables DHCP.
func (c *Client) EnableDHCP(ctx context.Context, enable bool) error {
	if enable {
		// Set IP to 0.0.0.0 to enable DHCP
		cmd := "IPAddress1 0.0.0.0"
		_, err := c.ExecuteCommand(ctx, cmd)
		return err
	}
	// To disable DHCP, you must set a static IP instead
	return NewError(ErrorTypeCommand, "to disable DHCP, use SetStaticIP", nil)
}

// SetDNSServer sets the DNS server address.
func (c *Client) SetDNSServer(ctx context.Context, dnsServer string) error {
	if net.ParseIP(dnsServer) == nil {
		return NewError(ErrorTypeCommand, "invalid DNS server address", nil)
	}
	cmd := fmt.Sprintf("IPAddress4 %s", dnsServer)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetWiFi configures WiFi credentials.
// slot should be 1 or 2 for AP1 or AP2.
func (c *Client) SetWiFi(ctx context.Context, ssid, password string, slot int) error {
	if slot < 1 || slot > 2 {
		return NewError(ErrorTypeCommand, "WiFi slot must be 1 or 2", nil)
	}
	if ssid == "" {
		return NewError(ErrorTypeCommand, "SSID cannot be empty", nil)
	}

	var commands []string
	commands = append(commands, fmt.Sprintf("SSId%d %s", slot, ssid))
	if password != "" {
		commands = append(commands, fmt.Sprintf("Password%d %s", slot, password))
	}

	_, err := c.ExecuteBacklog(ctx, commands...)
	return err
}

// GetSSID returns the configured SSIDs.
func (c *Client) GetSSID(ctx context.Context) ([]string, error) {
	raw, err := c.ExecuteCommand(ctx, "SSId")
	if err != nil {
		return nil, err
	}

	var result struct {
		SSId1 string `json:"SSId1"`
		SSId2 string `json:"SSId2"`
	}
	if err := unmarshalJSON(raw, &result); err != nil {
		return nil, err
	}

	var ssids []string
	if result.SSId1 != "" {
		ssids = append(ssids, result.SSId1)
	}
	if result.SSId2 != "" {
		ssids = append(ssids, result.SSId2)
	}

	return ssids, nil
}

// SetAPMode sets the WiFi access point mode.
// Values:
//   0 = Disable AP mode
//   1 = Enable AP mode (default)
//   2 = Enable AP mode with no authentication
func (c *Client) SetAPMode(ctx context.Context, mode int) error {
	if mode < 0 || mode > 2 {
		return NewError(ErrorTypeCommand, "AP mode must be 0, 1, or 2", nil)
	}
	cmd := fmt.Sprintf("AP %d", mode)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetWebPassword sets the web UI password.
func (c *Client) SetWebPassword(ctx context.Context, password string) error {
	cmd := fmt.Sprintf("WebPassword %s", password)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// GetWebPassword returns whether a web password is set.
func (c *Client) GetWebPassword(ctx context.Context) (bool, error) {
	raw, err := c.ExecuteCommand(ctx, "WebPassword")
	if err != nil {
		return false, err
	}

	var result struct {
		WebPassword int `json:"WebPassword"`
	}
	if err := unmarshalJSON(raw, &result); err != nil {
		return false, err
	}

	return result.WebPassword == 1, nil
}

// SetNetworkConfig applies multiple network configuration changes atomically using Backlog.
func (c *Client) SetNetworkConfig(ctx context.Context, cfg *NetworkConfig) error {
	if cfg == nil {
		return NewError(ErrorTypeCommand, "network config cannot be nil", nil)
	}

	var commands []string

	// Hostname
	if cfg.Hostname != "" {
		if len(cfg.Hostname) > 32 {
			return NewError(ErrorTypeCommand, "hostname cannot exceed 32 characters", nil)
		}
		commands = append(commands, fmt.Sprintf("Hostname %s", cfg.Hostname))
	}

	// IP configuration
	if cfg.UseDHCP {
		commands = append(commands, "IPAddress1 0.0.0.0")
	} else if cfg.IPAddress != "" && cfg.Gateway != "" && cfg.Subnet != "" {
		// Validate IPs
		if net.ParseIP(cfg.IPAddress) == nil {
			return NewError(ErrorTypeCommand, "invalid IP address", nil)
		}
		if net.ParseIP(cfg.Gateway) == nil {
			return NewError(ErrorTypeCommand, "invalid gateway address", nil)
		}
		if net.ParseIP(cfg.Subnet) == nil {
			return NewError(ErrorTypeCommand, "invalid subnet mask", nil)
		}

		commands = append(commands, fmt.Sprintf("IPAddress1 %s", cfg.IPAddress))
		commands = append(commands, fmt.Sprintf("IPAddress2 %s", cfg.Gateway))
		commands = append(commands, fmt.Sprintf("IPAddress3 %s", cfg.Subnet))
	}

	// DNS server
	if cfg.DNSServer != "" {
		if net.ParseIP(cfg.DNSServer) == nil {
			return NewError(ErrorTypeCommand, "invalid DNS server address", nil)
		}
		commands = append(commands, fmt.Sprintf("IPAddress4 %s", cfg.DNSServer))
	}

	// WiFi credentials
	if cfg.SSID1 != "" {
		commands = append(commands, fmt.Sprintf("SSId1 %s", cfg.SSID1))
		if cfg.Password1 != "" {
			commands = append(commands, fmt.Sprintf("Password1 %s", cfg.Password1))
		}
	}
	if cfg.SSID2 != "" {
		commands = append(commands, fmt.Sprintf("SSId2 %s", cfg.SSID2))
		if cfg.Password2 != "" {
			commands = append(commands, fmt.Sprintf("Password2 %s", cfg.Password2))
		}
	}

	if len(commands) == 0 {
		return NewError(ErrorTypeCommand, "no valid network configuration changes to apply", nil)
	}

	_, err := c.ExecuteBacklog(ctx, commands...)
	return err
}

// GetIPConfig returns the current IP configuration.
func (c *Client) GetIPConfig(ctx context.Context) (ip, gateway, subnet, dns string, err error) {
	netInfo, err := c.GetNetworkInfo(ctx)
	if err != nil {
		return "", "", "", "", err
	}

	return netInfo.IPAddress, netInfo.Gateway, netInfo.Subnetmask, netInfo.DNSServer, nil
}

// GetMACAddress returns the device MAC address.
func (c *Client) GetMACAddress(ctx context.Context) (string, error) {
	netInfo, err := c.GetNetworkInfo(ctx)
	if err != nil {
		return "", err
	}
	return netInfo.Mac, nil
}

// SetWiFiPower sets the WiFi transmit power in dBm.
// Values: 0-20.5 dBm (default: 17.0)
func (c *Client) SetWiFiPower(ctx context.Context, power float64) error {
	if power < 0 || power > 20.5 {
		return NewError(ErrorTypeCommand, "WiFi power must be between 0 and 20.5 dBm", nil)
	}
	cmd := fmt.Sprintf("WifiPower %.1f", power)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// GetWiFiPower returns the current WiFi transmit power in dBm.
func (c *Client) GetWiFiPower(ctx context.Context) (float64, error) {
	netInfo, err := c.GetNetworkInfo(ctx)
	if err != nil {
		return 0, err
	}
	return netInfo.WifiPower, nil
}

// SetWiFiConfig sets the WiFi configuration mode.
// Values:
//   0 = WIFI_RESTART: Reset Wi-Fi and restart
//   1 = WIFI_SMARTCONFIG: Start smart config for 1 minute
//   2 = WIFI_MANAGER: Start WiFi manager for 3 minutes
//   3 = WIFI_WPSCONFIG: Start WPS for 1 minute
//   4 = WIFI_RETRY: Disable Wi-Fi auto-restart
//   5 = WIFI_WAIT: Enable Wi-Fi auto-restart
func (c *Client) SetWiFiConfig(ctx context.Context, mode int) error {
	if mode < 0 || mode > 5 {
		return NewError(ErrorTypeCommand, "WiFi config mode must be between 0 and 5", nil)
	}
	cmd := fmt.Sprintf("WifiConfig %d", mode)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// Ping sends a ping to a host and returns whether it was successful.
func (c *Client) Ping(ctx context.Context, host string) (bool, error) {
	if host == "" {
		return false, NewError(ErrorTypeCommand, "ping host cannot be empty", nil)
	}

	cmd := fmt.Sprintf("Ping %s", host)
	raw, err := c.ExecuteCommand(ctx, cmd)
	if err != nil {
		return false, err
	}

	// Check if response contains "Success" or similar
	response := strings.ToLower(string(raw))
	success := strings.Contains(response, "success") || strings.Contains(response, "reply")

	return success, nil
}
