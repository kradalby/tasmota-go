package tasmota

import (
	"context"
	"fmt"
	"strings"
)

// DeviceConfig represents device configuration settings.
type DeviceConfig struct {
	FriendlyName []string
	DeviceName   string
	PowerOnState int
	LedState     int
	Sleep        int
	ButtonRetain int
	SwitchRetain int
	SensorRetain int
	PowerRetain  int
}

// GetConfig retrieves the current device configuration.
func (c *Client) GetConfig(ctx context.Context) (*DeviceConfig, error) {
	info, err := c.GetDeviceInfo(ctx)
	if err != nil {
		return nil, err
	}

	config := &DeviceConfig{
		FriendlyName: info.FriendlyName,
		DeviceName:   info.DeviceName,
		PowerOnState: info.PowerOnState,
		LedState:     info.LedState,
		ButtonRetain: info.ButtonRetain,
		SwitchRetain: info.SwitchRetain,
		SensorRetain: info.SensorRetain,
		PowerRetain:  info.PowerRetain,
	}

	return config, nil
}

// SetDeviceName sets the device name.
func (c *Client) SetDeviceName(ctx context.Context, name string) error {
	if name == "" {
		return NewError(ErrorTypeCommand, "device name cannot be empty", nil)
	}
	cmd := fmt.Sprintf("DeviceName %s", name)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetFriendlyName sets the first friendly name.
func (c *Client) SetFriendlyName(ctx context.Context, name string) error {
	return c.SetFriendlyNameN(ctx, 1, name)
}

// SetFriendlyNameN sets a specific friendly name (1-8).
func (c *Client) SetFriendlyNameN(ctx context.Context, index int, name string) error {
	if index < 1 || index > 8 {
		return NewError(ErrorTypeCommand, "friendly name index must be between 1 and 8", nil)
	}
	if name == "" {
		return NewError(ErrorTypeCommand, "friendly name cannot be empty", nil)
	}
	cmd := fmt.Sprintf("FriendlyName%d %s", index, name)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetPowerOnState sets the power state on boot.
// Values:
//   0 = Keep relay off after power on
//   1 = Turn relay on after power on
//   2 = Toggle relay after power on
//   3 = Set relay to last saved state after power on (default)
//   4 = Turn relay on and disable further relay control
//   5 = After a PulseTime period turn relay on (acts as inverted PulseTime mode)
func (c *Client) SetPowerOnState(ctx context.Context, state int) error {
	if state < 0 || state > 5 {
		return NewError(ErrorTypeCommand, "power on state must be between 0 and 5", nil)
	}
	cmd := fmt.Sprintf("PowerOnState %d", state)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetLedState sets the LED state.
// Values:
//   0 = Disable use of LED as much as possible
//   1 = Show power state on LED (default)
//   2 = Show MQTT subscriptions as a LED blink
//   3 = Show power state and MQTT subscriptions as a LED blink
//   4 = Show MQTT publications as a LED blink
//   5 = Show power state and MQTT publications as a LED blink
//   6 = Show all MQTT messages as a LED blink
//   7 = Show power state and MQTT messages as a LED blink
//   8 = LED on when Wi-Fi and MQTT are connected
func (c *Client) SetLedState(ctx context.Context, state int) error {
	if state < 0 || state > 8 {
		return NewError(ErrorTypeCommand, "LED state must be between 0 and 8", nil)
	}
	cmd := fmt.Sprintf("LedState %d", state)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetSleep sets the sleep mode.
// Values:
//   0 = Dynamic sleep disabled
//   1..250 = Set sleep duration in milliseconds (default = 50)
func (c *Client) SetSleep(ctx context.Context, duration int) error {
	if duration < 0 || duration > 250 {
		return NewError(ErrorTypeCommand, "sleep duration must be between 0 and 250", nil)
	}
	cmd := fmt.Sprintf("Sleep %d", duration)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetButtonRetain sets MQTT retain flag for button messages.
// 0 = disable retain, 1 = enable retain
func (c *Client) SetButtonRetain(ctx context.Context, retain bool) error {
	val := 0
	if retain {
		val = 1
	}
	cmd := fmt.Sprintf("ButtonRetain %d", val)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetSwitchRetain sets MQTT retain flag for switch messages.
// 0 = disable retain, 1 = enable retain
func (c *Client) SetSwitchRetain(ctx context.Context, retain bool) error {
	val := 0
	if retain {
		val = 1
	}
	cmd := fmt.Sprintf("SwitchRetain %d", val)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetSensorRetain sets MQTT retain flag for sensor messages.
// 0 = disable retain, 1 = enable retain
func (c *Client) SetSensorRetain(ctx context.Context, retain bool) error {
	val := 0
	if retain {
		val = 1
	}
	cmd := fmt.Sprintf("SensorRetain %d", val)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetPowerRetain sets MQTT retain flag for power messages.
// 0 = disable retain, 1 = enable retain
func (c *Client) SetPowerRetain(ctx context.Context, retain bool) error {
	val := 0
	if retain {
		val = 1
	}
	cmd := fmt.Sprintf("PowerRetain %d", val)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// ApplyConfig applies multiple configuration changes atomically using Backlog.
// This minimizes the number of HTTP requests and ensures all changes are applied together.
func (c *Client) ApplyConfig(ctx context.Context, cfg *DeviceConfig) error {
	if cfg == nil {
		return NewError(ErrorTypeCommand, "config cannot be nil", nil)
	}

	var commands []string

	// Device name
	if cfg.DeviceName != "" {
		commands = append(commands, fmt.Sprintf("DeviceName %s", cfg.DeviceName))
	}

	// Friendly names
	for i, name := range cfg.FriendlyName {
		if name != "" {
			commands = append(commands, fmt.Sprintf("FriendlyName%d %s", i+1, name))
		}
	}

	// Power on state
	if cfg.PowerOnState >= 0 && cfg.PowerOnState <= 5 {
		commands = append(commands, fmt.Sprintf("PowerOnState %d", cfg.PowerOnState))
	}

	// LED state
	if cfg.LedState >= 0 && cfg.LedState <= 8 {
		commands = append(commands, fmt.Sprintf("LedState %d", cfg.LedState))
	}

	// Sleep
	if cfg.Sleep >= 0 && cfg.Sleep <= 250 {
		commands = append(commands, fmt.Sprintf("Sleep %d", cfg.Sleep))
	}

	// Retain settings
	if cfg.ButtonRetain >= 0 && cfg.ButtonRetain <= 1 {
		commands = append(commands, fmt.Sprintf("ButtonRetain %d", cfg.ButtonRetain))
	}
	if cfg.SwitchRetain >= 0 && cfg.SwitchRetain <= 1 {
		commands = append(commands, fmt.Sprintf("SwitchRetain %d", cfg.SwitchRetain))
	}
	if cfg.SensorRetain >= 0 && cfg.SensorRetain <= 1 {
		commands = append(commands, fmt.Sprintf("SensorRetain %d", cfg.SensorRetain))
	}
	if cfg.PowerRetain >= 0 && cfg.PowerRetain <= 1 {
		commands = append(commands, fmt.Sprintf("PowerRetain %d", cfg.PowerRetain))
	}

	if len(commands) == 0 {
		return NewError(ErrorTypeCommand, "no valid configuration changes to apply", nil)
	}

	_, err := c.ExecuteBacklog(ctx, commands...)
	return err
}

// Restart restarts the device.
// reason: 1 = normal restart, 99 = reset to firmware defaults
func (c *Client) Restart(ctx context.Context, reason int) error {
	if reason != 1 && reason != 99 {
		return NewError(ErrorTypeCommand, "restart reason must be 1 (normal) or 99 (reset)", nil)
	}
	cmd := fmt.Sprintf("Restart %d", reason)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// Reset resets device configuration to defaults.
// Values:
//   1 = Reset relay settings
//   2 = Reset all settings except Wi-Fi
//   3 = Reset all settings except Wi-Fi and MQTT
//   4 = Reset Wi-Fi settings
//   5 = Erase all flash and reset parameters to firmware defaults but keep Wi-Fi settings
//   6 = Erase all flash and reset parameters to firmware defaults
//   99 = Reset device to firmware defaults and reboot (combines Reset 1 and Restart 1)
func (c *Client) Reset(ctx context.Context, level int) error {
	validLevels := []int{1, 2, 3, 4, 5, 6, 99}
	valid := false
	for _, v := range validLevels {
		if level == v {
			valid = true
			break
		}
	}
	if !valid {
		return NewError(ErrorTypeCommand, "invalid reset level", nil)
	}
	cmd := fmt.Sprintf("Reset %d", level)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// GetModule returns the module type and name.
func (c *Client) GetModule(ctx context.Context) (int, error) {
	info, err := c.GetDeviceInfo(ctx)
	if err != nil {
		return 0, err
	}
	return info.Module, nil
}

// SetOption sets a device option.
// Tasmota has many SetOption commands (0-150+) that control various behaviors.
func (c *Client) SetOption(ctx context.Context, option int, value interface{}) error {
	if option < 0 {
		return NewError(ErrorTypeCommand, "option number cannot be negative", nil)
	}

	var cmd string
	switch v := value.(type) {
	case bool:
		val := 0
		if v {
			val = 1
		}
		cmd = fmt.Sprintf("SetOption%d %d", option, val)
	case int:
		cmd = fmt.Sprintf("SetOption%d %d", option, v)
	case string:
		cmd = fmt.Sprintf("SetOption%d %s", option, v)
	default:
		return NewError(ErrorTypeCommand, "unsupported value type for SetOption", nil)
	}

	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// GetTelePeriod returns the telemetry period in seconds.
func (c *Client) GetTelePeriod(ctx context.Context) (int, error) {
	raw, err := c.ExecuteCommand(ctx, "TelePeriod")
	if err != nil {
		return 0, err
	}

	var result struct {
		TelePeriod int `json:"TelePeriod"`
	}
	if err := unmarshalJSON(raw, &result); err != nil {
		return 0, err
	}
	return result.TelePeriod, nil
}

// SetTelePeriod sets the telemetry period in seconds (10-3600).
func (c *Client) SetTelePeriod(ctx context.Context, seconds int) error {
	if seconds < 10 || seconds > 3600 {
		return NewError(ErrorTypeCommand, "telemetry period must be between 10 and 3600 seconds", nil)
	}
	cmd := fmt.Sprintf("TelePeriod %d", seconds)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}

// SetTemplate configures a device template.
// Template is a JSON string defining GPIO assignments.
func (c *Client) SetTemplate(ctx context.Context, template string) error {
	if template == "" {
		return NewError(ErrorTypeCommand, "template cannot be empty", nil)
	}
	// Escape quotes in template
	template = strings.ReplaceAll(template, `"`, `\"`)
	cmd := fmt.Sprintf("Template %s", template)
	_, err := c.ExecuteCommand(ctx, cmd)
	return err
}
