# tasmota-go

A comprehensive Go library for controlling and configuring Tasmota smart devices via HTTP API.

[![Test](https://github.com/kradalby/tasmota-go/workflows/Test/badge.svg)](https://github.com/kradalby/tasmota-go/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/kradalby/tasmota-go)](https://goreportcard.com/report/github.com/kradalby/tasmota-go)
[![GoDoc](https://godoc.org/github.com/kradalby/tasmota-go?status.svg)](https://godoc.org/github.com/kradalby/tasmota-go)
[![Coverage](https://codecov.io/gh/kradalby/tasmota-go/branch/main/graph/badge.svg)](https://codecov.io/gh/kradalby/tasmota-go)

## Features

- **Power Control**: Control up to 8 relays with on/off/toggle commands
- **Status Monitoring**: Query device status, firmware info, network info, and sensor data
- **Device Configuration**: Set friendly names, power-on state, LED state, and more
- **MQTT Configuration**: Configure MQTT broker, topics, authentication, and telemetry
- **Network Configuration**: Set hostname, static IP, DHCP, DNS, and WiFi credentials
- **Power Monitoring**: Read voltage, current, power, and energy consumption
- **Atomic Updates**: Use Backlog commands for atomic multi-setting updates
- **Context Support**: All operations support context for cancellation and timeouts
- **Type-Safe**: Comprehensive type definitions and error handling
- **Well Tested**: >89% test coverage with race detection

## Installation

```bash
go get github.com/kradalby/tasmota-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/kradalby/tasmota-go"
)

func main() {
    // Create a client
    client, err := tasmota.NewClient("192.168.1.100",
        tasmota.WithTimeout(10*time.Second),
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Get device status
    status, err := client.GetStatus(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Device: %s\n", status.FriendlyName[0])

    // Turn on the device
    if err := client.SetPower(ctx, tasmota.PowerOn, 1); err != nil {
        log.Fatal(err)
    }
    fmt.Println("Device turned on!")

    // Get power state
    state, err := client.GetPower(ctx, 1)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Power state: %s\n", state)
}
```

## Usage Examples

### Power Control

```go
// Turn on relay 1
client.SetPower(ctx, tasmota.PowerOn, 1)

// Turn off relay 2
client.SetPower(ctx, tasmota.PowerOff, 2)

// Toggle relay 1
client.SetPower(ctx, tasmota.PowerToggle, 1)

// Get power state
state, err := client.GetPower(ctx, 1)

// Control all relays at once
client.SetPower(ctx, tasmota.PowerOn, 0)
```

### Device Configuration

```go
// Set friendly name
client.SetFriendlyName(ctx, "Living Room Lamp", 1)

// Configure power-on state using typed constants
client.SetPowerOnState(ctx, tasmota.PowerOnStateSaved) // Restore last state on boot

// Set LED state using typed constants
client.SetLedState(ctx, tasmota.LedStatePower) // Show power state on LED

// Apply multiple settings atomically
config := &tasmota.DeviceConfig{
    FriendlyName: []string{"Bedroom Light"},
    PowerOnState: tasmota.PowerOnStateSaved,
    LedState:     tasmota.LedStatePower,
    Sleep:        50,
}
client.ApplyConfig(ctx, config)

// Available PowerOnState constants:
// - PowerOnStateOff: Keep relay off
// - PowerOnStateOn: Turn relay on
// - PowerOnStateToggle: Toggle relay
// - PowerOnStateSaved: Restore last state (default)
// - PowerOnStateOnLocked: Turn on and disable control
// - PowerOnStatePulse: Pulse mode

// Available LedState constants:
// - LedStateOff: Disable LED
// - LedStatePower: Show power state (default)
// - LedStateMQTTSub: Show MQTT subscriptions
// - LedStatePowerMQTTSub: Power + MQTT subscriptions
// - LedStateMQTTPub: Show MQTT publications
// - LedStatePowerMQTTPub: Power + MQTT publications
// - LedStateMQTTAll: Show all MQTT messages
// - LedStatePowerMQTTAll: Power + all MQTT
// - LedStateWiFiMQTT: LED on when connected
```

### MQTT Configuration

```go
// Set MQTT broker
client.SetMQTTHost(ctx, "mqtt.example.com", 1883)

// Set credentials
client.SetMQTTUser(ctx, "username")
client.SetMQTTPassword(ctx, "password")

// Set topic
client.SetMQTTTopic(ctx, "living_room_lamp")

// Apply complete MQTT configuration atomically
mqttConfig := &tasmota.MQTTConfig{
    Host:       "mqtt.example.com",
    Port:       1883,
    User:       "username",
    Password:   "password",
    Topic:      "living_room_lamp",
    FullTopic:  "%prefix%/%topic%/",
    TelePeriod: 300,
}
client.SetMQTTConfig(ctx, mqttConfig)
```

### Network Configuration

```go
// Set hostname
client.SetHostname(ctx, "tasmota-livingroom")

// Configure static IP (using typed IP addresses)
ip := tasmota.MustParseIPAddr("192.168.1.100")
gateway := tasmota.MustParseIPAddr("192.168.1.1")
subnet := tasmota.MustParseIPAddr("255.255.255.0")
client.SetStaticIP(ctx, ip, gateway, subnet)

// Enable DHCP
client.EnableDHCP(ctx, true)

// Set DNS server
dns := tasmota.MustParseIPAddr("8.8.8.8")
client.SetDNSServer(ctx, dns)

// Configure WiFi
client.SetWiFi(ctx, "MySSID", "password", 1)

// Apply complete network configuration atomically
netConfig := &tasmota.NetworkConfig{
    Hostname:  "tasmota-bedroom",
    IPAddress: tasmota.MustParseIPAddr("192.168.1.101"),
    Gateway:   tasmota.MustParseIPAddr("192.168.1.1"),
    Subnet:    tasmota.MustParseIPAddr("255.255.255.0"),
    DNSServer: tasmota.MustParseIPAddr("8.8.8.8"),
    SSID1:     "MySSID",
    Password1: "password",
    UseDHCP:   false,
}
client.SetNetworkConfig(ctx, netConfig)
```

### Status Monitoring

```go
// Get comprehensive status
status, err := client.GetStatus(ctx)

// Get specific status types
firmwareInfo, err := client.GetStatusFirmware(ctx)
networkInfo, err := client.GetNetworkInfo(ctx)
mqttInfo, err := client.GetMQTTInfo(ctx)
sensorData, err := client.GetSensorData(ctx)

// Get power monitoring data
powerInfo, err := client.GetPowerInfo(ctx)
fmt.Printf("Power: %.2fW, Voltage: %.2fV, Current: %.3fA\n",
    powerInfo.Power, powerInfo.Voltage, powerInfo.Current)
```

## API Overview

### Client Creation

- `NewClient(host string, opts ...ClientOption) (*Client, error)`
- `WithAuth(username, password string) ClientOption`
- `WithTimeout(timeout time.Duration) ClientOption`
- `WithHTTPClient(client *http.Client) ClientOption`
- `WithLogger(logger *slog.Logger) ClientOption`

### Power Control

- `SetPower(ctx, state PowerState, relay int) error`
- `GetPower(ctx, relay int) (string, error)`
- `GetPowerInfo(ctx) (*PowerInfo, error)`

### Status

- `GetStatus(ctx) (*StatusInfo, error)`
- `GetStatusFirmware(ctx) (*StatusFirmware, error)`
- `GetNetworkInfo(ctx) (*StatusNetwork, error)`
- `GetMQTTInfo(ctx) (*StatusMQTT, error)`
- `GetSensorData(ctx) (*SensorData, error)`

### Configuration

- `GetConfig(ctx) (*DeviceConfig, error)`
- `SetFriendlyName(ctx, name string, index int) error`
- `SetPowerOnState(ctx, state int) error`
- `SetLedState(ctx, state int) error`
- `SetTelePeriod(ctx, seconds int) error`
- `ApplyConfig(ctx, config *DeviceConfig) error`
- `Reset(ctx, resetType int) error`
- `Restart(ctx, restartType int) error`

### MQTT

- `GetMQTTConfig(ctx) (*MQTTConfig, error)`
- `SetMQTTHost(ctx, host string, port int) error`
- `SetMQTTUser(ctx, user string) error`
- `SetMQTTPassword(ctx, password string) error`
- `SetMQTTTopic(ctx, topic string) error`
- `SetFullTopic(ctx, fullTopic string) error`
- `SetMQTTConfig(ctx, config *MQTTConfig) error`

### Network

- `GetNetworkConfig(ctx) (*NetworkConfig, error)`
- `SetHostname(ctx, hostname string) error`
- `SetStaticIP(ctx, ip, gateway, subnet IPAddr) error`
- `EnableDHCP(ctx, enable bool) error`
- `SetDNSServer(ctx, dnsServer IPAddr) error`
- `SetWiFi(ctx, ssid, password string, slot int) error`
- `SetNetworkConfig(ctx, config *NetworkConfig) error`
- `GetMACAddress(ctx) (MACAddr, error)`
- `Ping(ctx, host string) (bool, error)`

## Development

This project uses Nix flakes for reproducible development environments.

### Setup Development Environment

```bash
# Enter development shell
nix develop

# Or use direnv
echo "use flake" > .envrc
direnv allow
```

The development shell includes:
- Go 1.25
- golangci-lint
- gopls (Go language server)
- delve (debugger)
- Other Go development tools

### Running Tests

```bash
# Run tests with coverage
go test -race -coverprofile=coverage.txt ./...

# View coverage
go tool cover -html=coverage.txt
```

### Linting

```bash
# Run golangci-lint
golangci-lint run

# Auto-fix issues
golangci-lint run --fix
```

## Examples

See the [examples/](examples/) directory for complete example applications:

- [basic](examples/basic/) - Simple power control
- [config](examples/config/) - Device configuration
- [mqtt](examples/mqtt/) - MQTT setup

Run an example:

```bash
TASMOTA_HOST=192.168.1.100 go run examples/basic/main.go
```

## Error Handling

The library provides typed errors for better error handling:

```go
err := client.SetPower(ctx, tasmota.PowerOn, 1)
if err != nil {
    switch {
    case tasmota.IsNetworkError(err):
        // Handle network error
    case tasmota.IsTimeoutError(err):
        // Handle timeout
    case tasmota.IsAuthError(err):
        // Handle authentication error
    default:
        // Handle other errors
    }
}
```

## Logging

The library supports structured logging using Go's `log/slog` package. Pass a custom logger to enable request/response logging:

```go
import (
    "log/slog"
    "os"
)

// Create a logger with debug level
logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

// Pass logger to client
client, err := tasmota.NewClient("192.168.1.100",
    tasmota.WithLogger(logger),
)

// Now all HTTP requests and responses will be logged
```

The logger will output structured logs like:
```
level=DEBUG msg="sending request" method=GET url=http://192.168.1.100/cm?cmnd=Power
level=DEBUG msg="received response" status_code=200 body_length=15 body={"POWER":"ON"}
```

If no logger is provided, no logging will be performed.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## References

- [Tasmota Commands Documentation](https://tasmota.github.io/docs/Commands/)
- [Tasmota HTTP API](https://tasmota.github.io/docs/Commands/#with-web-requests)
