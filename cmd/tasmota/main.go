package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/kradalby/tasmota-go"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func main() {
	var (
		rootFlagSet = flag.NewFlagSet("tasmota", flag.ExitOnError)
		host        = rootFlagSet.String("host", "", "Tasmota device host/IP (required)")
		username    = rootFlagSet.String("username", "", "Basic auth username")
		password    = rootFlagSet.String("password", "", "Basic auth password")
		timeout     = rootFlagSet.Duration("timeout", 10*time.Second, "Request timeout")
		debug       = rootFlagSet.Bool("debug", false, "Enable debug logging")
	)

	root := &ffcli.Command{
		Name:       "tasmota",
		ShortUsage: "tasmota [flags] <subcommand>",
		ShortHelp:  "CLI tool for controlling Tasmota devices",
		LongHelp: `tasmota - Control and configure Tasmota smart devices via HTTP API

This CLI provides comprehensive control over Tasmota devices including:
  - Power control (on/off/toggle for up to 8 relays)
  - Device status and information queries
  - Network configuration (hostname, static IP, DHCP, WiFi)
  - MQTT setup and testing
  - Real-time device information

Authentication:
  If your device requires authentication, use --username and --password flags.

Examples:
  # Get device information
  tasmota --host 192.168.1.100 info

  # Turn on a relay
  tasmota --host 192.168.1.100 power on

  # Configure network
  tasmota --host 192.168.1.100 network set-hostname --hostname tasmota-bedroom

  # Setup MQTT
  tasmota --host 192.168.1.100 mqtt set-config --mqtt-host mqtt.home --mqtt-topic bedroom

  # Enable debug logging
  tasmota --host 192.168.1.100 --debug status

Environment Variables:
  TASMOTA_HOST     - Default host (can be overridden with --host)
  TASMOTA_USERNAME - Default username (can be overridden with --username)
  TASMOTA_PASSWORD - Default password (can be overridden with --password)`,
		FlagSet: rootFlagSet,
		Subcommands: []*ffcli.Command{
			newStatusCmd(host, username, password, timeout, debug),
			newPowerCmd(host, username, password, timeout, debug),
			newInfoCmd(host, username, password, timeout, debug),
			newNetworkCmd(host, username, password, timeout, debug),
			newMQTTCmd(host, username, password, timeout, debug),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func newClient(host, username, password string, timeout time.Duration, debug bool) (*tasmota.Client, error) {
	if host == "" {
		return nil, fmt.Errorf("--host is required")
	}

	opts := []tasmota.ClientOption{
		tasmota.WithTimeout(timeout),
	}

	if username != "" || password != "" {
		opts = append(opts, tasmota.WithAuth(username, password))
	}

	if debug {
		logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
		opts = append(opts, tasmota.WithLogger(logger))
	}

	return tasmota.NewClient(host, opts...)
}
