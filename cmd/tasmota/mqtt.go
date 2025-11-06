package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/kradalby/tasmota-go"
	"github.com/peterbourgon/ff/v3/ffcli"
)

func newMQTTCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	return &ffcli.Command{
		Name:       "mqtt",
		ShortUsage: "tasmota mqtt <subcommand>",
		ShortHelp:  "MQTT configuration",
		Subcommands: []*ffcli.Command{
			newMQTTGetCmd(host, username, password, timeout, debug),
			newMQTTSetHostCmd(host, username, password, timeout, debug),
			newMQTTSetUserCmd(host, username, password, timeout, debug),
			newMQTTSetPasswordCmd(host, username, password, timeout, debug),
			newMQTTSetConfigCmd(host, username, password, timeout, debug),
			newMQTTEnableCmd(host, username, password, timeout, debug),
			newMQTTDisableCmd(host, username, password, timeout, debug),
			newMQTTTestCmd(host, username, password, timeout, debug),
		},
		Exec: func(ctx context.Context, args []string) error {
			return flag.ErrHelp
		},
	}
}

func newMQTTGetCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota mqtt get", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "get",
		ShortUsage: "tasmota mqtt get",
		ShortHelp:  "Get MQTT configuration",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			config, err := client.GetMQTTConfig(ctx)
			if err != nil {
				return fmt.Errorf("failed to get MQTT config: %w", err)
			}

			fmt.Printf("MQTT Configuration:\n")
			fmt.Printf("  Host: %s\n", config.Host)
			fmt.Printf("  Port: %d\n", config.Port)
			fmt.Printf("  User: %s\n", config.User)
			fmt.Printf("  Client: %s\n", config.Client)
			fmt.Printf("  Topic: %s\n", config.Topic)
			fmt.Printf("  Full Topic: %s\n", config.FullTopic)
			fmt.Printf("  Group Topic: %s\n", config.GroupTopic)
			fmt.Printf("  TelePeriod: %d seconds\n", config.TelePeriod)
			fmt.Printf("  Retain: %v\n", config.Retain)

			return nil
		},
	}
}

func newMQTTSetHostCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota mqtt set-host", flag.ExitOnError)
	mqttHost := fs.String("mqtt-host", "", "MQTT broker host (required)")
	mqttPort := fs.Int("mqtt-port", 1883, "MQTT broker port")

	return &ffcli.Command{
		Name:       "set-host",
		ShortUsage: "tasmota mqtt set-host --mqtt-host <host> [--mqtt-port <port>]",
		ShortHelp:  "Set MQTT broker host and port",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			if *mqttHost == "" {
				return fmt.Errorf("--mqtt-host is required")
			}

			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			if err := client.SetMQTTHost(ctx, *mqttHost); err != nil {
				return fmt.Errorf("failed to set MQTT host: %w", err)
			}

			if *mqttPort != 1883 {
				if err := client.SetMQTTPort(ctx, *mqttPort); err != nil {
					return fmt.Errorf("failed to set MQTT port: %w", err)
				}
			}

			fmt.Printf("MQTT broker set to: %s:%d\n", *mqttHost, *mqttPort)
			return nil
		},
	}
}

func newMQTTSetUserCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota mqtt set-user", flag.ExitOnError)
	mqttUser := fs.String("mqtt-user", "", "MQTT username (required)")

	return &ffcli.Command{
		Name:       "set-user",
		ShortUsage: "tasmota mqtt set-user --mqtt-user <username>",
		ShortHelp:  "Set MQTT username",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			if *mqttUser == "" {
				return fmt.Errorf("--mqtt-user is required")
			}

			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			if err := client.SetMQTTUser(ctx, *mqttUser); err != nil {
				return fmt.Errorf("failed to set MQTT user: %w", err)
			}

			fmt.Printf("MQTT user set to: %s\n", *mqttUser)
			return nil
		},
	}
}

func newMQTTSetPasswordCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota mqtt set-password", flag.ExitOnError)
	mqttPassword := fs.String("mqtt-password", "", "MQTT password (required)")

	return &ffcli.Command{
		Name:       "set-password",
		ShortUsage: "tasmota mqtt set-password --mqtt-password <password>",
		ShortHelp:  "Set MQTT password",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			if *mqttPassword == "" {
				return fmt.Errorf("--mqtt-password is required")
			}

			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			if err := client.SetMQTTPassword(ctx, *mqttPassword); err != nil {
				return fmt.Errorf("failed to set MQTT password: %w", err)
			}

			fmt.Println("MQTT password set")
			return nil
		},
	}
}

func newMQTTSetConfigCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota mqtt set-config", flag.ExitOnError)
	mqttHost := fs.String("mqtt-host", "", "MQTT broker host (required)")
	mqttPort := fs.Int("mqtt-port", 1883, "MQTT broker port")
	mqttUser := fs.String("mqtt-user", "", "MQTT username")
	mqttPassword := fs.String("mqtt-password", "", "MQTT password")
	mqttTopic := fs.String("mqtt-topic", "", "MQTT topic")
	telePeriod := fs.Int("tele-period", 300, "Telemetry period in seconds")

	return &ffcli.Command{
		Name:       "set-config",
		ShortUsage: "tasmota mqtt set-config --mqtt-host <host> [flags]",
		ShortHelp:  "Set complete MQTT configuration",
		LongHelp: `Set complete MQTT configuration atomically using Backlog.

This ensures all settings are applied together, which is useful when
configuring MQTT for the first time or changing multiple settings.`,
		FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			if *mqttHost == "" {
				return fmt.Errorf("--mqtt-host is required")
			}

			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			config := &tasmota.MQTTConfig{
				Host:       *mqttHost,
				Port:       *mqttPort,
				User:       *mqttUser,
				Password:   *mqttPassword,
				Topic:      *mqttTopic,
				TelePeriod: *telePeriod,
			}

			if err := client.SetMQTTConfig(ctx, config); err != nil {
				return fmt.Errorf("failed to set MQTT config: %w", err)
			}

			fmt.Println("MQTT configuration applied:")
			fmt.Printf("  Host: %s:%d\n", *mqttHost, *mqttPort)
			if *mqttUser != "" {
				fmt.Printf("  User: %s\n", *mqttUser)
			}
			if *mqttTopic != "" {
				fmt.Printf("  Topic: %s\n", *mqttTopic)
			}
			fmt.Printf("  TelePeriod: %d seconds\n", *telePeriod)

			return nil
		},
	}
}

func newMQTTEnableCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota mqtt enable", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "enable",
		ShortUsage: "tasmota mqtt enable",
		ShortHelp:  "Enable MQTT",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			if err := client.EnableMQTT(ctx, true); err != nil {
				return fmt.Errorf("failed to enable MQTT: %w", err)
			}

			fmt.Println("MQTT enabled")
			return nil
		},
	}
}

func newMQTTDisableCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota mqtt disable", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "disable",
		ShortUsage: "tasmota mqtt disable",
		ShortHelp:  "Disable MQTT",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			if err := client.EnableMQTT(ctx, false); err != nil {
				return fmt.Errorf("failed to disable MQTT: %w", err)
			}

			fmt.Println("MQTT disabled")
			return nil
		},
	}
}

func newMQTTTestCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota mqtt test", flag.ExitOnError)

	return &ffcli.Command{
		Name:       "test",
		ShortUsage: "tasmota mqtt test",
		ShortHelp:  "Test MQTT connection",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			if err := client.TestMQTTConnection(ctx); err != nil {
				return fmt.Errorf("MQTT connection test failed: %w", err)
			}

			fmt.Println("MQTT connection test: SUCCESS")
			return nil
		},
	}
}
