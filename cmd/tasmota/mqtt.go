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
			newMQTTSetTopicCmd(host, username, password, timeout, debug),
			newMQTTSetFullTopicCmd(host, username, password, timeout, debug),
			newMQTTSetGroupTopicCmd(host, username, password, timeout, debug),
			newMQTTSetRetainCmd(host, username, password, timeout, debug),
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
	mqttClient := fs.String("mqtt-client", "", "MQTT client name")
	mqttTopic := fs.String("mqtt-topic", "", "MQTT topic")
	mqttFullTopic := fs.String("mqtt-full-topic", "", "MQTT full topic (e.g., %prefix%/%topic%/)")
	mqttGroupTopic := fs.String("mqtt-group-topic", "", "MQTT group topic")
	mqttRetain := fs.Bool("mqtt-retain", false, "Enable MQTT retain")
	telePeriod := fs.Int("tele-period", 300, "Telemetry period in seconds")
	prefix1 := fs.String("prefix1", "", "Command prefix (default: cmnd)")
	prefix2 := fs.String("prefix2", "", "Status prefix (default: stat)")
	prefix3 := fs.String("prefix3", "", "Telemetry prefix (default: tele)")

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
				Client:     *mqttClient,
				Topic:      *mqttTopic,
				FullTopic:  *mqttFullTopic,
				GroupTopic: *mqttGroupTopic,
				Retain:     *mqttRetain,
				TelePeriod: *telePeriod,
				Prefix1:    *prefix1,
				Prefix2:    *prefix2,
				Prefix3:    *prefix3,
			}

			if err := client.SetMQTTConfig(ctx, config); err != nil {
				return fmt.Errorf("failed to set MQTT config: %w", err)
			}

			fmt.Println("MQTT configuration applied:")
			fmt.Printf("  Host: %s:%d\n", *mqttHost, *mqttPort)
			if *mqttUser != "" {
				fmt.Printf("  User: %s\n", *mqttUser)
			}
			if *mqttClient != "" {
				fmt.Printf("  Client: %s\n", *mqttClient)
			}
			if *mqttTopic != "" {
				fmt.Printf("  Topic: %s\n", *mqttTopic)
			}
			if *mqttFullTopic != "" {
				fmt.Printf("  Full Topic: %s\n", *mqttFullTopic)
			}
			if *mqttGroupTopic != "" {
				fmt.Printf("  Group Topic: %s\n", *mqttGroupTopic)
			}
			if *mqttRetain {
				fmt.Printf("  Retain: enabled\n")
			}
			fmt.Printf("  TelePeriod: %d seconds\n", *telePeriod)
			if *prefix1 != "" {
				fmt.Printf("  Prefix1 (command): %s\n", *prefix1)
			}
			if *prefix2 != "" {
				fmt.Printf("  Prefix2 (status): %s\n", *prefix2)
			}
			if *prefix3 != "" {
				fmt.Printf("  Prefix3 (telemetry): %s\n", *prefix3)
			}

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

func newMQTTSetTopicCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota mqtt set-topic", flag.ExitOnError)
	topic := fs.String("topic", "", "MQTT topic (required)")

	return &ffcli.Command{
		Name:       "set-topic",
		ShortUsage: "tasmota mqtt set-topic --topic <topic>",
		ShortHelp:  "Set MQTT topic",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			if *topic == "" {
				return fmt.Errorf("--topic is required")
			}

			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			if err := client.SetTopic(ctx, *topic); err != nil {
				return fmt.Errorf("failed to set topic: %w", err)
			}

			fmt.Printf("MQTT topic set to: %s\n", *topic)
			return nil
		},
	}
}

func newMQTTSetFullTopicCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota mqtt set-full-topic", flag.ExitOnError)
	fullTopic := fs.String("full-topic", "", "MQTT full topic (required, e.g., %prefix%/%topic%/)")

	return &ffcli.Command{
		Name:       "set-full-topic",
		ShortUsage: "tasmota mqtt set-full-topic --full-topic <topic>",
		ShortHelp:  "Set MQTT full topic",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			if *fullTopic == "" {
				return fmt.Errorf("--full-topic is required")
			}

			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			if err := client.SetFullTopic(ctx, *fullTopic); err != nil {
				return fmt.Errorf("failed to set full topic: %w", err)
			}

			fmt.Printf("MQTT full topic set to: %s\n", *fullTopic)
			return nil
		},
	}
}

func newMQTTSetGroupTopicCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota mqtt set-group-topic", flag.ExitOnError)
	groupTopic := fs.String("group-topic", "", "MQTT group topic (required)")

	return &ffcli.Command{
		Name:       "set-group-topic",
		ShortUsage: "tasmota mqtt set-group-topic --group-topic <topic>",
		ShortHelp:  "Set MQTT group topic",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			if *groupTopic == "" {
				return fmt.Errorf("--group-topic is required")
			}

			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			if err := client.SetGroupTopic(ctx, *groupTopic); err != nil {
				return fmt.Errorf("failed to set group topic: %w", err)
			}

			fmt.Printf("MQTT group topic set to: %s\n", *groupTopic)
			return nil
		},
	}
}

func newMQTTSetRetainCmd(host, username, password *string, timeout *time.Duration, debug *bool) *ffcli.Command {
	fs := flag.NewFlagSet("tasmota mqtt set-retain", flag.ExitOnError)
	retain := fs.Bool("retain", true, "Enable MQTT retain")

	return &ffcli.Command{
		Name:       "set-retain",
		ShortUsage: "tasmota mqtt set-retain [--retain=true|false]",
		ShortHelp:  "Set MQTT retain flag",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			client, err := newClient(*host, *username, *password, *timeout, *debug)
			if err != nil {
				return err
			}

			if err := client.SetMQTTRetain(ctx, *retain); err != nil {
				return fmt.Errorf("failed to set MQTT retain: %w", err)
			}

			fmt.Printf("MQTT retain: %v\n", *retain)
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
