// Package main provides an example of configuring MQTT on Tasmota devices.
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
	// Get configuration from environment
	host := os.Getenv("TASMOTA_HOST")
	if host == "" {
		log.Fatal("TASMOTA_HOST environment variable is required")
	}

	mqttHost := os.Getenv("MQTT_HOST")
	if mqttHost == "" {
		log.Fatal("MQTT_HOST environment variable is required")
	}

	mqttUser := os.Getenv("MQTT_USER")
	mqttPass := os.Getenv("MQTT_PASS")

	// Create client
	client, err := tasmota.NewClient(host,
		tasmota.WithTimeout(10*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get current MQTT configuration
	fmt.Println("Getting current MQTT configuration...")
	mqttConfig, err := client.GetMQTTConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to get MQTT config: %v", err) //nolint:gocritic // exitAfterDefer is acceptable in example code
	}
	fmt.Printf("Current MQTT Host: %s:%d\n", mqttConfig.Host, mqttConfig.Port)
	fmt.Printf("Current MQTT Topic: %s\n", mqttConfig.Topic)
	fmt.Printf("Current Full Topic: %s\n", mqttConfig.FullTopic)

	// Configure MQTT broker
	fmt.Println("\nConfiguring MQTT broker...")
	if err := client.SetMQTTHost(ctx, mqttHost); err != nil {
		log.Fatalf("Failed to set MQTT host: %v", err)
	}
	if err := client.SetMQTTPort(ctx, 1883); err != nil {
		log.Fatalf("Failed to set MQTT port: %v", err)
	}
	fmt.Println("MQTT host configured")

	// Set MQTT credentials if provided
	if mqttUser != "" {
		fmt.Println("\nConfiguring MQTT credentials...")
		if err := client.SetMQTTUser(ctx, mqttUser); err != nil {
			log.Fatalf("Failed to set MQTT user: %v", err)
		}
		if mqttPass != "" {
			if err := client.SetMQTTPassword(ctx, mqttPass); err != nil {
				log.Fatalf("Failed to set MQTT password: %v", err)
			}
		}
		fmt.Println("MQTT credentials configured")
	}

	// Configure MQTT topic
	fmt.Println("\nConfiguring MQTT topic...")
	if err := client.SetTopic(ctx, "living_room_lamp"); err != nil {
		log.Fatalf("Failed to set MQTT topic: %v", err)
	}
	fmt.Println("MQTT topic configured")

	// Configure full topic template
	fmt.Println("\nConfiguring full topic template...")
	if err := client.SetFullTopic(ctx, "%prefix%/%topic%/"); err != nil {
		log.Fatalf("Failed to set full topic: %v", err)
	}
	fmt.Println("Full topic template configured")

	// Set telemetry period (5 minutes)
	fmt.Println("\nConfiguring telemetry period...")
	if err := client.SetTelePeriod(ctx, 300); err != nil {
		log.Fatalf("Failed to set telemetry period: %v", err)
	}
	fmt.Println("Telemetry period set to 5 minutes")

	// Apply complete MQTT configuration atomically using Backlog
	fmt.Println("\nApplying complete MQTT configuration...")
	completeConfig := &tasmota.MQTTConfig{
		Host:       mqttHost,
		Port:       1883,
		User:       mqttUser,
		Password:   mqttPass,
		Client:     "tasmota_living_room",
		Topic:      "living_room_lamp",
		FullTopic:  "%prefix%/%topic%/",
		GroupTopic: "tasmotas",
		Retain:     false,
		TelePeriod: 300,
		Prefix1:    "cmnd",
		Prefix2:    "stat",
		Prefix3:    "tele",
	}
	if err := client.SetMQTTConfig(ctx, completeConfig); err != nil {
		log.Fatalf("Failed to apply MQTT config: %v", err)
	}
	fmt.Println("MQTT configuration applied successfully")

	// Get updated MQTT configuration
	fmt.Println("\nGetting updated MQTT configuration...")
	updatedConfig, err := client.GetMQTTConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to get updated MQTT config: %v", err)
	}
	fmt.Printf("MQTT Host: %s:%d\n", updatedConfig.Host, updatedConfig.Port)
	fmt.Printf("MQTT User: %s\n", updatedConfig.User)
	fmt.Printf("MQTT Client: %s\n", updatedConfig.Client)
	fmt.Printf("MQTT Topic: %s\n", updatedConfig.Topic)
	fmt.Printf("Full Topic: %s\n", updatedConfig.FullTopic)
	fmt.Printf("Group Topic: %s\n", updatedConfig.GroupTopic)
	fmt.Printf("Telemetry Period: %d seconds\n", updatedConfig.TelePeriod)

	// Example: Commands that will be published via MQTT
	fmt.Println("\nExample MQTT topics that will be used:")
	fmt.Printf("  Command topic: cmnd/%s/POWER\n", updatedConfig.Topic)
	fmt.Printf("  Status topic:  stat/%s/POWER\n", updatedConfig.Topic)
	fmt.Printf("  Telemetry topic: tele/%s/STATE\n", updatedConfig.Topic)

	fmt.Println("\nMQTT configuration example completed successfully!")
	fmt.Println("\nNote: Subscribe to the following topics in your MQTT client:")
	fmt.Printf("  stat/%s/#\n", updatedConfig.Topic)
	fmt.Printf("  tele/%s/#\n", updatedConfig.Topic)
}
