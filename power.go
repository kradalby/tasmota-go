package tasmota

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// PowerState represents the state of a power relay.
type PowerState string

const (
	// PowerOff turns the relay off.
	PowerOff PowerState = "OFF"
	// PowerOn turns the relay on.
	PowerOn PowerState = "ON"
	// PowerToggle toggles the relay state.
	PowerToggle PowerState = "TOGGLE"
	// PowerBlink blinks the relay.
	PowerBlink PowerState = "BLINK"
)

// PowerResponse represents the response from a power command.
type PowerResponse struct {
	Power  string `json:"POWER,omitempty"`
	Power1 string `json:"POWER1,omitempty"`
	Power2 string `json:"POWER2,omitempty"`
	Power3 string `json:"POWER3,omitempty"`
	Power4 string `json:"POWER4,omitempty"`
	Power5 string `json:"POWER5,omitempty"`
	Power6 string `json:"POWER6,omitempty"`
	Power7 string `json:"POWER7,omitempty"`
	Power8 string `json:"POWER8,omitempty"`
}

// IsOn checks if a specific relay is on.
// relayNum should be 0 for POWER, 1-8 for POWER1-POWER8.
func (p *PowerResponse) IsOn(relayNum int) bool {
	state := p.GetState(relayNum)
	return strings.ToUpper(state) == "ON"
}

// GetState returns the state of a specific relay.
// relayNum should be 0 for POWER, 1-8 for POWER1-POWER8.
func (p *PowerResponse) GetState(relayNum int) string {
	switch relayNum {
	case 0:
		return p.Power
	case 1:
		return p.Power1
	case 2:
		return p.Power2
	case 3:
		return p.Power3
	case 4:
		return p.Power4
	case 5:
		return p.Power5
	case 6:
		return p.Power6
	case 7:
		return p.Power7
	case 8:
		return p.Power8
	default:
		return ""
	}
}

// Power controls all relays or the main relay.
// state can be PowerOn, PowerOff, PowerToggle, or PowerBlink.
func (c *Client) Power(ctx context.Context, state PowerState) (*PowerResponse, error) {
	cmd := fmt.Sprintf("Power %s", state)
	return c.executePowerCommand(ctx, cmd)
}

// PowerN controls a specific relay (1-8).
// relayNum should be 1-8.
// state can be PowerOn, PowerOff, PowerToggle, or PowerBlink.
func (c *Client) PowerN(ctx context.Context, relayNum int, state PowerState) (*PowerResponse, error) {
	if relayNum < 1 || relayNum > 8 {
		return nil, NewError(ErrorTypeCommand, "relay number must be between 1 and 8", nil)
	}
	cmd := fmt.Sprintf("Power%d %s", relayNum, state)
	return c.executePowerCommand(ctx, cmd)
}

// GetPower returns the current power state of all relays.
func (c *Client) GetPower(ctx context.Context) (*PowerResponse, error) {
	return c.executePowerCommand(ctx, "Power")
}

// GetPowerN returns the current power state of a specific relay (1-8).
func (c *Client) GetPowerN(ctx context.Context, relayNum int) (*PowerResponse, error) {
	if relayNum < 1 || relayNum > 8 {
		return nil, NewError(ErrorTypeCommand, "relay number must be between 1 and 8", nil)
	}
	cmd := fmt.Sprintf("Power%d", relayNum)
	return c.executePowerCommand(ctx, cmd)
}

// IsPowerOn checks if a relay is currently on.
// relayNum should be 0 for main relay, 1-8 for specific relays.
func (c *Client) IsPowerOn(ctx context.Context, relayNum int) (bool, error) {
	var resp *PowerResponse
	var err error

	if relayNum == 0 {
		resp, err = c.GetPower(ctx)
	} else {
		resp, err = c.GetPowerN(ctx, relayNum)
	}

	if err != nil {
		return false, err
	}

	return resp.IsOn(relayNum), nil
}

// SetPowerOn turns on a relay.
// relayNum should be 0 for main relay, 1-8 for specific relays.
func (c *Client) SetPowerOn(ctx context.Context, relayNum int) error {
	if relayNum == 0 {
		_, err := c.Power(ctx, PowerOn)
		return err
	}
	_, err := c.PowerN(ctx, relayNum, PowerOn)
	return err
}

// SetPowerOff turns off a relay.
// relayNum should be 0 for main relay, 1-8 for specific relays.
func (c *Client) SetPowerOff(ctx context.Context, relayNum int) error {
	if relayNum == 0 {
		_, err := c.Power(ctx, PowerOff)
		return err
	}
	_, err := c.PowerN(ctx, relayNum, PowerOff)
	return err
}

// TogglePower toggles a relay state.
// relayNum should be 0 for main relay, 1-8 for specific relays.
func (c *Client) TogglePower(ctx context.Context, relayNum int) error {
	if relayNum == 0 {
		_, err := c.Power(ctx, PowerToggle)
		return err
	}
	_, err := c.PowerN(ctx, relayNum, PowerToggle)
	return err
}

// GetCurrentPower returns the current power consumption in Watts.
// This requires a device with power monitoring capability.
func (c *Client) GetCurrentPower(ctx context.Context) (float64, error) {
	// Use Status 10 for sensor data which includes power monitoring
	raw, err := c.ExecuteCommand(ctx, "Status 10")
	if err != nil {
		return 0, err
	}

	var result struct {
		StatusSNS struct {
			Energy struct {
				Power interface{} `json:"Power"`
			} `json:"ENERGY"`
		} `json:"StatusSNS"`
	}

	if err := json.Unmarshal(raw, &result); err != nil {
		return 0, NewError(ErrorTypeParse, "failed to parse status response", err)
	}

	// Handle both int and float values
	switch v := result.StatusSNS.Energy.Power.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case string:
		// Try to parse string as float
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, NewError(ErrorTypeParse, "power value is not a number", err)
		}
		return f, nil
	default:
		return 0, NewError(ErrorTypeParse, "power value has unexpected type", nil)
	}
}

// executePowerCommand is a helper to execute power commands and parse responses.
func (c *Client) executePowerCommand(ctx context.Context, cmd string) (*PowerResponse, error) {
	raw, err := c.ExecuteCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}

	var resp PowerResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, NewError(ErrorTypeParse, "failed to parse power response", err)
	}

	return &resp, nil
}
