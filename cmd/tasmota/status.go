package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func newStatusCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota status", flag.ExitOnError)
	category := fs.Int("category", 0, "Status category (0=all, 1-11=specific)")
	jsonOutput := fs.Bool("json", false, "Output raw JSON")

	return &ffcli.Command{
		Name:       "status",
		ShortUsage: "tasmota status [flags]",
		ShortHelp:  "Get device status information",
		LongHelp: `Get status information from a Tasmota device.

Categories:
  0  - All status information (default)
  1  - Device parameters
  2  - Firmware version
  3  - Logging information
  4  - Memory information
  5  - Network information
  6  - MQTT information
  7  - Time information
  8  - Sensor information
  9  - Power threshold
  10 - Sensor information
  11 - State information`,
		FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			resp, err := client.Status(ctx, *category)
			if err != nil {
				return fmt.Errorf("failed to get status: %w", err)
			}

			if *jsonOutput {
				data, err := json.MarshalIndent(resp, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal JSON: %w", err)
				}
				fmt.Println(string(data))
				return nil
			}

			// Print human-readable output
			if resp.Status != nil {
				fmt.Printf("Device Name: %s\n", resp.Status.DeviceName)
				if len(resp.Status.FriendlyName) > 0 {
					fmt.Printf("Friendly Name: %s\n", resp.Status.FriendlyName[0])
				}
				fmt.Printf("Module: %d\n", resp.Status.Module)
			}

			if resp.StatusFWR != nil {
				fmt.Printf("\nFirmware:\n")
				fmt.Printf("  Version: %s\n", resp.StatusFWR.Version)
				fmt.Printf("  Core: %s\n", resp.StatusFWR.Core)
				fmt.Printf("  SDK: %s\n", resp.StatusFWR.SDK)
			}

			if resp.StatusNET != nil {
				fmt.Printf("\nNetwork:\n")
				fmt.Printf("  Hostname: %s\n", resp.StatusNET.Hostname)
				fmt.Printf("  IP Address: %s\n", resp.StatusNET.IPAddress)
				fmt.Printf("  Gateway: %s\n", resp.StatusNET.Gateway)
				fmt.Printf("  MAC: %s\n", resp.StatusNET.Mac)
				fmt.Printf("  WiFi Power: %.1f dBm\n", resp.StatusNET.WifiPower)
			}

			if resp.StatusSTS != nil {
				fmt.Printf("\nState:\n")
				fmt.Printf("  Uptime: %s\n", resp.StatusSTS.Uptime)
				fmt.Printf("  Heap: %d KB\n", resp.StatusSTS.Heap)
				if resp.StatusSTS.POWER != "" {
					fmt.Printf("  Power: %s\n", resp.StatusSTS.POWER)
				}
				if resp.StatusSTS.Wifi != nil {
					fmt.Printf("\nWiFi:\n")
					fmt.Printf("  SSID: %s\n", resp.StatusSTS.Wifi.SSId)
					fmt.Printf("  RSSI: %d dBm\n", resp.StatusSTS.Wifi.RSSI)
					fmt.Printf("  Signal: %d%%\n", resp.StatusSTS.Wifi.Signal)
					fmt.Printf("  Channel: %d\n", resp.StatusSTS.Wifi.Channel)
				}
			}

			if resp.StatusMQT != nil {
				fmt.Printf("\nMQTT:\n")
				fmt.Printf("  Host: %s:%d\n", resp.StatusMQT.MqttHost, resp.StatusMQT.MqttPort)
				fmt.Printf("  Client: %s\n", resp.StatusMQT.MqttClient)
				fmt.Printf("  User: %s\n", resp.StatusMQT.MqttUser)
			}

			return nil
		},
	}
}
