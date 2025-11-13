// Package main provides a basic example of controlling Tasmota devices.
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

	// Create client with optional authentication
	client, err := tasmota.NewClient(host,
		tasmota.WithTimeout(10*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Use context with timeout for operations
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get device information
	fmt.Println("Getting device information...")
	status, err := client.GetDeviceInfo(ctx)
	if err != nil {
		log.Fatalf("Failed to get status: %v", err) //nolint:gocritic // exitAfterDefer is acceptable in example code
	}
	fmt.Printf("Device: %s (Module: %d)\n", status.FriendlyName[0], status.Module)

	firmwareInfo, err := client.GetFirmwareInfo(ctx)
	if err != nil {
		log.Fatalf("Failed to get firmware info: %v", err)
	}
	fmt.Printf("Firmware: %s\n", firmwareInfo.Version)

	// Get current power state
	fmt.Println("\nGetting current power state...")
	powerState, err := client.GetPowerN(ctx, 1)
	if err != nil {
		log.Fatalf("Failed to get power state: %v", err)
	}
	fmt.Printf("Power state: %s\n", powerState.Power1)

	// Turn on the device
	fmt.Println("\nTurning on...")
	if err := client.SetPowerOn(ctx, 1); err != nil {
		log.Fatalf("Failed to turn on: %v", err)
	}
	fmt.Println("Device is now ON")

	// Wait a bit
	time.Sleep(2 * time.Second)

	// Turn off the device
	fmt.Println("\nTurning off...")
	if err := client.SetPowerOff(ctx, 1); err != nil {
		log.Fatalf("Failed to turn off: %v", err)
	}
	fmt.Println("Device is now OFF")

	// Toggle the device
	fmt.Println("\nToggling power...")
	if err := client.TogglePower(ctx, 1); err != nil {
		log.Fatalf("Failed to toggle: %v", err)
	}

	// Get final power state
	finalState, err := client.GetPowerN(ctx, 1)
	if err != nil {
		log.Fatalf("Failed to get final power state: %v", err)
	}
	fmt.Printf("Final power state: %s\n", finalState.Power1)

	// Get power monitoring info (if supported)
	fmt.Println("\nGetting power monitoring data...")
	sensorData, err := client.GetSensorData(ctx)
	switch {
	case err != nil:
		log.Printf("Sensor data not available: %v", err)
	case sensorData.Energy != nil:
		fmt.Printf("Power: %.2fW, Voltage: %.2fV, Current: %.3fA\n",
			sensorData.Energy.Power, sensorData.Energy.Voltage, sensorData.Energy.Current)
	default:
		fmt.Println("Power monitoring not available on this device")
	}

	fmt.Println("\nExample completed successfully!")
}
