package tasmota

import (
	"context"
	"fmt"
)

// MQTTConfig represents MQTT broker configuration.
type MQTTConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	Client      string
	Topic       string
	FullTopic   string
	GroupTopic  string
	Retain      bool
	TelePeriod  int
	Prefix1     string // Command prefix (default: cmnd)
	Prefix2     string // Status prefix (default: stat)
	Prefix3     string // Telemetry prefix (default: tele)
}

// GetMQTTConfig retrieves the current MQTT configuration.
func (c *Client) GetMQTTConfig(ctx context.Context) (*MQTTConfig, error) {
	mqttInfo, err := c.GetMQTTInfo(ctx)
	if err != nil {
		return nil, err
	}

	deviceInfo, err := c.GetDeviceInfo(ctx)
	if err != nil {
		return nil, err
	}

	config := &MQTTConfig{
		Host:  mqttInfo.MqttHost,
		Port:  mqttInfo.MqttPort,
		User:  mqttInfo.MqttUser,
		Topic: deviceInfo.Topic,
	}

	return config, nil
}

// SetMQTTHost sets the MQTT broker hostname or IP address.
func (c *Client) SetMQTTHost(ctx context.Context, host string) error {
	if host == "" {
		return NewError(ErrorTypeCommand, "MQTT host cannot be empty", nil)
	}
	cmd := fmt.Sprintf("MqttHost %s", host)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetMQTTPort sets the MQTT broker port (1-65535).
func (c *Client) SetMQTTPort(ctx context.Context, port int) error {
	if port < 1 || port > 65535 {
		return NewError(ErrorTypeCommand, "MQTT port must be between 1 and 65535", nil)
	}
	cmd := fmt.Sprintf("MqttPort %d", port)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetMQTTUser sets the MQTT username.
func (c *Client) SetMQTTUser(ctx context.Context, username string) error {
	cmd := fmt.Sprintf("MqttUser %s", username)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetMQTTPassword sets the MQTT password.
func (c *Client) SetMQTTPassword(ctx context.Context, password string) error {
	cmd := fmt.Sprintf("MqttPassword %s", password)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetMQTTClient sets the MQTT client name.
func (c *Client) SetMQTTClient(ctx context.Context, clientName string) error {
	if clientName == "" {
		return NewError(ErrorTypeCommand, "MQTT client name cannot be empty", nil)
	}
	cmd := fmt.Sprintf("MqttClient %s", clientName)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetTopic sets the MQTT device topic.
func (c *Client) SetTopic(ctx context.Context, topic string) error {
	if topic == "" {
		return NewError(ErrorTypeCommand, "MQTT topic cannot be empty", nil)
	}
	cmd := fmt.Sprintf("Topic %s", topic)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetFullTopic sets the MQTT full topic template.
// Use tokens: %prefix%, %topic%, %hostname%, %id%
// Example: "%prefix%/%topic%/"
func (c *Client) SetFullTopic(ctx context.Context, fullTopic string) error {
	if fullTopic == "" {
		return NewError(ErrorTypeCommand, "MQTT full topic cannot be empty", nil)
	}
	cmd := fmt.Sprintf("FullTopic %s", fullTopic)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetGroupTopic sets the MQTT group topic for controlling multiple devices.
func (c *Client) SetGroupTopic(ctx context.Context, groupTopic string) error {
	if groupTopic == "" {
		return NewError(ErrorTypeCommand, "MQTT group topic cannot be empty", nil)
	}
	cmd := fmt.Sprintf("GroupTopic %s", groupTopic)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetMQTTRetain sets whether to retain MQTT messages.
func (c *Client) SetMQTTRetain(ctx context.Context, retain bool) error {
	val := 0
	if retain {
		val = 1
	}
	cmd := fmt.Sprintf("PowerRetain %d", val)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetPrefix sets an MQTT topic prefix.
// prefixNum: 1 = command prefix, 2 = status prefix, 3 = telemetry prefix
func (c *Client) SetPrefix(ctx context.Context, prefixNum int, prefix string) error {
	if prefixNum < 1 || prefixNum > 3 {
		return NewError(ErrorTypeCommand, "prefix number must be 1, 2, or 3", nil)
	}
	if prefix == "" {
		return NewError(ErrorTypeCommand, "prefix cannot be empty", nil)
	}
	cmd := fmt.Sprintf("Prefix%d %s", prefixNum, prefix)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// EnableMQTT enables or disables MQTT.
// Uses SetOption3: 0 = enable MQTT, 1 = disable MQTT
func (c *Client) EnableMQTT(ctx context.Context, enable bool) error {
	// SetOption3 is inverted: 0 = enable, 1 = disable
	val := 1
	if enable {
		val = 0
	}
	return c.SetOption(ctx, 3, val)
}

// SetMQTTConfig configures MQTT broker settings atomically using Backlog.
func (c *Client) SetMQTTConfig(ctx context.Context, cfg *MQTTConfig) error {
	if cfg == nil {
		return NewError(ErrorTypeCommand, "MQTT config cannot be nil", nil)
	}

	var commands []string

	// Enable MQTT first
	commands = append(commands, "SetOption3 0")

	// Host
	if cfg.Host != "" {
		commands = append(commands, fmt.Sprintf("MqttHost %s", cfg.Host))
	}

	// Port
	if cfg.Port > 0 && cfg.Port <= 65535 {
		commands = append(commands, fmt.Sprintf("MqttPort %d", cfg.Port))
	}

	// Authentication
	if cfg.User != "" {
		commands = append(commands, fmt.Sprintf("MqttUser %s", cfg.User))
	}
	if cfg.Password != "" {
		commands = append(commands, fmt.Sprintf("MqttPassword %s", cfg.Password))
	}

	// Client name
	if cfg.Client != "" {
		commands = append(commands, fmt.Sprintf("MqttClient %s", cfg.Client))
	}

	// Topics
	if cfg.Topic != "" {
		commands = append(commands, fmt.Sprintf("Topic %s", cfg.Topic))
	}
	if cfg.FullTopic != "" {
		commands = append(commands, fmt.Sprintf("FullTopic %s", cfg.FullTopic))
	}
	if cfg.GroupTopic != "" {
		commands = append(commands, fmt.Sprintf("GroupTopic %s", cfg.GroupTopic))
	}

	// Prefixes
	if cfg.Prefix1 != "" {
		commands = append(commands, fmt.Sprintf("Prefix1 %s", cfg.Prefix1))
	}
	if cfg.Prefix2 != "" {
		commands = append(commands, fmt.Sprintf("Prefix2 %s", cfg.Prefix2))
	}
	if cfg.Prefix3 != "" {
		commands = append(commands, fmt.Sprintf("Prefix3 %s", cfg.Prefix3))
	}

	// Retain
	if cfg.Retain {
		commands = append(commands, "PowerRetain 1")
	}

	// Telemetry period
	if cfg.TelePeriod >= 10 && cfg.TelePeriod <= 3600 {
		commands = append(commands, fmt.Sprintf("TelePeriod %d", cfg.TelePeriod))
	}

	if len(commands) <= 1 { // Only SetOption3
		return NewError(ErrorTypeCommand, "no valid MQTT configuration changes to apply", nil)
	}

	_, err := c.ExecuteBacklog(ctx, commands...)
	return err
}

// GetMQTTFingerprint returns the TLS fingerprint for MQTT.
func (c *Client) GetMQTTFingerprint(ctx context.Context) (string, error) {
	raw, err := c.ExecuteCommand(ctx, "MqttFingerprint")
	if err != nil {
		return "", err
	}

	var result struct {
		MqttFingerprint string `json:"MqttFingerprint"`
	}
	if err := unmarshalJSON(raw, &result); err != nil {
		return "", err
	}
	return result.MqttFingerprint, nil
}

// SetMQTTFingerprint sets the TLS fingerprint for MQTT.
// Use "00 00 00..." to disable fingerprint validation.
func (c *Client) SetMQTTFingerprint(ctx context.Context, fingerprint string) error {
	if fingerprint == "" {
		return NewError(ErrorTypeCommand, "MQTT fingerprint cannot be empty", nil)
	}
	cmd := fmt.Sprintf("MqttFingerprint %s", fingerprint)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// GetMQTTRetry returns the MQTT connection retry time in seconds.
func (c *Client) GetMQTTRetry(ctx context.Context) (int, error) {
	raw, err := c.ExecuteCommand(ctx, "MqttRetry")
	if err != nil {
		return 0, err
	}

	var result struct {
		MqttRetry int `json:"MqttRetry"`
	}
	if err := unmarshalJSON(raw, &result); err != nil {
		return 0, err
	}
	return result.MqttRetry, nil
}

// SetMQTTRetry sets the MQTT connection retry time in seconds (10-32000).
func (c *Client) SetMQTTRetry(ctx context.Context, seconds int) error {
	if seconds < 10 || seconds > 32000 {
		return NewError(ErrorTypeCommand, "MQTT retry must be between 10 and 32000 seconds", nil)
	}
	cmd := fmt.Sprintf("MqttRetry %d", seconds)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// TestMQTTConnection verifies MQTT connectivity by checking the MQTT count.
// A non-zero count indicates successful connection.
func (c *Client) TestMQTTConnection(ctx context.Context) error {
	mqttInfo, err := c.GetMQTTInfo(ctx)
	if err != nil {
		return err
	}

	if mqttInfo.MqttHost == "" {
		return NewError(ErrorTypeDevice, "MQTT host not configured", nil)
	}

	// MQTT count will be > 0 if connected
	if mqttInfo.MqttCount == 0 {
		return NewError(ErrorTypeDevice, "MQTT not connected", nil)
	}

	return nil
}
