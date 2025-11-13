// Package main provides an example of configuring Tasmota devices.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kradalby/tasmota-go"
)

func main() {
	// Get Tasmota device address from environment
	host := os.Getenv("TASMOTA_HOST")
	if host == "" {
		log.Fatal("TASMOTA_HOST environment variable is required")
	}

	// Create client
	client, err := tasmota.NewClient(host,
		tasmota.WithTimeout(10*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get current configuration
	fmt.Println("Getting current device configuration...")
	config, err := client.GetConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to get config: %v", err) //nolint:gocritic // exitAfterDefer is acceptable in example code
	}
	fmt.Printf("Device Name: %s\n", config.DeviceName)
	fmt.Printf("Friendly Names: %v\n", config.FriendlyName)
	fmt.Printf("Power-on State: %d\n", config.PowerOnState)
	fmt.Printf("LED State: %d\n", config.LedState)

	// Get network configuration
	fmt.Println("\nGetting network configuration...")
	netConfig, err := client.GetNetworkConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to get network config: %v", err)
	}
	fmt.Printf("Hostname: %s\n", netConfig.Hostname)
	fmt.Printf("IP Address: %s\n", netConfig.IPAddress)
	fmt.Printf("Gateway: %s\n", netConfig.Gateway)
	fmt.Printf("Subnet: %s\n", netConfig.Subnet)
	fmt.Printf("DNS: %s\n", netConfig.DNSServer)

	// Set friendly name for relay 1
	fmt.Println("\nSetting friendly name...")
	if err := client.SetFriendlyNameN(ctx, 1, "Living Room Lamp"); err != nil {
		log.Fatalf("Failed to set friendly name: %v", err)
	}
	fmt.Println("Friendly name updated")

	// Configure power-on state (restore last state)
	fmt.Println("\nConfiguring power-on state...")
	if err := client.SetPowerOnState(ctx, tasmota.PowerOnStateSaved); err != nil {
		log.Fatalf("Failed to set power-on state: %v", err)
	}
	fmt.Println("Power-on state set to restore last state")

	// Set LED state (show power state)
	fmt.Println("\nConfiguring LED state...")
	if err := client.SetLedState(ctx, tasmota.LedStatePower); err != nil {
		log.Fatalf("Failed to set LED state: %v", err)
	}
	fmt.Println("LED configured to show power state")

	// Configure telemetry period (10 minutes)
	fmt.Println("\nConfiguring telemetry period...")
	if err := client.SetTelePeriod(ctx, 600); err != nil {
		log.Fatalf("Failed to set telemetry period: %v", err)
	}
	fmt.Println("Telemetry period set to 10 minutes")

	// Apply multiple configuration changes atomically using Backlog
	fmt.Println("\nApplying multiple configuration changes...")
	deviceConfig := &tasmota.DeviceConfig{
		FriendlyName: []string{"Bedroom Light"},
		PowerOnState: tasmota.PowerOnStateSaved,
		LedState:     tasmota.LedStatePower,
		Sleep:        50,
	}
	if err := client.ApplyConfig(ctx, deviceConfig); err != nil {
		log.Fatalf("Failed to apply config: %v", err)
	}
	// Set telemetry period separately
	if err := client.SetTelePeriod(ctx, 300); err != nil {
		log.Fatalf("Failed to set telemetry period: %v", err)
	}
	fmt.Println("Configuration applied successfully")

	// Get updated configuration
	fmt.Println("\nGetting updated configuration...")
	updatedConfig, err := client.GetConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to get updated config: %v", err)
	}
	fmt.Printf("Device Name: %s\n", updatedConfig.DeviceName)
	fmt.Printf("Friendly Names: %v\n", updatedConfig.FriendlyName)
	fmt.Printf("Power-on State: %d\n", updatedConfig.PowerOnState)
	fmt.Printf("LED State: %d\n", updatedConfig.LedState)

	// Test network connectivity
	fmt.Println("\nTesting network connectivity...")
	pingSuccess, err := client.Ping(ctx, "8.8.8.8")
	if err != nil {
		log.Fatalf("Failed to ping: %v", err)
	}
	if pingSuccess {
		fmt.Println("Ping successful")
	} else {
		fmt.Println("Ping failed")
	}

	fmt.Println("\nConfiguration example completed successfully!")
}
