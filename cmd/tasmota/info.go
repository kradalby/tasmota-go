package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func newInfoCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota info", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "info",
		ShortUsage: "tasmota info",
		ShortHelp:  "Get comprehensive device information",
		LongHelp: `Get comprehensive information about a Tasmota device.

This command displays:
  - Device name and friendly names
  - Firmware version, build date, core version
  - Network information (hostname, IP, MAC, WiFi power)
  - Current state (uptime, heap memory, load average)
  - WiFi connection details (SSID, RSSI, signal strength, channel)

This is a convenient way to get a full overview of a device's status
without querying individual status categories.

Example:
  tasmota --host 192.168.1.100 info`,
		FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			// Get device info
			devInfo, err := client.GetDeviceInfo(ctx)
			if err != nil {
				return fmt.Errorf("failed to get device info: %w", err)
			}

			fmt.Printf("Device Information:\n")
			fmt.Printf("  Device Name: %s\n", devInfo.DeviceName)
			if len(devInfo.FriendlyName) > 0 {
				for i, name := range devInfo.FriendlyName {
					fmt.Printf("  Friendly Name %d: %s\n", i+1, name)
				}
			}
			fmt.Printf("  Module: %d\n", devInfo.Module)
			fmt.Printf("  Topic: %s\n", devInfo.Topic)

			// Get firmware info
			fwInfo, err := client.GetFirmwareInfo(ctx)
			if err != nil {
				return fmt.Errorf("failed to get firmware info: %w", err)
			}

			fmt.Printf("\nFirmware:\n")
			fmt.Printf("  Version: %s\n", fwInfo.Version)
			fmt.Printf("  Build Date: %s\n", fwInfo.BuildDateTime)
			fmt.Printf("  Core: %s\n", fwInfo.Core)
			fmt.Printf("  SDK: %s\n", fwInfo.SDK)
			fmt.Printf("  CPU Frequency: %d MHz\n", fwInfo.CpuFrequency)
			fmt.Printf("  Hardware: %s\n", fwInfo.Hardware)

			// Get network info
			netInfo, err := client.GetNetworkInfo(ctx)
			if err != nil {
				return fmt.Errorf("failed to get network info: %w", err)
			}

			fmt.Printf("\nNetwork:\n")
			fmt.Printf("  Hostname: %s\n", netInfo.Hostname)
			fmt.Printf("  IP Address: %s\n", netInfo.IPAddress)
			fmt.Printf("  Gateway: %s\n", netInfo.Gateway)
			fmt.Printf("  Subnet: %s\n", netInfo.Subnetmask)
			fmt.Printf("  DNS: %s\n", netInfo.DNSServer)
			fmt.Printf("  MAC Address: %s\n", netInfo.Mac)
			fmt.Printf("  WiFi Power: %.1f dBm\n", netInfo.WifiPower)

			// Get state
			state, err := client.GetState(ctx)
			if err != nil {
				return fmt.Errorf("failed to get state: %w", err)
			}

			fmt.Printf("\nState:\n")
			fmt.Printf("  Uptime: %s (%d seconds)\n", state.Uptime, state.UptimeSec)
			fmt.Printf("  Heap: %d KB\n", state.Heap)
			fmt.Printf("  Sleep Mode: %s\n", state.SleepMode)
			fmt.Printf("  Load Average: %d%%\n", state.LoadAvg)

			if state.Wifi != nil {
				fmt.Printf("\nWiFi:\n")
				fmt.Printf("  SSID: %s\n", state.Wifi.SSId)
				fmt.Printf("  BSSID: %s\n", state.Wifi.BSSId)
				fmt.Printf("  Channel: %d\n", state.Wifi.Channel)
				fmt.Printf("  Mode: %s\n", state.Wifi.Mode)
				fmt.Printf("  RSSI: %d dBm\n", state.Wifi.RSSI)
				fmt.Printf("  Signal: %d%%\n", state.Wifi.Signal)
				fmt.Printf("  Link Count: %d\n", state.Wifi.LinkCount)
				fmt.Printf("  Downtime: %s\n", state.Wifi.Downtime)
			}

			return nil
		},
	}
}
