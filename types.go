package tasmota

import "encoding/json"

// Version represents the library version.
const Version = "0.1.0"

// UserAgent is the User-Agent header sent with requests.
var UserAgent = "tasmota-go/" + Version

// unmarshalJSON is a helper function to unmarshal JSON with proper error handling.
func unmarshalJSON(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return NewError(ErrorTypeParse, "failed to parse JSON response", err)
	}
	return nil
}

// PowerOnState represents the power state behavior on device boot.
type PowerOnState int

const (
	// PowerOnStateOff keeps relay off after power on.
	PowerOnStateOff PowerOnState = 0
	// PowerOnStateOn turns relay on after power on.
	PowerOnStateOn PowerOnState = 1
	// PowerOnStateToggle toggles relay after power on.
	PowerOnStateToggle PowerOnState = 2
	// PowerOnStateSaved sets relay to last saved state after power on (default).
	PowerOnStateSaved PowerOnState = 3
	// PowerOnStateOnLocked turns relay on and disables further relay control.
	PowerOnStateOnLocked PowerOnState = 4
	// PowerOnStatePulse after a PulseTime period turn relay on (acts as inverted PulseTime mode).
	PowerOnStatePulse PowerOnState = 5
)

// LedState represents the LED behavior mode.
type LedState int

const (
	// LedStateOff disables use of LED as much as possible.
	LedStateOff LedState = 0
	// LedStatePower shows power state on LED (default).
	LedStatePower LedState = 1
	// LedStateMQTTSub shows MQTT subscriptions as a LED blink.
	LedStateMQTTSub LedState = 2
	// LedStatePowerMQTTSub shows power state and MQTT subscriptions as a LED blink.
	LedStatePowerMQTTSub LedState = 3
	// LedStateMQTTPub shows MQTT publications as a LED blink.
	LedStateMQTTPub LedState = 4
	// LedStatePowerMQTTPub shows power state and MQTT publications as a LED blink.
	LedStatePowerMQTTPub LedState = 5
	// LedStateMQTTAll shows all MQTT messages as a LED blink.
	LedStateMQTTAll LedState = 6
	// LedStatePowerMQTTAll shows power state and MQTT messages as a LED blink.
	LedStatePowerMQTTAll LedState = 7
	// LedStateWiFiMQTT LED on when Wi-Fi and MQTT are connected.
	LedStateWiFiMQTT LedState = 8
)

// APMode represents the WiFi access point mode.
type APMode int

const (
	// APModeDisabled disables AP mode.
	APModeDisabled APMode = 0
	// APModeEnabled enables AP mode (default).
	APModeEnabled APMode = 1
	// APModeOpen enables AP mode with no authentication.
	APModeOpen APMode = 2
)

// WiFiConfigMode represents the WiFi configuration mode.
type WiFiConfigMode int

const (
	// WiFiConfigRestart resets Wi-Fi and restarts.
	WiFiConfigRestart WiFiConfigMode = 0
	// WiFiConfigSmartConfig starts smart config for 1 minute.
	WiFiConfigSmartConfig WiFiConfigMode = 1
	// WiFiConfigManager starts WiFi manager for 3 minutes.
	WiFiConfigManager WiFiConfigMode = 2
	// WiFiConfigWPS starts WPS for 1 minute.
	WiFiConfigWPS WiFiConfigMode = 3
	// WiFiConfigRetryDisable disables Wi-Fi auto-restart.
	WiFiConfigRetryDisable WiFiConfigMode = 4
	// WiFiConfigRetryEnable enables Wi-Fi auto-restart.
	WiFiConfigRetryEnable WiFiConfigMode = 5
)

// ResetLevel represents the device reset level.
type ResetLevel int

const (
	// ResetLevelRelay resets relay settings.
	ResetLevelRelay ResetLevel = 1
	// ResetLevelAll resets all settings except Wi-Fi.
	ResetLevelAll ResetLevel = 2
	// ResetLevelAllKeepMQTT resets all settings except Wi-Fi and MQTT.
	ResetLevelAllKeepMQTT ResetLevel = 3
	// ResetLevelWiFi resets Wi-Fi settings.
	ResetLevelWiFi ResetLevel = 4
	// ResetLevelFlashKeepWiFi erases all flash and resets parameters to firmware defaults but keeps Wi-Fi settings.
	ResetLevelFlashKeepWiFi ResetLevel = 5
	// ResetLevelFlash erases all flash and resets parameters to firmware defaults.
	ResetLevelFlash ResetLevel = 6
	// ResetLevelFactoryReboot resets device to firmware defaults and reboots (combines Reset 1 and Restart 1).
	ResetLevelFactoryReboot ResetLevel = 99
)

// RestartReason represents the device restart reason.
type RestartReason int

const (
	// RestartReasonNormal performs a normal restart.
	RestartReasonNormal RestartReason = 1
	// RestartReasonReset performs a reset and restart.
	RestartReasonReset RestartReason = 99
)
